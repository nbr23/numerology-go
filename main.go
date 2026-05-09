package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"numerology/solver"
)

//go:embed static
var staticFS embed.FS

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
	corsOrigin := os.Getenv("CORS_ALLOWED_ORIGIN")

	http.Handle("/api", withCORS(corsOrigin, http.HandlerFunc(handler)))
	http.Handle("/api/", withCORS(corsOrigin, http.HandlerFunc(handler)))

	if os.Getenv("DISABLE_FRONTEND") == "" {
		staticRoot, err := fs.Sub(staticFS, "static")
		if err != nil {
			panic(err)
		}
		http.Handle("/", http.FileServer(http.FS(staticRoot)))
	}

	addr := host + ":" + port
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func withCORS(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if origin != "*" {
				w.Header().Set("Vary", "Origin")
			}
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api"), "/")
	query := r.URL.Query()

	var target int
	var input string
	var userProvidedInput bool

	if path == "" {
		targetStr := query.Get("target")
		if targetStr == "" {
			http.NotFound(w, r)
			return
		}
		t, err := strconv.Atoi(targetStr)
		if err != nil {
			http.Error(w, "invalid target", http.StatusBadRequest)
			return
		}
		target = t
		if d := query.Get("digits"); d != "" {
			input = d
			userProvidedInput = true
		} else {
			input = time.Now().Format("02012006")
		}
	} else {
		parts := strings.Split(path, "/")
		t, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "invalid target", http.StatusBadRequest)
			return
		}
		target = t
		if len(parts) > 1 && parts[1] != "" {
			input = parts[1]
			userProvidedInput = true
		} else {
			input = time.Now().Format("02012006")
		}
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

	writeResponse(w, parseFormat(query.Get("format")), input, target, expression, userProvidedInput)
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
	if userProvidedInput {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", secondsUntilMidnight()))
	}
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

func secondsUntilMidnight() int {
	now := time.Now()
	tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return int(tomorrow.Sub(now).Seconds())
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
