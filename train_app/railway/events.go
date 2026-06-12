package railway

type RailwayEvent string

const (
	// SIM //
	TrainEntered RailwayEvent = "TRAIN_ENTER"
	TrainExited  RailwayEvent = "TRAIN_EXIT"

	// STN //
	TrainDwellEnd RailwayEvent = "TRAIN_DWELL_END"
	TrainArrived  RailwayEvent = "TRAIN_ARRIVE"
	TrainDeparted RailwayEvent = "TRAIN_DEPART"

	// TRK //
	TrackReserved RailwayEvent = "TRACK_RESERVE"
	TrackOccupied RailwayEvent = "TRACK_OCCUPY"
	TrackReleased RailwayEvent = "TRACK_RELEASE"

	// SWT //
	SwitchSet RailwayEvent = "SWITCH_SET"
)
