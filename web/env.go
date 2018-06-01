package web

import (
	"os"
	"swarm-status/docker"
)

type EnvSettings struct {
	DisplayNameKey, GroupNameKey string
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

	return EnvSettings{
		DisplayNameKey: displayNameKey,
		GroupNameKey:   groupNameKey,
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
