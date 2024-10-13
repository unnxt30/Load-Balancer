package main

import (
	"fmt"
	"net/http"

	helper "github.com/unnxt30/Load-Balancer/helpers"
)


func ServerRespond(w http.ResponseWriter, r *http.Request){
	msg := "server_1 says hi"
	fmt.Println(msg)
	helper.RespondWithJSON(w, 200, msg)
}

