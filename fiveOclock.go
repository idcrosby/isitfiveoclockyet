package main

import (
	"fmt"
	"net/http/httputil"
	"net/http"
	"time"
	"bytes"
	"text/template"
	"strconv"
	"math"
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
		fmt.Printf("Zone: %s \n", zone)

		var hour = getLocalHour(zone)

		fmt.Printf("Adjusted local hour: %d \n", hour)
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
		mainTemplate.Execute(w, data)
	}
}

func getLocalHour(zone string) int {
	// Get Hour at GMT
	var hour = time.Now().UTC().Hour()
	var minute = time.Now().Minute()

	// TODO catch errors
	var hourDiff, err = strconv.Atoi(zone[3:6])
	if (err != nil) {
		return hour
	}
	var minDiff, err1 = strconv.Atoi(zone[7:9])
	if (err1 != nil) {
		return hour
	}

	if (zone[:3] == "GMT") {
		// TODO Make this smarter
		if (zone[3:4] == "-") {
			if (minDiff > minute) {
				hourDiff++
			}
		} else if (zone[3:4] == "+") {
			if (minDiff + minute > 60) {
				hourDiff++
			}
		} else {
			//Unknown format
		}
		hour = int(math.Mod(float64(hour + hourDiff) + 24.0, 24.0))
	} else {
		// Switch on all time zones?
	}

	return hour
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