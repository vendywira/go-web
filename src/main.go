package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
)

func main() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/greeting", handlerGreeting)
	http.HandleFunc("/form", handlerForm)
	http.HandleFunc("/result", handlerResult)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	var address = "localhost:8080"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}
func handlerResult(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		var filepath = path.Join("views", "message.html")
		var temp = template.Must(template.New("result").ParseFiles(filepath))
		if err := request.ParseForm(); err != nil {
			http.Error(writer, "form must not be empty", http.StatusBadRequest)
			return
		}

		name := request.Form.Get("name")
		message := request.Form.Get("message")
		var data = map[string]string{"name": name, "message": message}

		if err := temp.Execute(writer, data); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	http.Error(writer, "", http.StatusBadRequest)
}

func handlerForm(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		var filepath = path.Join("views", "form.html")
		var temp = template.Must(template.New("form").ParseFiles(filepath))

		if err := temp.Execute(writer, nil); err != nil{
			http.Error(writer, err.Error(), http.StatusBadRequest)
		}

		return
	}

	http.Error(writer, "", http.StatusBadRequest)
}

func handlerGreeting(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome!"))
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "index.html")
	var tmpl, err = template.ParseFiles(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data = map[string]interface{}{
		"title": "Learning Golang Web",
		"name":  "Vendy Wiranatha",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
