package railway

import "fmt"

type PathController struct {
	Id string

	path *Path

	train *Train
	sim   *Sim
}

func (pathCtrller *PathController) ReservePath() error {
	for _, edge := range pathCtrller.path.Edges {
		if edge.Track.IsReserved() || edge.Track.IsOccupied() {
			return fmt.Errorf("Cannot reserve path when %s is already reserved or occupied, bailing out.", edge.Track.Id)
		}

		edge.Track.Reserve(pathCtrller.train)
	}

	return nil
}
