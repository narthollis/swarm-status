package main

import (
	"net/http"
	"swarm-status/web"
	"swarm-status/docker"
	"os"
	"swarm-status/persist"
)

func main() {
	persistPath := os.Getenv("HISTORY_PERSIST_PATH")
	if persistPath == "" {
		persistPath = "history.json"
	}

	history := make(docker.HistoryArray, 288)
	persist.Load(persistPath, &history)

	// Start the background poll
	go docker.PollServiceStatus(&history, persistPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/history", web.MakeHistoryHandler(&history))
	http.HandleFunc("/current", web.MakeCurrentHandler())

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request){
		http.NotFound(w, req)
	})

	http.HandleFunc("/", web.MakeRootHandler())

	http.ListenAndServe(":8080", nil)
}
