package stream

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"go-media-stream/internal/config"
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
	r.Get("/stream/video/{id}", h.GetStreamVideo)
	r.Get("/stream/audio/{id}", h.GetStreamAudio)
}

func (h *handler) GetStreamVideo(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	video, err := h.db.GetVideo(r.Context(), id)
	if err != nil {
		http.Error(rw, "File not found", http.StatusNotFound)
		return
	}
	file, err := os.Open(path.Join("uploads", video.Name, "video.mp4"))
	if err != nil {
		http.Error(rw, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.ServeContent(rw, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func (h *handler) GetStreamAudio(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	audio, err := h.db.GetAudio(r.Context(), id)
	if err != nil {
		http.Error(rw, "File not found 1", http.StatusNotFound)
		return
	}
	video, err := h.db.GetVideo(r.Context(), audio.VideoId)
	if err != nil {
		http.Error(rw, "File not found 2", http.StatusNotFound)
		return
	}
	file, err := os.Open(path.Join("uploads", video.Name, fmt.Sprintf("%s.mp3", audio.Idx)))
	if err != nil {
		http.Error(rw, "File not found 3", http.StatusNotFound)
		return
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.ServeContent(rw, r, fileInfo.Name(), fileInfo.ModTime(), file)
}
