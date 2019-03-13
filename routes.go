package main

import (
	"log"
	"net/http"
)

func initRouts(r *room) {

	//first apprach with HandleFunc
	//http.HandleFunc("/", func (http.ResponseWriter, r *http.Request){w.Write([]byte(`html .....`})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.Handle("/login", &templateHandler{filename: "login.html"})

	http.HandleFunc("/auth/", loginHandler)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "index.html"}))

	//socket communication
	http.Handle("/room", r)

	logout()

}

func logout() {
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Logout user ...")
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}
