package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Question struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

var questions []Question

func loadQuestions() error {
	file, err := os.ReadFile("questions.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &questions)
	if err != nil {
		return err
	}
	return nil
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
	err := loadQuestions()
	if err != nil {
		log.Fatal("Error loading questions:", err)
	}

	http.HandleFunc("/", quizHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/result", resultHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
