package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type ConvertRequest struct {
	Value float64 `json:"value"`
	From  string  `json:"from"`
	To    string  `json:"to"`
}

type ConvertResponse struct {
	Result string `json:"result"`
}

func handlerConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var req ConvertRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	finalResult, err := convertLength(req.Value, req.From, req.To)

	strResponse := formatFloat(finalResult)
	strRequest := formatFloat(req.Value)

	responseText := fmt.Sprintf("%s %s = %s %s", strRequest, req.From, strResponse, req.To)

	w.Header().Set("Content-Type", "application/json")

	resp := ConvertResponse{Result: responseText}
	json.NewEncoder(w).Encode(resp)
}

func formatFloat(val float64) string {
	str := fmt.Sprintf("%.8f", val)
	str = strings.TrimRight(str, "0")
	return strings.TrimRight(str, ".")
}

// Коэффициенты перевода (сколько метров в одной единице)
var ratesLength = map[string]float64{
	// Метрическая система
	"mm":         0.001,
	"millimeter": 0.001,
	"cm":         0.01,
	"centimeter": 0.01,
	"m":          1.0,
	"meter":      1.0,
	"km":         1000.0,
	"kilometer":  1000.0,
	// Имперская система
	"in":   0.0254,
	"inch": 0.0254,
	"ft":   0.3048,
	"foot": 0.3048,
	"yd":   0.9144,
	"yard": 0.9144,
	"mi":   1609.344,
	"mile": 1609.344,
}

func convertLength(value float64, from, to string) (float64, error) {
	fromKey := strings.ToLower(strings.TrimSpace(from))
	toKey := strings.ToLower(strings.TrimSpace(to))

	fromRate, ok1 := ratesLength[fromKey]
	toRate, ok2 := ratesLength[toKey]

	if !ok1 || !ok2 {
		return 0, fmt.Errorf("неверная единица измерения: %s или %s", from, to)
	}
	return (value * fromRate) / toRate, nil
}

func main() {

	http.HandleFunc("/", handlerMain)
	http.HandleFunc("/api/form/length", handlerLength)
	http.HandleFunc("/api/form/weight", handlerWeight)
	http.HandleFunc("/api/form/temperature", handlerTemperature)
	http.HandleFunc("/api/convert", handlerConvert)
	http.ListenAndServe(":8080", nil)
}
