package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	helper "github.com/unnxt30/Load-Balancer/helpers"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func BalancerResponse(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load();
	port := os.Getenv("SERVER_PORT")
	if err != nil{
		log.Fatal("Couldn't load environment variables")
	}

	userIP := ReadUserIP(r)
	msg := fmt.Sprintf("Recieved request from: %v", userIP)
	go helper.RespondWithJSON(w, 200, map[string]string{"message": msg})

	server_route := fmt.Sprintf("http://localhost:%v/", port)

	client := &http.Client{}
	forwardRequest, err := http.NewRequest(r.Method, server_route, r.Body)
	if err != nil{
		helper.RespondWithError(w, 400, "Could not forward the request")
		return
	}

	forwardRequest.Header = r.Header

	resp , err := client.Do(forwardRequest)
	if err != nil{
		helper.RespondWithError(w, 400, "Could not forward the request")
		return
	}

	resp_msg, _ := io.ReadAll(resp.Body)

	body := string(resp_msg)

	helper.RespondWithJSON(w, 200, map[string]string{"server":body})

}

