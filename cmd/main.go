package main

import (
	"fmt"
	"net/http"

	"s3/internal/config"
	"s3/internal/rest/films"
	"s3/internal/rest/index"
	"s3/internal/rest/stream"
	"s3/internal/rest/upload"
	"s3/internal/store/local"
	"s3/internal/store/mysql"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	store := local.NewLocalStore()
	config := config.NewConfig()
	mysql, err := mysql.NewMysql(config.DBConnect)
	if err != nil {
		panic(err)
	}

	index.NewHandler(config, store, mysql).Register(r)
	films.NewHandler(config, store, mysql).Register(r)
	stream.NewHandler(config, store, mysql).Register(r)
	upload.NewHandler(config, store, mysql).Register(r)

	fmt.Println("Server started")
	http.ListenAndServe(":8080", r)
}
