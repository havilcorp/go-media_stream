package films

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"

	"s3/internal/config"
	"s3/internal/store/local"
	"s3/internal/store/mysql"

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
	// r.Post("/time", h.SetTime)
	r.Get("/films", h.GetFilms)
}

type GetFilmModel struct {
	Name string
	Url  string
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
	data := struct {
		Name  string  `json:"name"`
		Audio string  `json:"audio"`
		Time  float32 `json:"time"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	h.db.SetFilmTime(r.Context(), data.Name, data.Audio, data.Time)
	rw.WriteHeader(http.StatusOK)
}

func (h *handler) GetFilms(rw http.ResponseWriter, r *http.Request) {
	films, err := h.db.GetFilms(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "films.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(rw, struct{ Films []mysql.FilmModel }{films}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
