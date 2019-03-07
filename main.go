package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

// templateHendler - struct for html templates
type templateHendler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHendler) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
		err := t.template.Execute(w, nil)
		if err != nil {
			log.Fatalf("Templates not readed with exception: %v", err)
		}
	})
}

var myEnv map[string]string

func loadDotEnv() {
	var err error
	myEnv, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	loadDotEnv()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<meta http-equiv="X-UA-Compatible" content="ie=edge">
				<title>Chat app</title>
			</head>
			<body>
				<h1>Chat app stats here!</h1>
			</body>
			</html>
		`))
	})

	port := myEnv["PORT"]
	if len(port) == 0 {
		port = ":8080"
	}
	log.Printf("Main func running on port %s", port)
	// start the web server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
