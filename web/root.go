package web

import (
	"net/http"
	"fmt"
	"swarm-status/docker"
	"html/template"
	"time"
	"sort"
	"os"
)

type ServiceState struct {
	ID, Name, ClassName string
	Status              Status
	Replicas, Running   uint64
}

type ServiceStateGroup struct {
	Status   Status
	Services []ServiceState
}

type ServiceStateGroupMap map[string]ServiceStateGroup

type ServiceStatusPageData struct {
	Groups        ServiceStateGroupMap
	Timestamp     string
	OverallStatus Status
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

func RootHandler(w http.ResponseWriter, r *http.Request) {
	displayNameKey := os.Getenv("DISPLAY_NAME_LABEL_KEY")
	if displayNameKey == "" {
		displayNameKey = "com.example.display.name"
	}

	groupNameKey := os.Getenv("DISPLAY_GROUP_LABEL_KEY")
	if groupNameKey == "" {
		groupNameKey = "com.example.display.group"
	}

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

	overallStatus := Operational
	groups := make(ServiceStateGroupMap)
	for _, service := range services {
		name := service.Name
		if val, ok := service.Labels[displayNameKey]; ok {
			name = val
		}

		groupName := "Other"
		if val, ok := service.Labels[groupNameKey]; ok {
			groupName = val
		}

		status := computeStatus(service.Replicas, service.Running)
		s := append(groups[groupName].Services, ServiceState{
			ID:       service.ID,
			Name:     name,
			Status:   status,
			Replicas: service.Replicas,
			Running:  service.Running,
		})

		overallStatus = status.PickWorst(overallStatus)

		groups[groupName] = ServiceStateGroup{
			Status:   status.PickWorst(groups[groupName].Status),
			Services: s,
		}
	}

	for _, val := range groups {
		sort.Slice(val.Services, func(i, j int) bool {
			return val.Services[i].Name < val.Services[j].Name
		})
	}

	err = tmpl.Execute(w, ServiceStatusPageData{
		Groups:        groups,
		Timestamp:     time.Now().Format("Mon Jan _2 15:04:05"),
		OverallStatus: overallStatus,
	})

	if err != nil {
		fmt.Println(err)
	}
}
