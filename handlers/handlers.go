package handlers

import (
	ascii "ascii-art-web/program"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Data struct {
	Text   string
	Output string
}

var fonts = map[string]bool{
	"standard.txt":   true,
	"shadow.txt":     true,
	"thinkertoy.txt": true,
}

var (
	templates *template.Template
	err       error
)

func init() {
	templates, err = template.ParseGlob("./web/templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errors(w, "Ooops! Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, "Ooops! Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data := Data{
		Text:   "",
		Output: "",
	}
	err := templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		Errors(w, "Ooops! Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func Output(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ascii-art" {
		Errors(w, "Ooops! Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodPost {
		Errors(w, "Ooops! Not Found", http.StatusNotFound)
		return
	}
	text := r.FormValue("text")
	font := r.FormValue("banners")
	for _, char := range text {
		if char < 9 || char > 127 {
			Errors(w, "Ooops! Bad Request", http.StatusBadRequest)
			return
		}
	}
	if !fonts[font] || text == "" {
		Errors(w, "Ooops! Bad Request", http.StatusBadRequest)
		return
	}
	text = strings.ReplaceAll(text, "\r", "")
	asciiart, err := ascii.Convert(text, font)
	if err != nil {
		log.Println("ASCII conversion error:", err)
		Errors(w, "Ooops! Internal Server Error", http.StatusInternalServerError)
		return
	}
	data := Data{
		Text:   text,
		Output: asciiart,
	}
	err = templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		Errors(w, "Ooops! Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func Errors(w http.ResponseWriter, msgs string, status int) {
	Error := struct {
		Msg  string
		Code int
	}{
		msgs,
		status,
	}
	err := templates.ExecuteTemplate(w, "error.html", Error)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
