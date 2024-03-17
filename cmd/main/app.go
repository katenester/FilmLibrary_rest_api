package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/katenester/FilmLibrary_rest_api/internal/user"
	"log"
	"net"
	"net/http"
	"time"
)

func main() {
	log.Println("create router")
	// Роутер (маршрутезатор)
	router := httprouter.New()
	log.Println("register user handler")
	handler := user.NewHandler()
	handler.Register(router)
	start(router)
}

func start(router *httprouter.Router) {
	log.Println("start application")
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	//http.ListenAndServe(":8080", router)
	server := &http.Server{
		Handler: router,
		//Ставим таймаутд на чтение и запись
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("server is listening port 8080")
	log.Fatalln(server.Serve(listener))
}
