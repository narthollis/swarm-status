package docker

import (
	"time"
	"fmt"
)

type History struct {
	Timestamp time.Time
	Status    []ServiceState
}

type HistoryArray []*History
func PollServiceStatus(history *HistoryArray) {
	for {
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

		time.Sleep(time.Minute * 5)
	}
}

