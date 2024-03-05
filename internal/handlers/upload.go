package handlers

import (
	"errors"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/handlers/middleware"
	"go-media-stream/internal/services"

	"github.com/go-chi/chi/v5"
)

type UploadHandler struct {
	uploadService *services.UploadService
}

func NewUploadHandler(uploadService *services.UploadService) *UploadHandler {
	return &UploadHandler{
		uploadService: uploadService,
	}
}

func (h *UploadHandler) Register(r *chi.Mux) {
	r.With(middleware.JwtAuthMiddleware).Get("/upload", h.GetTemplateUpload)
	r.With(middleware.JwtAuthMiddleware).Post("/upload", h.UploadFile)
}

func (h *UploadHandler) GetTemplateUpload(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "upload.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct{}{}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UploadHandler) UploadFile(rw http.ResponseWriter, r *http.Request) {
	const sixGB = 6 * 1024 * 1024 * 1024
	r.ParseMultipartForm(sixGB)

	name := r.FormValue("name")
	if name == "" {
		http.Error(rw, errors.New("name").Error(), http.StatusBadRequest)
		return
	}

	film, handler, err := r.FormFile("film")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer film.Close()

	if handler.Size > (sixGB) {
		http.Error(rw, errors.New("file size exceeds the limit").Error(), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(rw.Header().Get("USER_ID"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.uploadService.Upload(r.Context(), userId, name, &film)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidVideoName) {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	rw.WriteHeader(http.StatusCreated)
}
