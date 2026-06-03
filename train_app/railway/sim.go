package railway

import (
	"fmt"
	"trainapp/des"
)

type Sim struct {
	des   *des.DES[RailwayEvent]
	world *World

	blockSecControllers map[string]*BlockSectionController
	stnControllers      map[string]*StationController
}

func (s *Sim) SetWorld(world *World) *Sim {
	s.world = world
	return s
}

func (s *Sim) InitWorld() {
	if s.world == nil {
		panic("s.world is nil, did you call SetWorld?")
	}
	s.des = &des.DES[RailwayEvent]{}
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
		s.des.Add(train.schedule[0].ArrTime-2.0, train.Name, TrainEnter, train)
	}
}

func (s *Sim) stnCtrller(stnCode string) *StationController {
	return s.stnControllers[stnCode]
}

func (s *Sim) bsecCtrller(bName string) *BlockSectionController {
	return s.blockSecControllers[bName]
}

func (s *Sim) NextEvent() (des.Event[RailwayEvent], bool) {
	if s.des == nil {
		panic("s.des is nil, did you call InitWorld?")
	}
	return s.des.NextEvent()
}

func (s *Sim) ScheduleEvent(t float64, name string, evtype RailwayEvent, data *TrainData) {
	s.des.Add(t, name, evtype, data)
}

func (s *Sim) CurTime() float64 {
	return s.des.CurTime
}

// TODO: move to TrainController, work out a way for it
func (s *Sim) Run() {
	for {
		ev, ok := s.NextEvent()
		if !ok {
			break
		}
		switch RailwayEvent(ev.Type) {
		case TrainEnter:
			train := ev.Data.(*TrainData)
			train.curSchedulePoint = 0 // enters the sim

			s.des.Add(s.des.CurTime+1.0, train.Name, TrainArrived, train)
			fmt.Printf("[%.2f] %s entered the simulation\n", s.des.CurTime, train.Name)
			// pf.track.OccupiedBy = train // occupy

		case TrainArrived:
			train := ev.Data.(*TrainData)
			trainSchedule := train.schedule
			curSchedule := trainSchedule[train.curSchedulePoint]

			curOccup := train.occupation

			curStn, ok := s.world.GetStation(curSchedule.StnCode)
			if !ok {
				panic("Station seems to be missing")
			}
			ctrller := s.stnCtrller(curStn.Code)
			_, ok = ctrller.Acquire(train, "0")
			if !ok {
				fmt.Printf("[%.2f] %s waiting to enter station %s - %s\n", s.des.CurTime, train.Name, curStn.Name, curStn.Code)
				break
			}
			if curOccup != nil {
				curOccup.ctrller.Release(curOccup)
			}

			dwellTime := curSchedule.ExpDwellTime()

			fmt.Printf("[%.2f] %s arrived at %s - %s\n", s.des.CurTime, train.Name, curStn.Name, curStn.Code)
			s.des.Add(s.des.CurTime+dwellTime, train.Name, TrainDwellEnd, train)

		case TrainDwellEnd:
			train := ev.Data.(*TrainData)
			s.des.Add(s.des.CurTime+1.0, train.Name, TrainDeparted, train)

		case TrackReleased:
			train := ev.Data.(*TrainData)
			// train is waked up here..
			// find the state
			fmt.Printf("[%.2f] %s track release received, waking up...\n", s.des.CurTime, train.Name)
			if _, ok := train.occupation.ctrller.(*BlockSectionController); ok {
				// so we found that we are waiting to enter station.
				s.des.Add(s.des.CurTime+1.0, train.Name, TrainArrived, train)
			} else if _, ok := train.occupation.ctrller.(*StationController); ok {
				// we are watiing to enter block section, just go back to traindeparting state
				s.des.Add(s.des.CurTime+1.0, train.Name, TrainDeparted, train)
			}

		case TrainDeparted:
			train := ev.Data.(*TrainData)
			curSchedule := train.schedule[train.curSchedulePoint]
			curStn, ok := s.world.GetStation(curSchedule.StnCode)
			curOccp := train.occupation
			if !ok {
				panic("Station seems to be missing")
			}
			fmt.Printf("[%.2f] %s departed from %s - %s\n", s.des.CurTime, train.Name, curStn.Name, curStn.Code)
			if train.curSchedulePoint+1 >= len(train.schedule) {
				train.curSchedulePoint += 1
				curOccp.ctrller.Release(curOccp)
				s.des.Add(s.des.CurTime+1.0, train.Name, TrainExit, train)
			} else {
				nextSchedule := train.schedule[train.curSchedulePoint+1]

				bSec, err := s.world.FindBlockBwStns(curSchedule.StnCode, nextSchedule.StnCode)
				if err != nil {
					panic(err)
				}
				ctrller := s.bsecCtrller(bSec.Id)
				if ctrller == nil {
					panic("Controller is not available for bsec.id")
				}
				_, ok := ctrller.Acquire(train, "0")
				if !ok {
					fmt.Printf("[%.2f] %s waiting for free track\n", s.des.CurTime, train.Name)
					// don't bother scheduling anything
					break
				}
				curOccp.ctrller.Release(curOccp)
				train.curSchedulePoint += 1
				s.des.Add(s.des.CurTime+100.0, train.Name, TrainArrived, train)
			}
		case TrainExit:
			train := ev.Data.(*TrainData)
			fmt.Printf("[%.2f] %s exited simulation\n", s.des.CurTime, train.Name)
		}
	}
}
