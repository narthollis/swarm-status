package main

import (
	"net/http"
	"swarm-status/web"
	"swarm-status/docker"
)

func main() {
	history := make(docker.HistoryArray, 288)

	// Start the background poll
	go docker.PollServiceStatus(&history)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/history", web.MakeHistoryHandler(&history))
	http.HandleFunc("/current", web.MakeCurrentHandler())

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request){
		http.NotFound(w, req)
	})

	http.HandleFunc("/", web.MakeRootHandler())

	http.ListenAndServe(":8080", nil)
}
