package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/middleware"

	"github.com/go-chi/chi/v5"
)

type VideoHandler struct {
	video VideoProvider
	audio AudioProvider
}

func NewVideoHandler(video VideoProvider, audio AudioProvider) *VideoHandler {
	return &VideoHandler{
		video: video,
		audio: audio,
	}
}

func (h *VideoHandler) Register(r *chi.Mux) {
	r.With(middleware.JwtAuthMiddleware).Get("/video/{id}", h.GetVideo)
	r.With(middleware.JwtAuthMiddleware).Post("/video/{id}/time", h.SetTime)
	r.With(middleware.JwtAuthMiddleware).Post("/video/{id}/audio", h.SetAudio)
}

func (h *VideoHandler) GetVideo(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	video, err := h.video.GetVideoById(r.Context(), id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNotFound)
		return
	}
	audio, err := h.audio.GetAudioByVideoId(r.Context(), id)
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
		Video domain.Video
		Audio []domain.Audio
	}{*video, *audio}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *VideoHandler) SetTime(rw http.ResponseWriter, r *http.Request) {
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
	h.video.SetTime(r.Context(), id, data.Time)
	rw.WriteHeader(http.StatusOK)
}

func (h *VideoHandler) SetAudio(rw http.ResponseWriter, r *http.Request) {
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
	err = h.video.SetAudio(r.Context(), id, jsonData.AudioId)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}
