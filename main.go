package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Question struct {
	Text   string
	Answer string
}

var questions = []Question{
	{"Cate culori are Drapelul?", "3"},
	{"Cat e 2 + 2?", "4"},
	{"Care e capitala Franței", "Paris"},
	{"Cât e radical de 64", "8"},
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/quiz.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, questions)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	score := 0
	for i, question := range questions {
		answer := r.FormValue("q" + strconv.Itoa(i))
		if answer == question.Answer {
			score++
		}
	}
	http.Redirect(w, r, "/result?score="+strconv.Itoa(score), http.StatusSeeOther)
}

func resultHandler(w http.ResponseWriter, r *http.Request) {
	score := r.URL.Query().Get("score")
	t, _ := template.ParseFiles("static/result.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, score)
}

func main() {
	http.HandleFunc("/", quizHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/result", resultHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
