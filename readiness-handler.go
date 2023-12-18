package main

import "net/http"

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	respondWtihJSON(w, 200, struct{}{})
}
