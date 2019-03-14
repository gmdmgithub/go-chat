package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Picture       string `json:"picture"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
}

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  string
)

func initOAuth() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost" + myEnv["PORT"] + "/auth/callback/google",
		ClientID:     myEnv["GOOGLE_CLIENT_ID"],
		ClientSecret: myEnv["GOOGLE_CLIENT_SECRET"],
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	oauthStateString = myEnv["SECRET_KEY"]
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Auth - start ...")
	cookie, _ := r.Cookie("auth")
	log.Printf("Coockie values %v", cookie)
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error - panic system stops
		panic(err.Error())
	} else {
		// success check cooki and finally call the next handler
		// Cookie data
		// data, err := base64.StdEncoding.DecodeString(cookie.Value)
		// if err != nil {
		// 	fmt.Printf("Error: %v", err)
		// }

		// var user googleUser
		// json.Unmarshal(data, &user)
		// log.Printf("User data are: %v", user)
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth -required user have to be login to enter
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		log.Println("String to short!!") // TODO ? redirect to not found page?
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		url := googleOauthConfig.AuthCodeURL(oauthStateString)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		log.Printf("Login with provider %s", provider)
	case "callback":
		user, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
		if err != nil {
			fmt.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//fmt.Fprintf(w, "Content: %s\n", content)
		userStr, err := json.Marshal(user)
		log.Printf("Content: %v\n", string(userStr))

		// data := "some string to decode"
		// decode := base64.StdEncoding.EncodeToString([]byte(data))
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: base64.StdEncoding.EncodeToString(userStr),
			Path:  "/"})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func getUserInfo(state string, code string) (*googleUser, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	var user *googleUser
	err = json.Unmarshal(contents, &user)
	log.Printf("JSON of user data is %v", user)
	if err != nil {
		log.Printf("Error unmarshaling Google user %s\n", err.Error())
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return user, nil
}
