package main

import "net/http"

func accessTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Token") == "" {
			errorHandler(w, r, http.StatusUnauthorized, "Invalid access token provided!")
			return
		}

		if r.Header.Get("Access-Token") != cdnConfig.AccessToken {
			errorHandler(w, r, http.StatusUnauthorized, "Invalid access token provided!")
			return
		}

		next.ServeHTTP(w, r)
	})
}
