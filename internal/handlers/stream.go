package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"

	"go-media-stream/internal/handlers/middleware"
	"go-media-stream/internal/log"
	"go-media-stream/internal/services"

	"github.com/go-chi/chi/v5"
)

type StreamHandler struct {
	logger *log.Logger
	video  *services.VideoService
	audio  *services.AudioService
}

func NewStreamHandler(
	logger *log.Logger,
	video *services.VideoService,
	audio *services.AudioService,
) *StreamHandler {
	return &StreamHandler{
		logger: logger,
		video:  video,
		audio:  audio,
	}
}

func (h *StreamHandler) Register(r *chi.Mux) {
	r.With(middleware.JwtAuthMiddleware).Get("/stream/video/{id}", h.GetVideoStream)
	r.With(middleware.JwtAuthMiddleware).Get("/stream/audio/{id}", h.GetAudioStream)
}

func (h *StreamHandler) GetVideoStream(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	video, err := h.video.GetVideoById(r.Context(), id)
	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusNotFound)
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
		h.logger.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.ServeContent(rw, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func (h *StreamHandler) GetAudioStream(rw http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	audio, err := h.audio.GetAudioById(r.Context(), id)
	if err != nil {
		http.Error(rw, "File not found 1", http.StatusNotFound)
		return
	}
	video, err := h.video.GetVideoById(r.Context(), audio.VideoId)
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
