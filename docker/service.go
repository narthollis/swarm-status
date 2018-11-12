package docker

type PollServiceStatus struct {
	quit chan bool
}

func (m *PollServiceStatus) Run(history *HistoryArray, persistPath string) {
	m.quit = make(chan bool)

	go pollServiceStatus(history, persistPath, m.quit)
}

func (m *PollServiceStatus) Shutdown() {
	m.quit <- true
}
