package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	helper "github.com/unnxt30/Load-Balancer/helpers"
)

type Server struct{
	serverURL string
	isHealthy bool
}

type BalancerConfig struct{
	currentServer Server 
	serverStack []Server
}


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

func createForwardRequest(route string, r *http.Request) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, route, r.Body)

	if err != nil {
		return nil, err
	}

	return req, nil
}

/* 
A list of all servers.
A list of active servers. Dynamic/Updated after each health check
*/


func (b *BalancerConfig) healthCheck(done chan bool) error{
	var deadServers int
	for _, v := range b.serverStack{
		resp, err := http.Get(v.serverURL)
		if err != nil{
			return err
		}
		if resp.StatusCode != 200 {
			code := fmt.Sprintf("%v", resp.StatusCode)
			v.isHealthy = false
			deadServers += 1
			return errors.New(code) 
		}
	}

	if deadServers == len(b.serverStack){
		done <- true
	}
	return nil
}


func (b *BalancerConfig) BalancerResponse(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load();
	if err != nil{
		log.Fatal("Couldn't load environment variables")
	}

	// Load-Balancer Response
	userIP := ReadUserIP(r)
	msg := fmt.Sprintf("Recieved request from: %v", userIP)
	go helper.RespondWithJSON(w, 200, map[string]string{"message": msg})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	
	for {
		select {
		case <-done:
			fmt.Println("All servers down!")
			return
		case <-ticker.C:
			go b.healthCheck(done)
		}

		i := 0
		for !b.serverStack[i].isHealthy {
			b.serverStack = append(b.serverStack, b.currentServer)	
			b.serverStack = b.serverStack[i:]
			i++
		}
	
		b.currentServer = b.serverStack[0]

		b.serverStack = append(b.serverStack, b.currentServer)	
		b.serverStack = b.serverStack[1:]


		// Round Robin Logic
		client := &http.Client{}
		forwardRequest, err := createForwardRequest(b.currentServer.serverURL, r)

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
}

