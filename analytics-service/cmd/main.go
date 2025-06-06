package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, "Analytics-service")
}

func main() {
	http.HandleFunc("/analytics", helloHandler)

	fmt.Println("Started server at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("error running server:", err)
	}
}
