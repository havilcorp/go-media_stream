package upload

import (
	"errors"
	"html/template"
	"net/http"
	"path"

	"go-media-stream/internal/common/ffmpeg"
	"go-media-stream/internal/config"
	"go-media-stream/internal/store/local"
	"go-media-stream/internal/store/mysql"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	config *config.Config
	store  *local.LocalStore
	db     *mysql.Mysql
	ff     *ffmpeg.FFMPEG
}

func NewHandler(
	config *config.Config,
	store *local.LocalStore,
	db *mysql.Mysql,
	ff *ffmpeg.FFMPEG,
) *handler {
	return &handler{
		config: config,
		store:  store,
		db:     db,
		ff:     ff,
	}
}

func (h *handler) Register(r *chi.Mux) {
	r.Get("/upload", h.GetTemplateUpload)
	r.Post("/upload", h.UploadFile)
}

func (h *handler) GetTemplateUpload(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "upload.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct{}{}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) UploadFile(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	film, handler, err := r.FormFile("film")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer film.Close()

	name := r.FormValue("name")
	if name == "" {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.db.IsValidVideoName(r.Context(), name)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if !ok {
		http.Error(rw, errors.New("name is exists").Error(), http.StatusBadRequest)
		return
	}

	preview, _, err := r.FormFile("preview")
	if err == nil {
		defer preview.Close()
		err = h.store.SaveFile(name, "preview.jpg", preview)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	} else if !errors.Is(err, http.ErrMissingFile) {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.store.SaveFile(name, handler.Filename, film)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	h.ff.Add(ffmpeg.JobModel{
		Name:     name,
		FileName: handler.Filename,
	})

	rw.WriteHeader(http.StatusCreated)
}
