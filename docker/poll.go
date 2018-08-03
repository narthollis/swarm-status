package docker

import (
	"time"
	"fmt"
	"swarm-status/persist"
)

type History struct {
	Timestamp time.Time
	Status    []ServiceState
}

type HistoryArray []*History
func pollServiceStatus(history *HistoryArray, persistPath string, quit chan bool) {
	for {
		select {
		case <- quit:
			return
		default:
			next, err := ReadServiceList()

			if err != nil {
				fmt.Println(err)
			} else {
				for i := len(*history) - 2; i >= 0; i-- {
					(*history)[i+1] = (*history)[i]
				}

				(*history)[0] = &History{
					Status: next,
					Timestamp: time.Now().UTC(),
				}
			}

			persist.Save(persistPath, *history)

			time.Sleep(time.Minute * 5)
		}
	}
}
