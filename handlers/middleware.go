package handlers

import (
	"encoding/base64"
	"net/http"
	"net/url"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			data, _ := base64.StdEncoding.DecodeString(cookie.Value)
			values, err := url.ParseQuery(string(data))
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusFound)
			} else {
				if values.Get("username") != "xxx" || values.Get("password") != "xxx" {
					http.Redirect(w, r, "/login", http.StatusFound)
				} else {
					next.ServeHTTP(w, r)
				}
			}
		}
	})
}
