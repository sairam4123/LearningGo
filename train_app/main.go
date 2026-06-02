package main

import (
	"fmt"
	"trainapp/des"
)

type TrackState string

const (
	Available TrackState = "TRACK_AVAILABLE"
	Reserved  TrackState = "TRACK_RESERVED"
	Occupied  TrackState = "TRACK_OCCUPIED"
)

type TrackDirection string

const (
	Up    TrackDirection = "UP"
	Down  TrackDirection = "DOWN"
	Bidir TrackDirection = "BIDIR"
)

type TrainEvent string

const (
	// SIM //
	TrainEnter TrainEvent = "TRAIN_ENTER"
	TrainExit  TrainEvent = "TRAIN_EXIT"

	// STN //
	TrainDwellEnd TrainEvent = "TRAIN_DWELL"
	TrainArrived  TrainEvent = "TRAIN_ARRIVE"
	TrainDeparted TrainEvent = "TRAIN_DEPART"

	// BS //
	TrainEnterBSec TrainEvent = "TRAIN_ENTER_BSEC"
	TrainExitBSec  TrainEvent = "TRAIN_EXIT_BSEC"

	TrackReleased TrainEvent = "TRACK_RELEASED"
)

func m[T comparable](length T) T {
	return length
}

type TrackData struct {
	Id string

	// Resource related
	ReservedBy *TrainData
	OccupiedBy *TrainData

	Direction TrackDirection
	Length    int
}

type Platform struct {
	Id string

	Length int

	track *TrackData
}

type TrainData struct {
	Name   string
	Number string

	curSchedulePoint int
	schedule         []*SchedulePoint

	occupation *OccupationData
}

type SchedulePoint struct {
	trainNumber string
	stnCode     string
	arrTime     float64
	deptTime    float64
	spPfNo      int
}

type Station struct {
	Code string
	Name string

	Platforms []*Platform
}

// func (stn *Station) PutInQueue(train *TrainData) {
// 	stn.Waiting = append(stn.Waiting, train)
// }

// func (stn *Station) Release() (*TrainData, bool) {
// 	nextTrain := stn.Waiting[0]
// 	stn.Waiting = stn.Waiting[1:]
// 	return nextTrain, nextTrain != nil
// }

type BlockSection struct {
	stnA *Station
	stnB *Station
	Id   string

	tracks []*TrackData

	// Waiting []*TrainData
}

type World struct {
	stations map[string]*Station
	trains   map[string]*TrainData

	bsections map[string]*BlockSection

	occupiedTracks map[*TrainData]*TrackData
}

type IController interface {
	FindAvailableTrack() (*TrackData, error)
	Acquire(train *TrainData, expectedTrackId string) (*OccupationData, bool)
	Release(occup *OccupationData)
}

type StationController struct {
	station *Station
	sim     *Sim

	waiting []*TrainData // Must be FIFO
}

type BlockSectionController struct {
	bsec *BlockSection
	sim  *Sim

	waiting []*TrainData
}

type OccupationData struct {
	track   *TrackData
	ctrller IController
	train   *TrainData
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

func (ctrl *StationController) FindAvailableTrack() (*TrackData, error) {
	for _, pf := range ctrl.station.Platforms {
		// ignore direction for now
		// fmt.Printf("Track Available %s -> %t\n", pf.track.Id, pf.track.IsAvailable())
		if pf.track.IsAvailable() {
			return pf.track, nil
		}
	}
	return nil, fmt.Errorf("Cannot find any available tracks")
}

func (ctrl *StationController) Acquire(train *TrainData, expectedTrackId string) (*OccupationData, bool) {
	// ignore expectedTrackId for now
	track, err := ctrl.FindAvailableTrack()
	if err != nil {
		// queue the train to the block section
		ctrl.waiting = append(ctrl.waiting, train)
		return nil, false
	}

	// create an occupation object
	occup := OccupationData{
		track:   track,
		ctrller: ctrl,
		train:   train,
	}

	track.Request(train)
	train.occupation = &occup

	return &occup, true
}

func (ctrl *StationController) Release(occup *OccupationData) {
	if ctrl != occup.ctrller {
		return
	}

	occup.track.Release()
	if occup == occup.train.occupation {
		occup.train.occupation = nil
	}
	if len(ctrl.waiting) > 0 {
		waitingTrain := ctrl.waiting[0]
		ctrl.waiting = ctrl.waiting[1:]

		ctrl.sim.des.Add(ctrl.sim.des.CurTime+1, waitingTrain.Name, TrackReleased, waitingTrain)
	}
}

func (ctrl *BlockSectionController) FindAvailableTrack() (*TrackData, error) {
	for _, trk := range ctrl.bsec.tracks {
		// fmt.Printf("Track Available %s -> %t\n", trk.Id, trk.IsAvailable())
		// ignore direction for now
		if trk.IsAvailable() {
			return trk, nil
		}
	}
	return nil, fmt.Errorf("Cannot find any available tracks")
}

func (ctrl *BlockSectionController) Acquire(train *TrainData, expectedTrackId string) (*OccupationData, bool) {
	// ignore expectedTrackId for now
	track, err := ctrl.FindAvailableTrack()
	if err != nil {
		// fmt.Println("Waiting Len", len(ctrl.waiting))
		// queue the train to the block section
		ctrl.waiting = append(ctrl.waiting, train)
		return nil, false
	}

	// create an occupation object
	occup := OccupationData{
		track:   track,
		ctrller: ctrl,
		train:   train,
	}

	track.Request(train)
	train.occupation = &occup

	return &occup, true
}

func (ctrl *BlockSectionController) Release(occup *OccupationData) {
	if ctrl != occup.ctrller {
		return
	}

	// fmt.Println("Releasing block section")

	occup.track.Release()
	if occup == occup.train.occupation {
		occup.train.occupation = nil
	}

	// fmt.Println("Waiting Len", len(ctrl.waiting))
	if len(ctrl.waiting) > 0 {

		waitingTrain := ctrl.waiting[0]
		ctrl.waiting = ctrl.waiting[1:]

		ctrl.sim.des.Add(ctrl.sim.des.CurTime+1, waitingTrain.Name, TrackReleased, waitingTrain)
	}
}

// func (bsec *BlockSection) Release() *TrainData {
// 	nextTrain := bsec.Waiting[0]
// 	bsec.Waiting = bsec.Waiting[1:]
// 	return nextTrain
// }

// func (stn *Station) FindAvailablePlatform() (*Platform, error) {
// 	for _, pf := range stn.Platforms {
// 		if pf.track.IsAvailable() {
// 			return pf, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("Cannot find any available platforms")
// }

func (w *World) FindBlockBwStns(aStnCode string, bStnCode string) (*BlockSection, error) {
	for _, bsec := range w.bsections {
		if bsec.stnA.Code == aStnCode && bsec.stnB.Code == bStnCode {
			return bsec, nil
		}
		if bsec.stnA.Code == bStnCode && bsec.stnB.Code == aStnCode {
			return bsec, nil
		}
	}
	return nil, fmt.Errorf("Cannot find any block sections between aStnCode %s <-> bStnCode %s", aStnCode, bStnCode)
}

func (bsec *BlockSection) Init() {
	bsec.tracks = make([]*TrackData, 0)
}

func (bsec *BlockSection) AddTrack(td *TrackData) {
	bsec.tracks = append(bsec.tracks, td)
}

func (stn *Station) Init() {
	stn.Platforms = make([]*Platform, 0)
}

func (stn *Station) AddPlatform(pfData *Platform) {
	stn.Platforms = append(stn.Platforms, pfData)
}

func (w *World) GetStation(code string) (*Station, bool) {
	stn, ok := w.stations[code]
	return stn, ok
}

func (w *World) AddTrain(train *TrainData) {
	w.trains[train.Number] = train
}
func (w *World) RemoveTrain(trainNumber string) {
	delete(w.trains, trainNumber)
}

func (w *World) AddStation(stn *Station) {
	w.stations[stn.Code] = stn
}
func (w *World) RemoveStation(stnCode string) {
	delete(w.stations, stnCode)
}

func (w *World) AddBlockSection(bsec *BlockSection) {
	w.bsections[bsec.Id] = bsec
}

func (t *TrainData) AddSchedule(sp *SchedulePoint) {
	t.schedule = append(t.schedule, sp)
}

type Sim struct {
	des   *des.DES[TrainEvent]
	world *World

	blockSecControllers map[string]*BlockSectionController
	stnControllers      map[string]*StationController
}

func (w *World) Init() {
	w.stations = make(map[string]*Station)
	w.trains = make(map[string]*TrainData)
	w.bsections = make(map[string]*BlockSection)
	w.occupiedTracks = make(map[*TrainData]*TrackData)

}

func (s *Sim) SetWorld(world *World) *Sim {
	s.world = world
	return s
}

func (s *Sim) InitWorld() {
	if s.world == nil {
		panic("s.world is nil, did you call SetWorld?")
	}
	s.des = &des.DES[TrainEvent]{}
	s.des.Init()

	s.blockSecControllers = make(map[string]*BlockSectionController)
	s.stnControllers = make(map[string]*StationController)

	for _, stn := range s.world.stations {
		s.stnControllers[stn.Code] = &StationController{
			station: stn,
			sim:     s,
			waiting: make([]*TrainData, 0),
		}
	}

	for _, bsec := range s.world.bsections {
		s.blockSecControllers[bsec.Id] = &BlockSectionController{
			bsec:    bsec,
			sim:     s,
			waiting: make([]*TrainData, 0),
		}
	}

	for _, train := range s.world.trains {
		s.des.Add(train.schedule[0].arrTime-2.0, train.Name, TrainEnter, train)
	}

}

func (s *Sim) stnCtrller(stnCode string) *StationController {
	return s.stnControllers[stnCode]
}

func (s *Sim) bsecCtrller(bName string) *BlockSectionController {
	return s.blockSecControllers[bName]
}

func (s *Sim) NextEvent() (des.Event[TrainEvent], bool) {
	if s.des == nil {
		panic("s.des is nil, did you call InitWorld?")
	}
	return s.des.NextEvent()
}

func main() {
	sim := Sim{}

	world := World{}
	world.Init()

	_12606 := TrainData{
		Name:   "Pallavan",
		Number: "12606",
	}

	_12605 := TrainData{
		Name:   "RMM Express",
		Number: "12605",
	}

	pdkt := Station{
		Code: "PDKT",
		Name: "Pudukkottai",
	}

	pdkt.Init()

	kkdi := Station{
		Code: "KKDI",
		Name: "Karaikudi",
	}

	kkdi.Init()

	trTPJ_KKDI := TrackData{
		Id:        "TPJ_KKDI",
		Direction: Bidir,
	}

	bsTPJ_KKDI := BlockSection{
		stnA: &kkdi,
		stnB: &pdkt,
	}

	bsTPJ_KKDI.AddTrack(&trTPJ_KKDI)

	trKkdi0 := TrackData{
		Id:        "kkdi0",
		Direction: Bidir,
		Length:    m(700),
	}
	trKkdi1 := TrackData{
		Id:        "kkdi1",
		Direction: Bidir,
		Length:    m(700),
	}

	trPdkt0 := TrackData{
		Id:        "pdkt0",
		Direction: Bidir,
		Length:    m(700),
	}
	trPdkt1 := TrackData{
		Id:        "pdkt1",
		Direction: Bidir,
		Length:    m(500),
	}

	world.AddStation(&kkdi)
	world.AddStation(&pdkt)

	world.AddBlockSection(&bsTPJ_KKDI)
	kkdi0 := Platform{
		Id:     trKkdi0.Id,
		track:  &trKkdi0,
		Length: m(500),
	}
	kkdi1 := Platform{
		Id:     trKkdi1.Id,
		track:  &trKkdi1,
		Length: m(500),
	}
	pdkt0 := Platform{
		Id:     trPdkt0.Id,
		track:  &trPdkt0,
		Length: m(500),
	}
	pdkt1 := Platform{
		Id:     trPdkt1.Id,
		track:  &trPdkt1,
		Length: m(500),
	}

	kkdi.AddPlatform(&kkdi0)
	kkdi.AddPlatform(&kkdi1)
	pdkt.AddPlatform(&pdkt0)
	pdkt.AddPlatform(&pdkt1)

	world.AddTrain(&_12605)
	world.AddTrain(&_12606)

	_12606.AddSchedule(&SchedulePoint{
		trainNumber: _12606.Number,
		stnCode:     kkdi.Code,
		arrTime:     120,
		deptTime:    125,
		spPfNo:      0,
	})

	_12606.AddSchedule(&SchedulePoint{
		trainNumber: _12606.Number,
		stnCode:     pdkt.Code,
		arrTime:     145,
		deptTime:    150,
		spPfNo:      0,
	})

	_12605.AddSchedule(&SchedulePoint{
		trainNumber: _12605.Number,
		stnCode:     pdkt.Code,
		arrTime:     130,
		deptTime:    140,
		spPfNo:      0,
	})

	_12605.AddSchedule(&SchedulePoint{
		trainNumber: _12605.Number,
		stnCode:     kkdi.Code,
		arrTime:     320,
		deptTime:    340,
		spPfNo:      0,
	})

	sim.SetWorld(&world)
	sim.InitWorld()

	for {
		ev, ok := sim.NextEvent()
		if !ok {
			break
		}
		switch TrainEvent(ev.Type) {
		case TrainEnter:
			train := ev.Data.(*TrainData)
			train.curSchedulePoint = 0 // enters the sim

			sim.des.Add(sim.des.CurTime+1.0, train.Name, TrainArrived, train)
			fmt.Printf("[%.2f] %s entered the simulation\n", sim.des.CurTime, train.Name)
			// pf.track.OccupiedBy = train // occupy

		case TrainArrived:
			train := ev.Data.(*TrainData)
			trainSchedule := train.schedule
			curSchedule := trainSchedule[train.curSchedulePoint]

			curOccup := train.occupation

			curStn, ok := sim.world.GetStation(curSchedule.stnCode)
			if !ok {
				panic("Station seems to be missing")
			}
			ctrller := sim.stnCtrller(curStn.Code)
			_, ok = ctrller.Acquire(train, "0")
			if !ok {
				fmt.Printf("[%.2f] %s waiting to enter station %s - %s\n", sim.des.CurTime, train.Name, curStn.Name, curStn.Code)
				break
			}
			if curOccup != nil {
				curOccup.ctrller.Release(curOccup)
			}

			dwellTime := curSchedule.deptTime - curSchedule.arrTime

			fmt.Printf("[%.2f] %s arrived at %s - %s\n", sim.des.CurTime, train.Name, curStn.Name, curStn.Code)
			sim.des.Add(sim.des.CurTime+dwellTime, train.Name, TrainDwellEnd, train)

		case TrainDwellEnd:
			train := ev.Data.(*TrainData)
			sim.des.Add(sim.des.CurTime+1.0, train.Name, TrainDeparted, train)

		case TrackReleased:
			train := ev.Data.(*TrainData)
			// train is waked up here..
			// find the state
			fmt.Printf("[%.2f] %s track release received, waking up...\n", sim.des.CurTime, train.Name)
			if _, ok := train.occupation.ctrller.(*BlockSectionController); ok {
				// so we found that we are waiting to enter station.
				sim.des.Add(sim.des.CurTime+1.0, train.Name, TrainArrived, train)
			} else if _, ok := train.occupation.ctrller.(*StationController); ok {
				// we are watiing to enter block section, just go back to traindeparting state
				sim.des.Add(sim.des.CurTime+1.0, train.Name, TrainDeparted, train)
			}

		// todo implement
		case TrainDeparted:
			train := ev.Data.(*TrainData)
			curSchedule := train.schedule[train.curSchedulePoint]
			curStn, ok := sim.world.GetStation(curSchedule.stnCode)
			curOccp := train.occupation
			if !ok {
				panic("Station seems to be missing")
			}
			fmt.Printf("[%.2f] %s departed from %s - %s\n", sim.des.CurTime, train.Name, curStn.Name, curStn.Code)
			if train.curSchedulePoint+1 >= len(train.schedule) {
				train.curSchedulePoint += 1
				curOccp.ctrller.Release(curOccp)
				sim.des.Add(sim.des.CurTime+1.0, train.Name, TrainExit, train)
			} else {
				nextSchedule := train.schedule[train.curSchedulePoint+1]

				bSec, err := world.FindBlockBwStns(curSchedule.stnCode, nextSchedule.stnCode)
				if err != nil {
					panic(err)
				}
				ctrller := sim.bsecCtrller(bSec.Id)
				if ctrller == nil {
					panic("Controller is not available for bsec.id")
				}
				_, ok := ctrller.Acquire(train, "0")
				if !ok {
					fmt.Printf("[%.2f] %s waiting for free track\n", sim.des.CurTime, train.Name)
					// don't bother scheduling anything
					break
				}
				curOccp.ctrller.Release(curOccp)
				train.curSchedulePoint += 1
				sim.des.Add(sim.des.CurTime+100.0, train.Name, TrainArrived, train)
			}
		case TrainExit:
			train := ev.Data.(*TrainData)
			fmt.Printf("[%.2f] %s exited simulation\n", sim.des.CurTime, train.Name)
		}
	}
}
