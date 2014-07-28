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
	Messages []Message
}

var pages = make(map[string][]Message)

func main() {
	pages["Test"] = []Message{{User: "system", Text: "Hi, write your message"}}
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/send/", sendHandler)
	http.ListenAndServe(":8080", nil)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	thread := r.URL.Path[len("/view/"):]
	if pages[thread] == nil {
		log.Println("nil")
	}
	p := &Page{Thread: thread, Messages: pages[thread]}
	temp, _ := template.ParseFiles("index.html")
	err := temp.Execute(w, p)
	if err != nil {
		log.Println(err)

	}
}

func sendHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("send")
	thread := r.URL.Path[len("/send/"):]
	message := r.FormValue("message")
	pages[thread] = append(pages[thread], Message{User: "user1", Text: message})
	http.Redirect(w, r, "/view/"+thread, http.StatusFound)
}
