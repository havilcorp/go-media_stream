package middleware

import (
	"net/http"

	"go-media-stream/internal/utils"

	"github.com/sirupsen/logrus"
)

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bearer, err := r.Cookie("Bearer")
		if err != nil {
			http.Redirect(rw, r, "/auth", http.StatusSeeOther)
			return
		}
		payload, err := utils.VerifyJWT(bearer.Value)
		if err != nil {
			logrus.Error(err)
			http.Redirect(rw, r, "/auth", http.StatusSeeOther)
			return
		}
		userID, err := payload.GetSubject()
		if err != nil {
			logrus.Error(err)
			http.Redirect(rw, r, "/auth", http.StatusSeeOther)
			return
		}
		rw.Header().Add("USER_ID", userID)
		next.ServeHTTP(rw, r)
	})
}
