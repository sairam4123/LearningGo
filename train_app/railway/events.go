package railway

type RailwayEvent string

const (
	// SIM //
	TrainEnter RailwayEvent = "TRAIN_ENTER"
	TrainExit  RailwayEvent = "TRAIN_EXIT"

	// STN //
	TrainDwellEnd RailwayEvent = "TRAIN_DWELL"
	TrainArrived  RailwayEvent = "TRAIN_ARRIVE"
	TrainDeparted RailwayEvent = "TRAIN_DEPART"

	// BS //
	TrainEnterBSec RailwayEvent = "TRAIN_ENTER_BSEC"
	TrainExitBSec  RailwayEvent = "TRAIN_EXIT_BSEC"

	TrackReleased RailwayEvent = "TRACK_RELEASED"
)
