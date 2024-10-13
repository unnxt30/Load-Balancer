package main

import (
	"fmt"
	"net/http"
)


func ServerRespond(w http.ResponseWriter, r *http.Request){
	msg := "server_3 says hi"
	fmt.Println(msg)
	w.Write([]byte(msg))
}

// func HealthCheck(w http.ResponseWriter, r *http.Request){
// 	helper.RespondWithJSON(w, 200, "Healthy :-)")
// }

