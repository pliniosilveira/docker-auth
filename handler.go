package main

import (
	"fmt"
	"net/http"
)

var isBlock bool

func blockHandler(w http.ResponseWriter, r *http.Request) {
	isBlock = true
	statusHandler(w, r)
}

func unblockHandler(w http.ResponseWriter, r *http.Request) {
	isBlock = false
	statusHandler(w, r)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if isBlocked() {
		fmt.Fprintf(w, "New Contaienr create is Blocked")
	} else {
		fmt.Fprintf(w, "New Contaienr create is Unblocked")
	}
}

func isBlocked() bool {
	return isBlock
}
