package main

import (
	"log"
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
	poller := &docker.PollServiceStatus{}
	poller.Run(&history, persistPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/history", web.MakeHistoryHandler(&history))
	http.HandleFunc("/current", web.MakeCurrentHandler())

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, req *http.Request){
		http.NotFound(w, req)
	})

	http.HandleFunc("/", web.MakeRootHandler())

	certPath := os.Getenv("TLS_CERT")
	keyPath := os.Getenv("TLS_KEY")

	var err error
	if certPath != "" && keyPath != "" {
		err = http.ListenAndServeTLS(":8080", certPath, keyPath, nil)
	} else {
		err = http.ListenAndServe(":8080", nil)
	}
	if err != nil {
		log.Fatal(err)
	}

	poller.Shutdown()
}
