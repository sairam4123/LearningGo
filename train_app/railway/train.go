package railway

type Train struct {
	Name   string
	Number string

	curSchedulePoint int
	schedule         []*SchedulePoint

	FacingToward *TrackPoint

	occupation *OccupationData
}

type SchedulePoint struct {
	TrainNumber string
	StnCode     string
	ArrTime     float64
	DeptTime    float64
	SpPfNo      string
}

func (t *Train) AddSchedule(sp *SchedulePoint) {
	t.schedule = append(t.schedule, sp)
}

func (s *SchedulePoint) ExpDwellTime() float64 {
	return s.DeptTime - s.ArrTime
}

type TrainController struct {
	sim   *Sim
	train *Train
}
