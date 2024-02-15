package upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
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

func Contains(items []string, item string) bool {
	for _, elem := range items {
		if elem == item {
			return true
		}
	}
	return false
}

func (h *handler) Register(r *chi.Mux) {
	// r.Post("/folder/create", h.CreateFolder)
	// r.Post("/folder/upload", h.UploadFile)
	r.Get("/upload", h.UploadFilm)
}

func (h *handler) UploadFilm(w http.ResponseWriter, r *http.Request) {
	// fp := path.Join("templates", "upload.html")
	// tmpl, err := template.ParseFiles(fp)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "upload.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(w, struct{}{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) CreateFolder(rw http.ResponseWriter, r *http.Request) {
	data := struct {
		Folder string `json:"folder"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.store.FolderCreate(data.Folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) UploadFile(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	folderFolder := r.FormValue("folder")
	if folderFolder == "" {
		http.Error(rw, errors.New("NOT_FOUND path").Error(), http.StatusBadGateway)
		return
	}
	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		cType := fileHeader.Header["Content-Type"][0]
		fmt.Println(cType)
		if !Contains(h.config.AllowTypes, cType) {
			http.Error(rw, errors.New("TYPE_NOT_VALID").Error(), http.StatusBadGateway)
			return
		}
	}
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		err = h.store.Create(&file, folderFolder, fileHeader.Filename)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	rw.WriteHeader(http.StatusOK)
}
