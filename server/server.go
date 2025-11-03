package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	ShutdownServer *http.Server
	AgentId        = make(chan string)
)

func init() {
	http.HandleFunc("/Dahlia/ticket", postHandler)
	http.HandleFunc("/web/nginx", VisitNginx)
	http.HandleFunc("/web/tomcat", VisitTomcat)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only use post..", http.StatusMethodNotAllowed)
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("get post failed")
	}
	defer r.Body.Close()
	pid := struct {
		Id string `json:"id"`
	}{}
	err = json.Unmarshal([]byte(string(body)), &pid)
	if err != nil {
		log.Fatalf("body unmarshal faild :%v", err)
	}
	AgentId <- pid.Id
	log.Println(pid.Id)
}

func Server() {
	server := &http.Server{
		Addr:         ":1021",
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	ShutdownServer = server

	log.Println("server open.. post is 1021...")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server failed, err:%v\n", err)
	}

}

// web/nginx路由转发触发的func
func VisitNginx(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/web/nginx.html")
}

// web/tomcat路由转发触发的func
func VisitTomcat(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./server/web/tomcat.html")
}
