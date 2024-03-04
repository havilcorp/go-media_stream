package video

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"go-media-stream/internal/config"
	"go-media-stream/internal/entity"
	"go-media-stream/internal/store/local"
	"go-media-stream/internal/store/mysql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	config *config.Config
	store  *local.LocalStore
	db     *mysql.Mysql
}

func NewHandler(config *config.Config, store *local.LocalStore, db *mysql.Mysql) *handler {
	return &handler{
		config: config,
		store:  store,
		db:     db,
	}
}

func (h *handler) Register(r *chi.Mux) {
	// r.Get("/film/*", h.GetFilm)
	r.Post("/video/{id}/time", h.SetTime)
	r.Post("/video/{id}/audio", h.SetAudio)
	r.Get("/video/{id}", h.GetVideo)
}

func (h *handler) GetVideo(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	video, err := h.db.GetVideo(r.Context(), id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	audio, err := h.db.GetAudioByVideoId(r.Context(), id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "video.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct {
		Video entity.VideoModel
		Audio []entity.AudioModel
	}{video, audio}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) SetAudio(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var jsonData struct {
		AudioId int `json:"audio_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.db.SetVideoAudio(r.Context(), id, jsonData.AudioId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) GetFilm(rw http.ResponseWriter, r *http.Request) {
	prefix := chi.URLParam(r, "*")
	prefixDecode, err := url.QueryUnescape(prefix)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "film.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(rw, struct{ Name string }{prefixDecode}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) SetTime(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	data := struct {
		Time float32 `json:"time"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	h.db.SetVideoTime(r.Context(), id, data.Time)
	rw.WriteHeader(http.StatusOK)
}
