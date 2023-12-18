package main

import "net/http"

func errHandler(w http.ResponseWriter, r *http.Request) {
	respondWtihError(w, 400, "Something went wrong")
}
