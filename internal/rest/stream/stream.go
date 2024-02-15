package stream

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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
	r.Get("/stream/video", h.GetStreamVideo)
	r.Get("/stream/audio", h.GetStreamAudio)
}

func (h *handler) GetStreamVideo(rw http.ResponseWriter, r *http.Request) {
	folder := r.URL.Query().Get("folder")
	if folder == "" {
		http.Error(rw, errors.New("FOLDER_NOT_FOUND").Error(), http.StatusBadGateway)
		return
	}
	file, err := os.Open(path.Join("uploads", folder, "video.mp4"))
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
	folder := r.URL.Query().Get("folder")
	audio := r.URL.Query().Get("audio")
	if folder == "" {
		http.Error(rw, errors.New("FOLDER_NOT_FOUND").Error(), http.StatusBadGateway)
		return
	}
	if audio == "" {
		http.Error(rw, errors.New("AUDIO_NOT_FOUND").Error(), http.StatusBadGateway)
		return
	}
	file, err := os.Open(path.Join("uploads", folder, fmt.Sprintf("%s.mp3", audio)))
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
