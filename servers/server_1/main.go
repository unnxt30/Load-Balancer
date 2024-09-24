package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


func main(){
	err := godotenv.Load("../../.env")

	if err != nil{
		log.Fatal(".env error")
	}

	SERVER_PORT := os.Getenv("SERVER_PORT_1")

	mux := http.NewServeMux()
	var server http.Server

	mux.HandleFunc("/", ServerRespond)

	server.Addr = fmt.Sprintf(":%v", SERVER_PORT);
	server.Handler = mux
	server.ListenAndServe()
}