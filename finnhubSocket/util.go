package main

import "net/http"

// return true to allow connection, false to dismiss
func checkOrigin(r *http.Request) bool {
	return true
	/*
		origin := r.Header.Get("Origin")
		switch origin {
		case "https://www.socketchat.app":
			return true
		default:
			return false
		}
	*/
}
