package index

import (
	"net/http"
	"net/url"
	"path"
	"text/template"

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
	// r.Get("/*", h.MainPage)
}

func (h *handler) MainPage(rw http.ResponseWriter, r *http.Request) {
	pathEncode := chi.URLParam(r, "*")
	pathDecode, err := url.QueryUnescape(pathEncode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	folders, files, err := h.store.Ls(pathDecode)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	data := struct{ Folders, Files map[string]string }{folders, files}
	if err := tmpl.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}
