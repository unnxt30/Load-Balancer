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
	port_1 := os.Getenv("SERVER_PORT_1")
	port_2 := os.Getenv("SERVER_PORT_2")
	
	server_route_1 := fmt.Sprintf("http://localhost:%v/", port_1)
	server_route_2 := fmt.Sprintf("http://localhost:%v/", port_2)
	
	// serverList := map[string]bool{
	// 	server_route_1 : true,
	// 	server_route_2 : false,
	// }

	var serverList []Server;

	server1 := Server{
		serverURL: server_route_1,
		isHealthy: true,
	}

	server2 := Server{
		serverURL: server_route_2,
		isHealthy: true,
	}

	serverList = append(serverList, server1)
	serverList = append(serverList, server2)


	balConfig := BalancerConfig{
		currentServer: serverList[0],
		serverStack: serverList,

	};

	fmt.Println(BALANCER_PORT)

	mux := http.NewServeMux()
	var server http.Server

	mux.HandleFunc("/", balConfig.BalancerResponse);

	server.Addr = ":8080"
	server.Handler = mux
	server.ListenAndServe()
}
