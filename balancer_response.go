package main

import (
	"errors"
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


func (b *BalancerConfig) healthCheck(quit chan bool) error {
    deadServers := 0
    for i := range b.serverStack {
        resp, err := http.Get(b.serverStack[i].serverURL)
        if err != nil {
            b.serverStack[i].isHealthy = false
            deadServers++
            continue
        }
        if resp.StatusCode != 200 {
            b.serverStack[i].isHealthy = false
            deadServers++
        } else {
            b.serverStack[i].isHealthy = true
        }
    }

	if deadServers == len(b.serverStack){
		quit <- false
		return errors.New("all servers dead nigga")
	}

	return nil
}


func (b *BalancerConfig) BalancerResponse(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load();
	if err != nil{
		log.Fatal("Couldn't load environment variables")
	}

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)
	go func() {
		for {
		   select {
			case <- ticker.C:
				b.healthCheck(quit)
			case <- quit:
				ticker.Stop()
				return
			}
		}
	 }()	


	i := 0
    for i < len(b.serverStack) {
        if !b.serverStack[i].isHealthy {
            b.serverStack = append(b.serverStack, b.serverStack[i])
            b.serverStack = append(b.serverStack[:i], b.serverStack[i+1:]...)
        } else {
            b.currentServer = b.serverStack[i]
            i++
            break
        }
    }

    b.serverStack = append(b.serverStack, b.currentServer)   
    b.serverStack = b.serverStack[1:]		


	client := &http.Client{}
	forwardRequest, err := createForwardRequest(b.currentServer.serverURL, r)

	if err != nil{
		helper.RespondWithError(w, 400, "Could not forward the request")
		return
	}

	forwardRequest.Header = r.Header

	resp , err := client.Do(forwardRequest)
	if err != nil{
		helper.RespondWithError(w, 404, err.Error())
		return
	}

	resp_msg, _ := io.ReadAll(resp.Body)

	body := string(resp_msg)

	if resp.StatusCode != 200{
		helper.RespondWithError(w, 400, "server not live.")
	}

	helper.RespondWithJSON(w, 200, map[string]string{"server":body, "code":resp.Status})

}

