package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"text/template"

	"go-media-stream/internal/domain"
	"go-media-stream/internal/log"
	"go-media-stream/internal/services"

	"github.com/go-chi/chi/v5"
)

type AuthHandler struct {
	logger      *log.Logger
	authService *services.AuthService
}

func NewAuthHandler(logger *log.Logger, authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		logger:      logger,
		authService: authService,
	}
}

func (h *AuthHandler) Register(r *chi.Mux) {
	r.Get("/auth", h.GetAuthPage)
	r.Post("/auth", h.Auth)
}

func (h *AuthHandler) GetAuthPage(rw http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		path.Join("templates", "authorization.html"),
		path.Join("templates", "header.html"),
		path.Join("templates", "footer.html"),
	))
	if err := tmpl.Execute(rw, struct{}{}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Auth(rw http.ResponseWriter, r *http.Request) {
	data := struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: переделать функцю на две
	token, err := h.authService.SignUp(r.Context(), data.Login, data.Password)
	if err != nil {
		if errors.Is(err, domain.ErrWrongPassword) {
			http.Error(rw, err.Error(), http.StatusUnauthorized)
			return
		} else {
			h.logger.Error(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	resp := struct {
		Token string `json:"token"`
	}{Token: token}
	b, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(b)
	if err != nil {
		h.logger.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
