package main

import (
	"fmt"
	"net/http/httputil"
	"net/http"
	"time"
	"bytes"
	"text/template"
	"strconv"
)

func main() {
	http.HandleFunc("/", defaultHandler)
	http.ListenAndServe(":8081", nil)	
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("Request: %s", requestAsString(r))

	if (r.Method == "GET") {
		var mainTemplate, err  = template.ParseFiles("fiveOclock.html")
		if err != nil {
			w.WriteHeader(404)
			return
		}
		mainTemplate.Execute(w, nil)
	} else {
		r.ParseForm()
		var zone = r.FormValue("zone")
		//fmt.Printf("Zone: %s", zone)

		// Get Hour at GMT
		// TODO assuming EST zone
		var hour = time.Now().Hour() + 5
		var minute = time.Now().Minute()

		// TODO catch errors
		var hourDiff, _ = strconv.Atoi(zone[4:6])
		var minDiff, _ = strconv.Atoi(zone[7:9])

		if (zone[:3] == "GMT") {
			if (zone[3:4] == "-") {
				if (minDiff > minute) {
					hourDiff++
				}
				hour = hour - hourDiff % 24
			} else if (zone[3:4] == "+") {
				if (minDiff + minute > 60) {
					hourDiff++
				}
				hour = hour + hourDiff % 24
			} else {
				//error
			}
		} else {
			// Switch on all time zones?
		}

		//fmt.Printf("Adjusted local hour: %d", hour)
		// fmt.Printf("PostForm: %s", r.Form)
		var mainTemplate, err  = template.ParseFiles("fiveOclockWithZone.html")
		if err != nil {
			w.WriteHeader(404)
			return
		}
		var message string
		var offset = 17 - hour
		if offset <= 0 {
			message = "Yup"
		} else {
			message = "It's five o'clock somewhere"
		}
		var data = Data{Message: message}
		// var name, offset = time.Now().Zone()
		// fmt.Printf("Zone is %s with offset of %d", name, offset)
		mainTemplate.Execute(w, data)
	}
}

func getTimeZone() {
	// time.Now().Zone
}

// Create a string which contains all important request data
func requestAsString(request *http.Request) []byte {
	var buffer bytes.Buffer
	buffer.WriteString("\n")
	buffer.WriteString("Current Time: ")
	buffer.WriteString(time.Now().String())
	buffer.WriteString("\n")
	requestBytes, _ := httputil.DumpRequest(request, true)
	buffer.Write(requestBytes)

	return buffer.Bytes()
}

type Data struct {
	Message string
}