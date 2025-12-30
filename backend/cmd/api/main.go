package main

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

var allowedOrigins = []string{
	"http://localhost:5173",
	"https://bookmarks-web-qbt.pages.dev",
}

func isOriginAllowed(origin string) bool {
	// Check exact matches
	if slices.Contains(allowedOrigins, origin) {
		return true
	}

	// Check if it's a Cloudflare Pages preview URL
	// Pattern: https://<hash>.bookmarks-web-qbt.pages.dev
	if strings.HasSuffix(origin, ".bookmarks-web-qbt.pages.dev") && strings.HasPrefix(origin, "https://") {
		return true
	}

	return false
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Add("Vary", "Origin")
		} else {
			fmt.Printf("Received request from non-allowedOrigin: %v\n", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		fmt.Fprintf(w, "Hello!")
	})

	fmt.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
