package main

import (
	"html/template"
	"log"
	"net/http"
)

type Message struct {
	User string
	Text string
}

type Page struct {
	Thread   string
	Username string
	Messages []Message
}

var pages = make(map[string][]Message)

func main() {
	pages["Test"] = []Message{{User: "system", Text: "Hi, write your message"}}
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/send/", sendHandler)
	http.ListenAndServe(":8080", nil)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("index.html")
	temp.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	cookie := http.Cookie{Name: "username", Value: username}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/view/Test", http.StatusFound)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	thread := r.URL.Path[len("/view/"):]
	username, _ := r.Cookie("username")
	if username.Value == "" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
	}
	if pages[thread] == nil {
		log.Println("nil")
	}
	p := &Page{Thread: thread, Username: username.Value, Messages: pages[thread]}
	temp, _ := template.ParseFiles("room.html")
	err := temp.Execute(w, p)
	if err != nil {
		log.Println(err)

	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("send")
	thread := r.URL.Path[len("/send/"):]
	message := r.FormValue("message")
	username, _ := r.Cookie("username")
	pages[thread] = append(pages[thread], Message{User: username.Value, Text: message})
	viewHandler(w, r)

}
