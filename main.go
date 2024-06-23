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

type Test struct {
	Name      string     `json:"name"`
	Questions []Question `json:"questions"`
}

type Tests struct {
	Tests []Test `json:"tests"`
}

var tests Tests

func loadTests() error {
	file, err := os.ReadFile("questions.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &tests)
	if err != nil {
		return err
	}
	return nil
}

func testsHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/tests.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, tests)
}

func quizHandler(w http.ResponseWriter, r *http.Request) {
	testIndex, err := strconv.Atoi(r.URL.Query().Get("test"))
	if err != nil || testIndex < 0 || testIndex >= len(tests.Tests) {
		http.Error(w, "Invalid test index", http.StatusBadRequest)
		return
	}

	type QuizData struct {
		Index     int
		Questions []Question
	}

	quizData := QuizData{
		Index:     testIndex,
		Questions: tests.Tests[testIndex].Questions,
	}

	t, err := template.ParseFiles("static/quiz.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, quizData)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	testIndex, err := strconv.Atoi(r.FormValue("test"))
	if err != nil || testIndex < 0 || testIndex >= len(tests.Tests) {
		http.Error(w, "Invalid test index", http.StatusBadRequest)
		return
	}

	test := tests.Tests[testIndex]
	score := 0
	for i, question := range test.Questions {
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
	err := loadTests()
	if err != nil {
		log.Fatal("Error loading tests:", err)
	}

	http.HandleFunc("/", testsHandler)
	http.HandleFunc("/quiz", quizHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/result", resultHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
