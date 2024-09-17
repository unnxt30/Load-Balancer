package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)


func main() {
	err := godotenv.Load();
	if err != nil{
		log.Fatal("Error fetching environment variables!\n");
	}

	BALANCER_PORT := os.Getenv("BALANCER_PORT");

	fmt.Println(BALANCER_PORT)

	mux := http.NewServeMux()
	var server http.Server

	mux.HandleFunc("/", BalancerResponse);

	server.Addr = ":8080"
	server.Handler = mux
	server.ListenAndServe()
}
