package fileserver

import "net/http"

func AuthMiddleware(fs *FileServer, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if fs.auth == nil {
			next.ServeHTTP(w, r)
			return
		}

		// verify basic auth
		if !fs.auth.authenticate(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
