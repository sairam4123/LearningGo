package railway

import "trainapp/units"

type Station struct {
	Code string
	Name string

	Platforms []*Platform
}

type Platform struct {
	Id string

	Length units.Meters

	Track *TrackData
}

func (stn *Station) Init() {
	stn.Platforms = make([]*Platform, 0)
}

func (stn *Station) AddPlatform(pfData *Platform) {
	stn.Platforms = append(stn.Platforms, pfData)
}
