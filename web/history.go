package web

import (
	"swarm-status/docker"
	"net/http"
	"encoding/json"
	"time"
)

type FormattedHistory struct {
	Timestamp string
	Groups    ServiceStateGroupMap
}

func MakeHistoryHandler(history *docker.HistoryArray) http.HandlerFunc {
	settings := GetEnvSettings()

	return func(w http.ResponseWriter, r *http.Request) {
		hist := (*history)[:]

		var output []FormattedHistory
		for _, h := range hist {
			if h != nil {
				output = append(output, FormattedHistory{
					Timestamp: h.Timestamp.UTC().Format(time.RFC3339),
					Groups:    MakeServiceStateGroupMap(&settings, h.Status),
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(output)
	}
}
