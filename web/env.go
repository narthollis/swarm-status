package web

import (
	"os"
	"swarm-status/docker"
)

type EnvSettings struct {
	DisplayNameKey, GroupNameKey, MetricPageSource string
}

func GetEnvSettings() EnvSettings {
	displayNameKey := os.Getenv("DISPLAY_NAME_LABEL_KEY")
	if displayNameKey == "" {
		displayNameKey = "com.example.display.name"
	}

	groupNameKey := os.Getenv("DISPLAY_GROUP_LABEL_KEY")
	if groupNameKey == "" {
		groupNameKey = "com.example.display.group"
	}

	metricPageSource := os.Getenv("METRIC_PAGE_SRC")

	return EnvSettings{
		DisplayNameKey:   displayNameKey,
		GroupNameKey:     groupNameKey,
		MetricPageSource: metricPageSource,
	}
}

func (settings EnvSettings) GetDisplayName(service docker.ServiceState) string {
	name := service.Name
	if val, ok := service.Labels[settings.DisplayNameKey]; ok {
		name = val
	}

	return name
}

func (settings EnvSettings) GetGroupName(service docker.ServiceState) string {
	groupName := "Other"
	if val, ok := service.Labels[settings.GroupNameKey]; ok {
		groupName = val
	}

	return groupName
}
