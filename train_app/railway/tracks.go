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

type TrackSegment struct {
	Id string

	// Resource related
	ReservedBy *Train
	OccupiedBy *Train

	// Track related
	Direction TrackDirection
	Length    units.Meters
}

func (t *TrackSegment) Reserve(train *Train) bool {
	if t.IsAvailable() {
		t.ReservedBy = train
		return true
	}
	return false
}

func (t *TrackSegment) IsOccupied() bool {
	return t.OccupiedBy != nil
}

func (t *TrackSegment) IsAvailable() bool {
	return t.OccupiedBy == nil && t.ReservedBy == nil
}

func (t *TrackSegment) IsReserved() bool {
	return t.OccupiedBy == nil && t.ReservedBy != nil
}

func (t *TrackSegment) Request(train *Train) bool {
	if t.IsAvailable() {
		t.OccupiedBy = train
		return true
	}
	return false
}

func (t *TrackSegment) Release() (bool, error) {
	if t.IsAvailable() {
		fmt.Printf("Cannot release an empty track\n")
		return false, fmt.Errorf("Cannot release an empty track")
	}

	t.OccupiedBy = nil
	return true, nil
}
