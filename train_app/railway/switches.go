package railway

import "fmt"

type Switch struct {
	Id string

	Routes map[string]*SwitchRoute

	lockedRoute *SwitchRoute
}

type SwitchRoute struct {
	Id string

	FromTrack *TrackData
	ToTrack   *TrackData
}

func (s *Switch) Init() {
	s.Routes = make(map[string]*SwitchRoute)
}

func (s *Switch) AddRoute(fromTrack *TrackData, toTrack *TrackData) string {
	route := &SwitchRoute{
		Id:        fmt.Sprintf("%s-%s", fromTrack.Id, toTrack.Id),
		FromTrack: fromTrack,
		ToTrack:   toTrack,
	}
	s.Routes[route.Id] = route
	return route.Id
}

func (s *Switch) RemoveRoute(routeId string) {
	delete(s.Routes, routeId)
}

func (s *Switch) LockRoute(routeId string) bool {
	lockedRoute, ok := s.Routes[routeId]
	s.lockedRoute = lockedRoute
	return ok
}
