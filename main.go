package main

import (
	"fmt"
	"html/template"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
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

func create(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("pages/create.html",
		"pages/header.html",
		"pages/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	temp.ExecuteTemplate(w, "create", nil)
}
func save_article(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	shortDesc := r.FormValue("short_desc")
	content := r.FormValue("content")

	if title == "" || shortDesc == "" || content == "" {
		fmt.Fprintf(w, "Всі поля повинні бути заповненими")
	}
	db, err := sql.Open("mysql", "root:mysql@/newsportal_golang_db")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := fmt.Sprintf("insert into `articles` (`title`, `short_desc`, `content`) values('%s', '%s', '%s')",
		title, shortDesc, content)

	insert, err := db.Query(query)

	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	fs := http.FileServer(http.Dir("./res/"))
	http.Handle("/res/", http.StripPrefix("/res", fs))

	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
