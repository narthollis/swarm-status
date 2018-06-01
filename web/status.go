package web

import "encoding/json"

type Status int

const (
	Operational Status = iota
	Unhealthy
	Critical
)

func (status Status) String() string {
	names := [...]string{
		"Operational",
		"Unhealthy",
		"Critical",
	}
	if status < Operational || status > Critical {
		return "Unknown"
	}

	return names[status]
}

func (status Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Value, ClassName, Icon string
	}{
		Value:     status.String(),
		ClassName: status.ClassName(),
		Icon:      status.Icon(),
	})
}

func (status Status) ClassName() string {
	classNames := [...]string{
		"success",
		"warning",
		"danger",
	}
	if status < Operational || status > Critical {
		return ""
	}

	return classNames[status]
}

func (status Status) Icon() string {
	icons := [...]string{
		"fa-check-circle",
		"fa-bell",
		"fa-exclamation-triangle",
	}
	if status < Operational || status > Critical {
		return ""
	}

	return icons[status]
}

func (status Status) Overview() string {
	classNames := [...]string{
		"All Systems Operational",
		"System Unhealthy",
		"System Critical",
	}
	if status < Operational || status > Critical {
		return ""
	}

	return classNames[status]
}

func (status Status) PickWorst(other Status) Status {
	if status > other {
		return status
	} else if status < other {
		return other
	}

	return status
}
