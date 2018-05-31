package web

import (
	"net/http"
	"fmt"
	"swarm-stats/docker"
	"html/template"
	"time"
	"sort"
	"os"
)

type ServiceState struct {
	ID, Name, Status, ClassName string
	Replicas, Running           uint64
}

type ServiceStateGroup struct {
	Status, ClassName string
	Services          []ServiceState
}

type ServiceStateGroupMap map[string]ServiceStateGroup

type ServiceStatusPageData struct {
	Groups    ServiceStateGroupMap
	Timestamp string
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

func getClassNameFromStatusText(status string) string {
	switch status {
	case "OK":
		return "success"
	case "Warning":
		return "warning"
	case "Error":
		return "danger"
	}
	return ""
}

func pickWorstText(a string, b string) string {
	if a == "Error" || b == "Error" {
		return "Error"
	}

	if a == "Warning" || b == "Warning" {
		return "Warning"
	}

	return "OK"
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

		statusText := getStatusText(service.Replicas, service.Running)

		s := append(groups[groupName].Services, ServiceState{
			ID:        service.ID,
			Name:      name,
			Status:    statusText,
			ClassName: getClassNameFromStatusText(statusText),
			Replicas:  service.Replicas,
			Running:   service.Running,
		})

		groupStatusText := pickWorstText(statusText, groups[groupName].Status)

		groups[groupName] = ServiceStateGroup{
			Status:    groupStatusText,
			ClassName: getClassNameFromStatusText(groupStatusText),
			Services:  s,
		}
	}

	for _, val := range groups {
		sort.Slice(val.Services, func(i, j int) bool {
			return val.Services[i].Name < val.Services[j].Name
		})
	}

	err = tmpl.Execute(w, ServiceStatusPageData{
		Groups:    groups,
		Timestamp: time.Now().Format("Mon Jan _2 15:04:05"),
	})

	if err != nil {
		fmt.Println(err)
	}
}
