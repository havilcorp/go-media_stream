package index

import (
	"net/http"
	"path"
	"text/template"

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
	r.Get("/", h.MainPage)
}

func (h *handler) MainPage(rw http.ResponseWriter, r *http.Request) {
	video, err := h.db.GetVideos(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "index.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct{ Videos []entity.VideoModel }{video}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
