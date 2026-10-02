package main

import (
	"net/http"
)

// Отдает главную страницу (где находятся табы)
func handlerMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

// Отдает только фрагмент формы для длины
func handlerLength(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "length-form.html")
}

// Отдает только фрагмент формы для веса
func handlerWeight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "weight-form.html")
}

// Отдает только фрагмент формы для температуры
func handlerTemperature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "temperature-form.html")
}

func main() {

	http.HandleFunc("/", handlerMain)
	http.HandleFunc("/api/form/length", handlerLength)
	http.HandleFunc("/api/form/weight", handlerWeight)
	http.HandleFunc("/api/form/temperature", handlerTemperature)
	http.ListenAndServe(":8080", nil)
}
