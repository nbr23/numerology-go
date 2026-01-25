package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"numerology/solver"
)

type outputFormat string

const (
	formatRaw  outputFormat = "raw"
	formatJSON outputFormat = "json"
	formatText outputFormat = "text"
)

type jsonResponse struct {
	Input      string `json:"input"`
	Target     int    `json:"target"`
	Expression string `json:"expression"`
	Result     int    `json:"result"`
}

func main() {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)

	addr := host + ":" + port
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	target, err := strconv.Atoi(parts[0])
	if err != nil {
		http.Error(w, "invalid target", http.StatusBadRequest)
		return
	}

	format := parseFormat(r.URL.Query().Get("format"))

	var input string
	var userProvidedInput bool
	if len(parts) > 1 && parts[1] != "" {
		input = parts[1]
		userProvidedInput = true
	} else {
		input = time.Now().Format("02012006")
		userProvidedInput = false
	}

	digits := filterDigits(input)
	if len(digits) == 0 {
		http.NotFound(w, r)
		return
	}

	expression, found := solver.Solve(digits, target)
	if !found {
		http.NotFound(w, r)
		return
	}

	writeResponse(w, format, input, target, expression, userProvidedInput)
}

func parseFormat(s string) outputFormat {
	switch strings.ToLower(s) {
	case "json":
		return formatJSON
	case "text":
		return formatText
	default:
		return formatRaw
	}
}

func writeResponse(w http.ResponseWriter, format outputFormat, input string, target int, expression string, userProvidedInput bool) {
	switch format {
	case formatJSON:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jsonResponse{
			Input:      input,
			Target:     target,
			Expression: expression,
			Result:     target,
		})
	case formatText:
		w.Header().Set("Content-Type", "text/plain")
		if userProvidedInput {
			fmt.Fprintf(w, "%s is %s, which equals %d", input, expression, target)
		} else {
			formatted := input[:2] + "/" + input[2:4] + "/" + input[4:]
			fmt.Fprintf(w, "Today is %s and %s = %d", formatted, expression, target)
		}
	default:
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, expression)
	}
}

func filterDigits(s string) []int {
	var digits []int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			digits = append(digits, int(c-'0'))
		}
	}
	return digits
}
