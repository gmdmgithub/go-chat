package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

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

	tm := time.Now()

	formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second())

	userData := map[string]interface{}{
		"Host": r.Host,
		"Date": formatted,
	}
	cookie, err := r.Cookie("auth")
	if err == nil {
		data, err := base64.StdEncoding.DecodeString(cookie.Value)
		if err == nil {

			var user googleUser
			json.Unmarshal(data, &user)
			// log.Printf("User data are: %v", user)
			userData["UserData"] = user
			t.template.Execute(w, userData)
		} else {
			fmt.Printf("Error: %v\n", err)
			t.template.Execute(w, r)
		}

	} else {
		fmt.Printf("Error: %v\n", err)
		t.template.Execute(w, r)
	}

}

func main() {

	loadDotEnv()
	port := myEnv["PORT"]
	if len(port) == 0 {
		port = ":8080"
	}

	initOAuth()
	// r := newRoom(UseAuthAvatar)
	r := newRoom(UseGravatar)
	r.tracer = trace.New(os.Stdout)

	initRouts(r)

	// start a room conversation (single room in new thread)
	go r.run()

	log.Printf("Main func running on port %s", port)
	// start the web server
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
