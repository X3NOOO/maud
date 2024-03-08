package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/X3NOOO/maud/db"
	"github.com/X3NOOO/maud/types"
)

const CONFIG_FILE string = "./maud.toml"

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

	response, rerr := ctx.db.Register(account)
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
		log.Println("Failed login attempt from", r.RemoteAddr)
		return
	}

	response, rerr := ctx.db.Login(account)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		log.Println("Failed login attempt from", r.RemoteAddr)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		log.Println("Failed login attempt from", r.RemoteAddr)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) statusGET(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")
	response, rerr := ctx.db.Status(authorization)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) alivePOST(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")
	response, rerr := ctx.db.UpdateAlive(authorization)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) switchesPOST(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")

	var switch_body types.SwitchesPOST
	err := json.NewDecoder(r.Body).Decode(&switch_body)
	if err != nil {
		http.Error(w, "failed to decode the body of your request", http.StatusBadRequest)
		log.Println("Failed login attempt from", r.RemoteAddr)
		return
	}

	response, rerr := ctx.db.AddSwitch(authorization, switch_body)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) switchesGET(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")

	id_str := r.PathValue("id")
	if id_str == "" {
		id_str = "-1"
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	response, rerr := ctx.db.GetSwitch(authorization, int64(id))
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) switchesDELETE(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")

	id_str := r.PathValue("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	response, rerr := ctx.db.DeleteSwitch(authorization, int64(id))
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func (ctx *maud_context) switchesPATCH(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")

	id_str := r.PathValue("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var switch_body types.SwitchesPATCH
	err = json.NewDecoder(r.Body).Decode(&switch_body)
	if err != nil {
		http.Error(w, "failed to decode the body of your request", http.StatusBadRequest)
		log.Println("Failed login attempt from", r.RemoteAddr)
		return
	}

	response, rerr := ctx.db.UpdateSwitch(authorization, int64(id), switch_body)
	if rerr != nil {
		http.Error(w, rerr.Error(), rerr.StatusCode)
		return
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "couldn't construct response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response_json)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	db, err := db.InitDatabase("anon@/maud", "maud")
	if err != nil {
		log.Fatalf("Failed to initialise the database: %v\n", err)
	}

	ctx := maud_context{}
	ctx.db = db
	defer ctx.db.Close()

	http.Handle("POST /register", jsonMiddleware(http.HandlerFunc(ctx.registerPOST)))
	http.Handle("POST /login", jsonMiddleware(http.HandlerFunc(ctx.loginPOST)))
	http.Handle("GET /status", ctx.authorizationMiddleware(http.HandlerFunc(ctx.statusGET)))
	http.Handle("POST /alive", ctx.authorizationMiddleware(http.HandlerFunc(ctx.alivePOST)))
	http.Handle("POST /switches", jsonMiddleware(ctx.authorizationMiddleware(http.HandlerFunc(ctx.switchesPOST))))
	http.Handle("GET /switches/{id...}", ctx.authorizationMiddleware(http.HandlerFunc(ctx.switchesGET)))
	http.Handle("DELETE /switches/{id...}", ctx.authorizationMiddleware(http.HandlerFunc(ctx.switchesDELETE)))
	http.Handle("PATCH /switches/{id...}", ctx.authorizationMiddleware(http.HandlerFunc(ctx.switchesPATCH)))

	http.HandleFunc("/tea", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("I run on a teapot."))
	})

	go ctx.watchdog()

	log.Fatalln(http.ListenAndServe(":1337", nil))
}
