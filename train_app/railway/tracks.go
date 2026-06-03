package railway

import (
	"fmt"
	"trainapp/units"
)

type TrackDirection string

const (
	Up    TrackDirection = "UP"
	Down  TrackDirection = "DOWN"
	Bidir TrackDirection = "BIDIR"
)

type TrackData struct {
	Id string

	// Resource related
	ReservedBy *TrainData
	OccupiedBy *TrainData

	Direction TrackDirection
	Length    units.Meters
}

func (t *TrackData) IsAvailable() bool {
	return t.OccupiedBy == nil && t.ReservedBy == nil
}

func (t *TrackData) IsReserved() bool {
	return t.OccupiedBy == nil && t.ReservedBy != nil
}

func (t *TrackData) Request(train *TrainData) bool {
	if t.IsAvailable() {
		t.OccupiedBy = train
		return true
	} else {
		return false
	}
}

func (t *TrackData) Release() (bool, error) {
	if t.IsAvailable() {
		fmt.Printf("Cannot release an empty track\n")
		return false, fmt.Errorf("Cannot release an empty track")
	}

	t.OccupiedBy = nil
	return true, nil
}
