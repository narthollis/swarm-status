package web

import "encoding/json"

type Status int

const (
	Stopped		Status = iota
	Operational
	Unhealthy
	Critical
)

func (status Status) String() string {
	names := [...]string{
		"Stopped",
		"Operational",
		"Unhealthy",
		"Critical",
	}
	if status < Stopped || status > Critical {
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
		"secondary",
		"success",
		"warning",
		"danger",
	}
	if status < Stopped || status > Critical {
		return ""
	}

	return classNames[status]
}

func (status Status) Icon() string {
	icons := [...]string{
		"",
		"fa-check-circle",
		"fa-bell",
		"fa-exclamation-triangle",
	}
	if status < Stopped || status > Critical {
		return ""
	}

	return icons[status]
}

func (status Status) Overview() string {
	classNames := [...]string{
		"Stopped",
		"All Systems Operational",
		"System Unhealthy",
		"System Critical",
	}
	if status < Stopped || status > Critical {
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
