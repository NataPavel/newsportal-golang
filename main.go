package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("pages/index.html",
		"pages/header.html",
		"pages/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	temp.ExecuteTemplate(w, "index", nil)
}

func handleFunc() {
	fs := http.FileServer(http.Dir("./res/"))
	http.Handle("/res/", http.StripPrefix("/res", fs))
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
