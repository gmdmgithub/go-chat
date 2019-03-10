package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var myEnv map[string]string

func loadDotEnv() {
	var err error
	myEnv, err = godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// templateHendler - struct for html templates
type templateHendler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHendler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("cient templatehendler")
	t.once.Do(func() {
		log.Printf("read file once %v", t.filename)
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	err := t.template.Execute(w, r)

	if err != nil {
		log.Fatalf("Templates not loaded with error: %v", err)
	}

}

func main() {

	loadDotEnv()
	//first apprach with HandleFunc
	//http.HandleFunc("/", func (http.ResponseWriter, r *http.Request){w.Write([]byte(`html .....`})

	r := newRoom()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/", &templateHendler{filename: "index.html"})

	http.Handle("/room", r)
	// start a room conversation (single room in new thread)
	go r.run()

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
