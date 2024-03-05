package handlers

import (
	"net/http"
	"path"
	"text/template"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/middleware"
	"go-media-stream/internal/utils"

	"github.com/go-chi/chi/v5"
)

type HomeHandler struct {
	video VideoProvider
	audio AudioProvider
}

func NewHomeHandler(video VideoProvider, audio AudioProvider) *HomeHandler {
	return &HomeHandler{
		video: video,
		audio: audio,
	}
}

func (h *HomeHandler) Register(r *chi.Mux) {
	r.With(middleware.JwtAuthMiddleware).Get("/", h.MainPage)
}

func (h *HomeHandler) MainPage(rw http.ResponseWriter, r *http.Request) {
	videos, err := h.video.GetVideos(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// jsonData, err := json.Marshal(videos)
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// rw.Write(jsonData)
	root, err := utils.GetProjectRoot()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles(
		path.Join(root, "templates", "index.html"),
		path.Join(root, "templates", "header.html"),
		path.Join(root, "templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct{ Videos []domain.Video }{*videos}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
