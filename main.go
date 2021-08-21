package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Topic struct {
	ID          int
	Title       string
	Body        string
	CreatedDate string
}

type Comment struct {
	ID          int
	Body        string
	CreatedDate string
}

var DbConnection *sql.DB
var layout = "2006-01-02 15:04:05"
var templates = template.Must(template.ParseFiles("./templates/index.html", "./templates/new.html", "./templates/topic.html"))
var validPath = regexp.MustCompile("^/(save|topic)/([a-zA-Z0-9]+)$")

// func (topic *Topic)save()error {

// }

func rootHandler(w http.ResponseWriter, r *http.Request) {
	cmd := "SELECT * FROM topic"
	topic_rows, _ := DbConnection.Query(cmd)
	var topics []Topic
	for topic_rows.Next() {
		var topic Topic
		err := topic_rows.Scan(&topic.Title, &topic.Body, &topic.CreatedDate)
		if err != nil {
			log.Println(err)
		}
		topics = append(topics, topic)
	}
	err := templates.ExecuteTemplate(w, "index.html", topics)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
	}

	topic_title := m[2]
	cmd := "SELECT * FROM topic WHERE title = ?"
	topic_row := DbConnection.QueryRow(cmd, topic_title)

	var topic Topic
	err := topic_row.Scan(&topic.Title, &topic.Body, &topic.CreatedDate)
	if err != nil {
		log.Println(err)
	}

	err = templates.ExecuteTemplate(w, "topic.html", topic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./templates/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	t.Execute(w, "aaa")
}

func (topic *Topic) save() error {
	cmd := `INSERT INTO topic (title, body, created_date) VALUES (?,?,?)`
	_, err := DbConnection.Exec(cmd, topic.Title, topic.Body, topic.CreatedDate)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	body := r.FormValue("body")
	value_array := [2]string{title, body}
	for _, v := range value_array {
		if v == "" {
			return
		}
	}
	created_date := time.Now().Format(layout)
	topic := &Topic{Title: title, Body: body, CreatedDate: created_date}
	err := topic.save()

	if err != nil {
		log.Fatalln(err)
	}
	http.Redirect(w, r, "/topic/"+title, http.StatusFound)
}

func commentHandler(w http.ResponseWriter, r *http.Request) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	topic_title := m[2]

	cmd := "SELECT id FROM topics WHERE title = ?"
	topic_id := DbConnection.QueryRow(cmd, topic_title)
	cmd = "INSERT INTO comments (id, body, created_date) VALUES (?, ?, ?)"

	DbConnection.Exec(cmd)

	body := r.FormValue("comment")
	cmd := "INSERT INTO commen"
	DbConnection.Exec()
}

func main() {
	DbConnection, _ = sql.Open("sqlite3", "./internet_forum.sql")
	cmd := `CREATE TABLE IF NOT EXISTS topics (
		title STRING,
		body  STRING,
		created_date STRING
	)`

	_, err := DbConnection.Exec(cmd)

	cmd = `CREATE TABLE IF NOT EXISTS comments (
		id            INT,
		body          STRING,
		created_date  STRING,
		topic_id      INT
	)`
	_, err = DbConnection.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/topic/", topicHandler)
	http.HandleFunc("/new/", newHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/comment/", commentHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
