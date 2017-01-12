package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zhangmingkai4315/weichat-golang-backend/config"
	"github.com/zhangmingkai4315/weichat-golang-backend/routers"
	"log"
	"net/http"
	"time"
)

func main() {
	configObj, err := config.NewConfig()
	if err != nil {
		log.Printf("Loading Config File Error: %s", err)
		return
	}
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/api", routers.PostMessage).Methods("POST")
	r.HandleFunc("/api", routers.WeiChatValidate).Methods("GET")
	http.Handle("/",r)

	hostAndPort := fmt.Sprintf("%s:%d", configObj.Server.Host, configObj.Server.Port)
	server := &http.Server{
		Handler:      r,
		Addr:         hostAndPort,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Println("Server listen at", hostAndPort)
	log.Fatal(server.ListenAndServe())
}
