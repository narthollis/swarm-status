package docker

import (
	"golang.org/x/net/context"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"os"
	"time"
)

type ServiceState struct {
	ID, Name          string
	Replicas, Running uint64
	Labels            map[string]string
	UpdatedAt         time.Time
}

func setApiVersion() {
	if os.Getenv("DOCKER_API_VERSION") == "" {
		os.Setenv("DOCKER_API_VERSION", "1.35")
	}
}

func ReadServiceList() ([]ServiceState, error) {
	setApiVersion()

	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})

	if err != nil {
		return nil, err
	}

	var taskFilters = filters.NewArgs();
	taskFilters.Add("desired-state", "running")

	tasks, err := cli.TaskList(ctx, types.TaskListOptions{
		Filters: taskFilters,
	})

	if err != nil {
		return nil, err
	}

	var servicesState []ServiceState

	for _, service := range services {
		if service.Spec.Mode.Replicated != nil && service.Spec.TaskTemplate.ContainerSpec != nil {
			var runningTasks uint64 = 0
			var t time.Time

			for _, task := range tasks {
				if task.ServiceID == service.ID && task.Status.State == "running" {
					runningTasks += 1
					if t.Before(task.Status.Timestamp) {
						t = task.Status.Timestamp
					}
				}
			}

			var serviceState = ServiceState{
				ID:        service.ID,
				Name:      service.Spec.Name,
				Labels:    service.Spec.Labels,
				Replicas:  *service.Spec.Mode.Replicated.Replicas,
				Running:   runningTasks,
				UpdatedAt: t,
			}

			servicesState = append(servicesState, serviceState)
		}
	}

	return servicesState, nil
}
