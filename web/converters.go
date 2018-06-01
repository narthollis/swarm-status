package web

import (
	"swarm-status/docker"
	"sort"
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

func MakeServiceStateGroupMap(settings *EnvSettings, services []docker.ServiceState) ServiceStateGroupMap {
	groups := make(ServiceStateGroupMap)
	for _, service := range services {
		groupName := settings.GetGroupName(service)
		name := settings.GetDisplayName(service)

		status := computeStatus(service.Replicas, service.Running)
		s := append(groups[groupName].Services, ServiceState{
			ID:       service.ID,
			Name:     name,
			Status:   status,
			Replicas: service.Replicas,
			Running:  service.Running,
		})

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

	return groups
}
