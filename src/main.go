package main

import (
	"fmt"
	"net/http"
)

func main()  {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/greeting", handlerGreeting)


	var address = "localhost:8080"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func handlerGreeting(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}
