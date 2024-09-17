package main

import (
	"fmt"
	"net/http"
)


func ServerRespond(w http.ResponseWriter, r *http.Request){
	msg := "server says hi"
	fmt.Println(msg)
	w.Write([]byte(msg))
}

