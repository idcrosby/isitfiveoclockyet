package main

import (
	"net/http"
	"time"
	"text/template"
)

func main() {
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8081", nil)	
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	var mainTemplate, _  = template.ParseFiles("drink.html")
	var message string
	if time.Now().Hour() > 16 {
		message = "Yup"
	} else {
		message = "Nope"
	}
	var data = Data{Message: message}
	mainTemplate.Execute(w, data)
}

type Data struct {
	Message string
}