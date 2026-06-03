package railway

type TrainData struct {
	Name   string
	Number string

	curSchedulePoint int
	schedule         []*SchedulePoint

	occupation *OccupationData
}

type SchedulePoint struct {
	TrainNumber string
	StnCode     string
	ArrTime     float64
	DeptTime    float64
	SpPfNo      int
}

func (t *TrainData) AddSchedule(sp *SchedulePoint) {
	t.schedule = append(t.schedule, sp)
}

func (s *SchedulePoint) ExpDwellTime() float64 {
	return s.DeptTime - s.ArrTime
}
