package web

import (
	"net/http"
	"fmt"
	"swarm-stats/docker"
	"html/template"
	"time"
)

type ServiceState struct {
	ID, Name, Status, ClassName string
	Replicas, Running uint64
}

type ServiceStatusPageData struct {
	Services []ServiceState
	Timestamp time.Time
}

func getStatusText(replicas uint64, running uint64) string {
	if replicas == running {
		return "OK"
	} else if running > 0 {
		return "Warning"
	} else {
		return "Error"
	}
}

func getClassName(replicas uint64, running uint64) string {
	if replicas == running {
		return "success"
	} else if running > 0 {
		return "warning"
	} else {
		return "danger"
	}
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/serviceStatus.html")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	services, err := docker.ReadServiceList()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	var states []ServiceState
	for _, service := range services {
		states = append(states, ServiceState{
			ID:        service.ID,
			Name:      service.Name,
			Status:    getStatusText(service.Replicas, service.Running),
			ClassName: getClassName(service.Replicas, service.Running),
			Replicas:  service.Replicas,
			Running:   service.Running,
		})
	}

	fmt.Println(states)

	tmpl.Execute(w, ServiceStatusPageData{Services: states, Timestamp: time.Now().UTC() })
}
