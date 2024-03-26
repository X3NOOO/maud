package main

import (
	"log"
	"mime"
	"net/http"
)

func optionsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func (ctx *maud_context) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", ctx.config.Maud.ACAO)
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Add("Access-Control-Allow-Credentials", "true")

		next.ServeHTTP(w, r)
	})
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		content_type := r.Header.Get("Content-Type")
		if content_type == "" {
			http.Error(w, "Content-Type not provided", http.StatusUnauthorized)
			return
		}

		mt, _, err := mime.ParseMediaType(content_type)
		if err != nil {
			http.Error(w, "malformed Content-Type header", http.StatusBadRequest)
			return
		}

		if mt != "application/json" {
			http.Error(w, "expected application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (ctx *maud_context) authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")
		if authorization == "" {
			authorization_cookie, err := r.Cookie("Authorization")
			if err == http.ErrNoCookie {
				http.Error(w, "authorization not provided", http.StatusUnauthorized)
				return
			} else if err != nil {
				log.Println("error while authorizing:", err)
				return
			}
			r.Header.Add("authorization", authorization_cookie.Value) // todo: other functions depend on this header
			authorization = authorization_cookie.Value
		}

		ok, rerr := ctx.db.Authorize(authorization)
		if rerr != nil {
			http.Error(w, rerr.Err.Error(), rerr.StatusCode)
			return
		}
		if !ok {
			http.Error(w, "invalid authorization token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
