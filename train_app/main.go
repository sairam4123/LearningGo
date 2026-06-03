package main

import (
	"trainapp/railway"
	"trainapp/units"
)

func main() {
	sim := railway.Sim{}

	world := railway.World{}
	world.Init()

	pdkt := railway.Station{
		Code: "PDKT",
		Name: "Pudukkottai",
	}

	pdkt.Init()

	kkdi := railway.Station{
		Code: "KKDI",
		Name: "Karaikudi",
	}

	kkdi.Init()

	trPdkt_Kkdi := railway.TrackData{
		Id:        "PDKT_KKDI",
		Direction: railway.Bidir,
	}

	bsPdkt_Kkdi := railway.BlockSection{}
	bsPdkt_Kkdi.Init(&pdkt, &kkdi)

	bsPdkt_Kkdi.AddTrack(&trPdkt_Kkdi)

	trKkdi0 := railway.TrackData{
		Id:        "kkdi0",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}
	trKkdi1 := railway.TrackData{
		Id:        "kkdi1",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}

	trPdkt0 := railway.TrackData{
		Id:        "pdkt0",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}
	trPdkt1 := railway.TrackData{
		Id:        "pdkt1",
		Direction: railway.Bidir,
		Length:    units.M(500),
	}

	world.AddStation(&kkdi)
	world.AddStation(&pdkt)

	world.AddBlockSection(&bsPdkt_Kkdi)

	kkdi0 := railway.Platform{
		Id:     trKkdi0.Id,
		Track:  &trKkdi0,
		Length: units.M(500),
	}
	kkdi1 := railway.Platform{
		Id:     trKkdi1.Id,
		Track:  &trKkdi1,
		Length: units.M(500),
	}
	pdkt0 := railway.Platform{
		Id:     trPdkt0.Id,
		Track:  &trPdkt0,
		Length: units.M(500),
	}
	pdkt1 := railway.Platform{
		Id:     trPdkt1.Id,
		Track:  &trPdkt1,
		Length: units.M(500),
	}

	kkdi.AddPlatform(&kkdi0)
	kkdi.AddPlatform(&kkdi1)
	pdkt.AddPlatform(&pdkt0)
	pdkt.AddPlatform(&pdkt1)

	_12606 := railway.TrainData{
		Name:   "Pallavan",
		Number: "12606",
	}

	_12605 := railway.TrainData{
		Name:   "RMM Express",
		Number: "12605",
	}

	world.AddTrain(&_12605)
	world.AddTrain(&_12606)

	_12606.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12606.Number,
		StnCode:     kkdi.Code,
		ArrTime:     120,
		DeptTime:    125,
		SpPfNo:      0,
	})

	_12606.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12606.Number,
		StnCode:     pdkt.Code,
		ArrTime:     145,
		DeptTime:    150,
		SpPfNo:      0,
	})

	_12605.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12605.Number,
		StnCode:     pdkt.Code,
		ArrTime:     130,
		DeptTime:    140,
		SpPfNo:      0,
	})

	_12605.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12605.Number,
		StnCode:     kkdi.Code,
		ArrTime:     320,
		DeptTime:    340,
		SpPfNo:      0,
	})

	sim.SetWorld(&world)
	sim.InitWorld()

	sim.Run()
}
