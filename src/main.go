package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", handlerIndex)
	http.HandleFunc("/greeting", handlerGreeting)
	http.HandleFunc("/form", handlerForm)
	http.HandleFunc("/result", handlerResult)
	http.HandleFunc("/upload/form", formUpload)
	http.HandleFunc("/upload-process", uploadProcess)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	var address = "localhost:8080"
	fmt.Printf("server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func uploadProcess(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		file, handler, err := request.FormFile("file")
		defer file.Close()

		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		fileName := handler.Filename

		if alias := request.Form.Get("alias"); alias != "" {
			fileName = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
		}

		dir, err := os.Getwd()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		fileLocation :=  filepath.Join(dir, "files", fileName)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(targetFile, file); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		defer targetFile.Close()

		writer.Write([]byte("done"))
	}

	http.Error(writer, "", http.StatusBadRequest)
}

func formUpload(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		formUpload := path.Join("views", "form-upload.html")
		temp := template.Must(template.New("upload").ParseFiles(formUpload))

		if err := temp.Execute(writer, ""); err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		return
	}

	http.Error(writer, "", http.StatusInternalServerError)
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
