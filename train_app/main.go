package main

import (
	"trainapp/railway"
	"trainapp/units"
)

func mainOld() {
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

	trPdkt_Kkdi := railway.TrackSegment{
		Id:        "PDKT_KKDI",
		Direction: railway.Bidir,
		Length:    units.KM(20),
	}

	bsPdkt_Kkdi := railway.BlockSection{
		Id: trPdkt_Kkdi.Id,
	}
	bsPdkt_Kkdi.Init(&pdkt, &kkdi)

	bsPdkt_Kkdi.AddTrack(&trPdkt_Kkdi)

	trKkdi0 := railway.TrackSegment{
		Id:        "kkdi0",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}
	trKkdi1 := railway.TrackSegment{
		Id:        "kkdi1",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}

	trPdkt0 := railway.TrackSegment{
		Id:        "pdkt0",
		Direction: railway.Bidir,
		Length:    units.M(700),
	}
	trPdkt1 := railway.TrackSegment{
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

	tpj := railway.Station{
		Code: "TPJ",
		Name: "Tiruchy",
	}

	tpj.Init()
	world.AddStation(&tpj)

	trTpj0 := railway.TrackSegment{
		Id:        "tpj0",
		Direction: railway.Bidir,
		Length:    units.M(500),
	}
	trTpj1 := railway.TrackSegment{
		Id:        "tpj1",
		Direction: railway.Bidir,
		Length:    units.M(500),
	}

	tpj0 := railway.Platform{
		Id:     trTpj0.Id,
		Length: units.M(400),
		Track:  &trTpj0,
	}
	tpj1 := railway.Platform{
		Id:     trTpj1.Id,
		Length: units.M(400),
		Track:  &trTpj1,
	}

	trTpj_Pdkt := railway.TrackSegment{
		Id:        "TPJ_PDKT",
		Direction: railway.Bidir,
		Length:    units.KM(40),
	}

	bsTpj_Pdkt := railway.BlockSection{
		Id: trTpj_Pdkt.Id,
	}
	bsTpj_Pdkt.Init(&tpj, &pdkt)
	bsTpj_Pdkt.AddTrack(&trTpj_Pdkt)

	world.AddBlockSection(&bsTpj_Pdkt)

	kkdi.AddPlatform(&kkdi0)
	kkdi.AddPlatform(&kkdi1)

	pdkt.AddPlatform(&pdkt0)
	pdkt.AddPlatform(&pdkt1)

	tpj.AddPlatform(&tpj0)
	tpj.AddPlatform(&tpj1)

	_12606 := railway.Train{
		Name:   "Pallavan",
		Number: "12606",
	}

	_12605 := railway.Train{
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

	_12606.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12605.Number,
		StnCode:     tpj.Code,
		ArrTime:     165,
		DeptTime:    180,
		SpPfNo:      0,
	})

	_12605.AddSchedule(&railway.SchedulePoint{
		TrainNumber: _12605.Number,
		StnCode:     tpj.Code,
		ArrTime:     100,
		DeptTime:    110,
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

func main() {

	PLATFORM_TRACK_LENGTH := 800.0
	PLATFORM_LENGTH := 600.0

	sim := railway.Sim{}

	world := &railway.World{}
	world.Init()
	sim.SetWorld(world)

	tpj := world.NewStation("TPJ", "Tiruchy Jn")
	tpj.Init()
	pdkt := world.NewStation("PDKT", "Pudukkottai")
	tpj.Init()

	tpjPf1S := world.NewTrackPoint("tpjPf1S").SetDeadEnd(true).SetSimBoundary(true)
	tpjPf1E := world.NewTrackPoint("tpjPf1E")

	tpjPf2S := world.NewTrackPoint("tpjPf2S").SetDeadEnd(true).SetSimBoundary(true)
	tpjPf2E := world.NewTrackPoint("tpjPf2E")

	tpjPf3S := world.NewTrackPoint("tpjPf3S").SetDeadEnd(true).SetSimBoundary(true)
	tpjPf3E := world.NewTrackPoint("tpjPf3E")

	tpjPf1 := &railway.TrackSegment{
		Id:     "tpjPf1",
		Length: units.M(PLATFORM_TRACK_LENGTH),
	}

	tpjPf2 := &railway.TrackSegment{
		Id:     "tpjPf2",
		Length: units.M(PLATFORM_TRACK_LENGTH),
	}

	tpjPf3 := &railway.TrackSegment{
		Id:     "tpjPf3",
		Length: units.M(PLATFORM_TRACK_LENGTH),
	}

	world.TrackGraph.AddTrack(tpjPf1S, tpjPf1E, tpjPf1)
	world.TrackGraph.AddTrack(tpjPf2S, tpjPf2E, tpjPf2)
	world.TrackGraph.AddTrack(tpjPf3S, tpjPf3E, tpjPf3)

	tpj.AddPlatform(&railway.Platform{
		Track:  tpjPf1,
		Length: units.M(PLATFORM_LENGTH),
	})

	tpj.AddPlatform(&railway.Platform{
		Track:  tpjPf2,
		Length: units.M(PLATFORM_LENGTH),
	})

	tpj.AddPlatform(&railway.Platform{
		Track:  tpjPf3,
		Length: units.M(PLATFORM_LENGTH),
	})

	pdktPf1S := world.NewTrackPoint("pdktPf1S")
	pdktPf1E := world.NewTrackPoint("pdktPf1E").SetDeadEnd(true).SetSimBoundary(true)

	pdktPf2S := world.NewTrackPoint("pdktPf2S")
	pdktPf2E := world.NewTrackPoint("pdktPf2E").SetDeadEnd(true).SetSimBoundary(true)

	pdktPf1 := railway.TrackSegment{
		Id:     "pdktPf1",
		Length: units.M(PLATFORM_TRACK_LENGTH),
	}

	pdktPf2 := railway.TrackSegment{
		Id:     "pdktPf2",
		Length: units.M(PLATFORM_TRACK_LENGTH),
	}

	world.TrackGraph.AddTrack(pdktPf1S, pdktPf1E, &pdktPf1)
	world.TrackGraph.AddTrack(pdktPf2S, pdktPf2E, &pdktPf2)

	tpjSw1 := world.NewTrackPoint("tpjSw1")
	tpjPf1ESw1 := railway.TrackSegment{
		Id:     "tpjPf1ESw1",
		Length: units.M(100),
	}
	tpjPf2ESw1 := railway.TrackSegment{
		Id:     "tpjPf2ESw1",
		Length: units.M(100),
	}

	tpjPf3ESw1 := railway.TrackSegment{
		Id:     "tpjPf3ESw1",
		Length: units.M(100),
	}

	world.TrackGraph.AddTrack(tpjPf1E, tpjSw1, &tpjPf1ESw1)
	world.TrackGraph.AddTrack(tpjPf2E, tpjSw1, &tpjPf2ESw1)
	world.TrackGraph.AddTrack(tpjPf3E, tpjSw1, &tpjPf3ESw1)

	pdktSw1 := world.NewTrackPoint("pdktSw1")

	pdktPf1SSw1 := railway.TrackSegment{
		Id:     "pdktPf1SSw1",
		Length: units.M(100),
	}

	pdktPf2SSw1 := railway.TrackSegment{
		Id:     "pdktPf2SSw1",
		Length: units.M(100),
	}

	world.TrackGraph.AddTrack(pdktPf1S, pdktSw1, &pdktPf1SSw1)
	world.TrackGraph.AddTrack(pdktPf2S, pdktSw1, &pdktPf2SSw1)

	bsecTpjPdkt := world.NewBlockSection("bsecTpjPdkt")
	bsecTpjPdkt.Init(tpj, pdkt)
	bsTpjPdkt0 := railway.TrackSegment{
		Id:     "bsTpjPdkt0",
		Length: units.KM(5),
	}
	bsTpjPdkt1 := railway.TrackSegment{
		Id:     "bsTpjPdkt1",
		Length: units.KM(16),
	}
	bsTpjPdkt2 := railway.TrackSegment{
		Id:     "bsTpjPdkt2",
		Length: units.KM(5),
	}
	krurCp1 := world.NewTrackPoint("krurCp1")
	tpjCp1 := world.NewTrackPoint("tpjCp1")
	pdktCp1 := world.NewTrackPoint("pdktCp1")

	bsTpjKrur0 := railway.TrackSegment{
		Id:     "bsTpjKrur",
		Length: units.KM(7),
	}
	bsKrurPdkt0 := railway.TrackSegment{
		Id:     "bsKrurPdkt0",
		Length: units.KM(7),
	}

	world.TrackGraph.AddTrack(tpjSw1, tpjCp1, &bsTpjPdkt0)
	world.TrackGraph.AddTrack(tpjCp1, krurCp1, &bsTpjKrur0)
	world.TrackGraph.AddTrack(krurCp1, pdktCp1, &bsKrurPdkt0)
	world.TrackGraph.AddTrack(tpjCp1, pdktCp1, &bsTpjPdkt1)
	world.TrackGraph.AddTrack(pdktCp1, pdktSw1, &bsTpjPdkt2)

	bsecTpjPdkt.AddTrack(&bsTpjPdkt0)
	bsecTpjPdkt.AddTrack(&bsTpjPdkt1)
	bsecTpjPdkt.AddTrack(&bsTpjPdkt2)

	// this is a temp call -> TODO: Move it to Graph.FindPath() call instead or something
	world.TrackGraph.BuildCacheMap()

	sim.Run()

}
