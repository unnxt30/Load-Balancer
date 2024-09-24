package main

import (
	"fmt"
	"net/http"
)


func ServerRespond(w http.ResponseWriter, r *http.Request){
	msg := "server_2 says hi"
	fmt.Println(msg)
	w.Write([]byte(msg))
}

