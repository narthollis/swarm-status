package web

import (
	"net/http"
	"fmt"
	"swarm-status/docker"
	"html/template"
	"time"
)

type ServiceStatusPageData struct {
	Groups               ServiceStateGroupMap
	Timestamp            string
	OverallStatus        Status
	OverallStatusVerbose string
}

func computeStatus(replicas uint64, running uint64) Status {
	if replicas == running {
		return Operational
	} else if running > 0 {
		return Unhealthy
	} else {
		return Critical
	}
}

func MakeRootHandler() http.HandlerFunc {
	settings := GetEnvSettings()

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/serviceStatus.html")

		if err != nil {
			panic(err)
		}

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

		err = tmpl.Execute(w, ServiceStatusPageData{
			Groups:               groups,
			Timestamp:            time.Now().Format("Mon Jan _2 15:04:05"),
			OverallStatus:        overallStatus,
			OverallStatusVerbose: overallStatus.Overview(),
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}
