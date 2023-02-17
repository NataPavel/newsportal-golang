package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Article struct {
	Id                        uint
	Title, Content, ShortDesc string
}

func index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("pages/index.html",
		"pages/header.html",
		"pages/footer.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Db conection
	db, err := sql.Open("mysql", "root:mysql@/newsportal_golang_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Selecting data
	//Db query
	query := "select * from `articles`"
	res, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	// Scanning DB and Filling it into a post list
	posts := []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.ShortDesc, &post.Content)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}
	temp.ExecuteTemplate(w, "index", posts)
}

func contacts(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("pages/contacts.html",
		"pages/header.html",
		"pages/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	temp.ExecuteTemplate(w, "contacts", nil)
}

func about(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("pages/about.html",
		"pages/header.html",
		"pages/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	temp.ExecuteTemplate(w, "about", nil)
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

	// validation
	if title == "" || shortDesc == "" || content == "" {
		fmt.Fprintf(w, "Всі поля повинні бути заповненими")
	} else {
		// Db conection
		db, err := sql.Open("mysql", "root:mysql@/newsportal_golang_db")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Db query
		query := fmt.Sprintf("insert into `articles` (`title`, `short_desc`, `content`) values('%s', '%s', '%s')",
			title, shortDesc, content)

		insert, err := db.Query(query)
		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func article_details(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	temp, err := template.ParseFiles("pages/article_details.html",
		"pages/header.html",
		"pages/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	// Db conection
	db, err := sql.Open("mysql", "root:mysql@/newsportal_golang_db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Selecting data
	//Db query
	query := fmt.Sprintf("select * from `articles` where `id` = '%s'", vars["id"])
	res, err := db.Query(query)
	if err != nil {
		panic(err)
	}

	// Scanning DB
	var showArticle = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.ShortDesc, &post.Content)
		if err != nil {
			panic(err)
		}
		showArticle = post
	}

	temp.ExecuteTemplate(w, "article_details", showArticle)
}

func handleFunc() {
	fs := http.FileServer(http.Dir("./res/"))
	http.Handle("/res/", http.StripPrefix("/res", fs))

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/article_details/{id:[0-9]+}", article_details).Methods("GET")
	rtr.HandleFunc("/contacts", contacts).Methods("GET")
	rtr.HandleFunc("/about", about).Methods("GET")
	http.Handle("/", rtr)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
