package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/X3NOOO/maud/db"
	"github.com/X3NOOO/maud/types"
)

type maud_context struct {
	db *db.DB
}

func (ctx *maud_context) registerPOST(w http.ResponseWriter, r *http.Request) {
	var account types.RegisterPOST

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil || account.Nick == "" || account.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("failed to decode the body of your request"))
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}
	
	response, rerr := ctx.db.Register(&account)
	if rerr != nil {
		w.WriteHeader(rerr.StatusCode)
		w.Write([]byte(rerr.Error()))
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("couldn't construct response"))
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) loginPOST(w http.ResponseWriter, r *http.Request) {
	var account types.LoginPOST

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil || account.Nick == "" || account.Password == "" {
		http.Error(w, "failed to decode the body of your request", http.StatusBadRequest)
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}
	
	response, rerr := ctx.db.Login(&account)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		log.Println("Failed register attempt from", r.RemoteAddr)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) alivePOST(w http.ResponseWriter, r *http.Request) {

}

func (ctx *maud_context) switchesPOST(w http.ResponseWriter, r *http.Request) {

}

func (ctx *maud_context) switchesGET(w http.ResponseWriter, r *http.Request) {
	if id := r.PathValue("id"); id != "" {
		w.Write([]byte(id))
	} else {
		w.Write([]byte("duap"))
	}
}

func jsonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")

		content_type := r.Header.Get("Content-Type")
		
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

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := db.InitDatabase("")
	if err != nil {
		log.Fatalf("Failed to initialise the database: %v\n", err)
	}

	ctx := maud_context{db}

	http.Handle("POST /register", jsonMiddleware(http.HandlerFunc(ctx.registerPOST)))
	http.Handle("POST /login", jsonMiddleware(http.HandlerFunc(ctx.loginPOST)))
	http.Handle("POST /alive", jsonMiddleware(http.HandlerFunc(ctx.alivePOST)))
	http.Handle("POST /switches", jsonMiddleware(http.HandlerFunc(ctx.switchesPOST)))
	http.Handle("GET /switches/{id...}", jsonMiddleware(http.HandlerFunc(ctx.switchesGET)))

	http.HandleFunc("/tea", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("Im a teapot."))
	})

	log.Fatalln(http.ListenAndServe(":1337", nil))
}
