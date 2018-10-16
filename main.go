package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	http.Handle("/css/", http.FileServer(http.Dir("static/")))
	http.Handle("/js/", http.FileServer(http.Dir("static/")))
	http.Handle("/font-awesome/", http.FileServer(http.Dir("static/")))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
