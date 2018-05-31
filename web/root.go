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
	Replicas, Running uint64
}

type ServiceStateGroup map[string][]ServiceState

type ServiceStatusPageData struct {
	Groups ServiceStateGroup
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

	fmt.Println(services)

	groups := make(ServiceStateGroup)
	for _, service := range services {
		name := service.Name
		if val, ok := service.Labels[displayNameKey]; ok {
			name = val
		}

		groupName := "Other"
		if val, ok := service.Labels[groupNameKey]; ok {
			groupName = val
		}

		groups[groupName] = append(groups[groupName], ServiceState{
			ID:        service.ID,
			Name:      name,
			Status:    getStatusText(service.Replicas, service.Running),
			ClassName: getClassName(service.Replicas, service.Running),
			Replicas:  service.Replicas,
			Running:   service.Running,
		})
	}

	for _, val := range groups {
		sort.Slice(val, func(i, j int) bool {
			return val[i].Name < val[j].Name
		})
	}

	err = tmpl.Execute(w, ServiceStatusPageData{
		Groups: groups,
		Timestamp: time.Now().Format("Mon Jan _2 15:04:05"),
	})

	if err != nil {
		fmt.Println(err)
	}
}
