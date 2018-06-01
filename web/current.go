package web

import (
	"net/http"
	"swarm-status/docker"
	"fmt"
	"encoding/json"
	"time"
)

func MakeCurrentHandler() http.HandlerFunc {
	settings := GetEnvSettings()

	return func(w http.ResponseWriter, r *http.Request) {
		services, err := docker.ReadServiceList()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		groups := MakeServiceStateGroupMap(&settings, services)
		overallStatus := Operational
		for _, g := range groups {
			for _, s := range g.Services {
				overallStatus = s.Status.PickWorst(overallStatus)
			}
		}

		json.NewEncoder(w).Encode(ServiceStatusPageData{
			Groups:               groups,
			Timestamp:            time.Now().UTC().Format(time.RFC3339),
			OverallStatus:        overallStatus,
			OverallStatusVerbose: overallStatus.Overview(),
		})
	}
}
