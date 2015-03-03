package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var cfg Config
var db *sqlx.DB

func sendError(rw http.ResponseWriter, status int, message string) {
	envelope := map[string]interface{}{
		"errors": [1]interface{}{
			map[string]interface{}{
				"status": fmt.Sprintf("%d",status),
				"title": message,
				"detail": message,
			},
		},
	}
	render.New().JSON(rw, status, envelope)
}

func sendData(rw http.ResponseWriter, status int, data interface{}) {
	envelope := map[string]interface{}{
		"data": data,
	}
	render.New().JSON(rw, status, envelope)
}

func getIP(r *http.Request) string {
    if ipProxy := r.Header.Get("X-Forwarded-For"); len(ipProxy) > 0 {
    	ips := strings.Split(ipProxy, ", ")
        return ips[0]
    }
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}

func main() {
	err := LoadConfigInto(&cfg, "config.gcfg")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Loaded config.")

	// set up database connection
	db, err = sqlx.Open("mysql", cfg.Database.User + ":" + cfg.Database.Password + "@/" + cfg.Database.Database)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Opened database.")

	// DB fields are capitalized in the same way as Go structs, so mapper is a no-op
	db.MapperFunc(func(s string) string {
		return s
	})

    // Heroku uses env var to specify port
    port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Server.Port
	}

	// any requests on the regular HTTP port get automatically redirected to the secure home page
	// don't just redirect HTTP to HTTPS, because that doesn't train the user to not use HTTP on the first request
	// TODO: commented out so I can run Apache on 80 on my computer
	// TODO: this didn't work, maybe I have too many web servers running, also Chrome gets pretty confused with port numbers
	// go http.ListenAndServe(cfg.Server.Host + ":80", http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
	// 	// TODO: secure app shouldn't have to be at the root of the domain
	// 	http.Redirect(w, req, "https://" + cfg.Server.Host + ":" + port, http.StatusMovedPermanently)
	// }))

    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)

    api := r.PathPrefix("/api").Subrouter()

    api.HandleFunc("/user", PostUserHandler).Methods("POST")
	// api.HandleFunc("/user/{handle}", func (rw http.ResponseWriter, r *http.Request) {
	// 	// TODO: 
	// }).Methods("GET")
	// api.HandleFunc("/user/{handle}", func (rw http.ResponseWriter, r *http.Request) {
	// 	// TODO: 
	// }).Methods("PUT")
	// api.HandleFunc("/user/{handle}", func (rw http.ResponseWriter, r *http.Request) {
	// 	// TODO: 
	// }).Methods("DELETE")

	api.HandleFunc("/note", ListNotesHandler).Methods("GET")
	api.HandleFunc("/note", PostNoteHandler).Methods("POST")
	api.HandleFunc("/note/{id}", GetNoteHandler).Methods("GET")
	api.HandleFunc("/note/{id}", PutNoteHandler).Methods("PUT")
	api.HandleFunc("/note/{id}", DeleteNoteHandler).Methods("DELETE")

    api.HandleFunc("/token", PostTokenHandler).Methods("POST")
    api.HandleFunc("/token/{token}", DeleteTokenHandler).Methods("DELETE")

    api.HandleFunc("/user/{handle}/host/{host}", GetUserHostHandler).Methods("GET")
    api.HandleFunc("/user/{handle}/host", PostUserHostHandler).Methods("POST")
    api.HandleFunc("/guest", PostGuestHandler).Methods("POST")

    hostname := cfg.Server.Host + ":" + port
    log.Println("Listening on " + hostname + ".")
	http.ListenAndServeTLS(hostname, cfg.Server.Certificate, cfg.Server.Key, r)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	rend := render.New()
	rend.HTML(rw, http.StatusOK, "login", nil)
}
