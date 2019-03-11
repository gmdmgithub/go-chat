package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gmdmgithub/chat/trace"
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

// templateHandler - struct for html templates
type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("cient templatehandler")
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
	port := myEnv["PORT"]
	if len(port) == 0 {
		port = ":8080"
	}

	oauthInit()

	//first apprach with HandleFunc
	//http.HandleFunc("/", func (http.ResponseWriter, r *http.Request){w.Write([]byte(`html .....`})

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/room", r)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "index.html"}))

	// start a room conversation (single room in new thread)
	go r.run()

	log.Printf("Main func running on port %s", port)
	// start the web server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
