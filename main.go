package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func getCurrentTime(w http.ResponseWriter, r *http.Request) {

	utc := time.Now().UTC()
	keys, ok := r.URL.Query()["tz"]
	timeMap := map[string]string{}

	// if tz is not provided
	if !ok {
		w.Header().Add("Content-Type", "application/json")
		timeMap["current_time"] = utc.String()
		jsonStr, _ := json.Marshal(timeMap)
		fmt.Fprint(w, string(jsonStr))
		return
	}

	tzs := strings.Split(keys[0], ",")

	for i, tz := range tzs {
		tzs[i] = strings.TrimSpace(tz)
	}

	for _, tz := range tzs {
		loc, err := time.LoadLocation(tz)

		if err == nil {
			timeMap[tz] = utc.In(loc).String()
		}

	}

	// no valid timezones provided
	if len(timeMap) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "Invalid timezone")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonStr, _ := json.Marshal(timeMap)
	fmt.Fprint(w, string(jsonStr))

}

func main() {
	http.HandleFunc("/api/time", getCurrentTime)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
