package railway

type RailwayEvent string

const (
	// SIM //
	WorldEntered RailwayEvent = "WORLD_ENTER"
	WorldExited  RailwayEvent = "WORLD_EXIT"

	// STN //
	TrainDwellEnd RailwayEvent = "TRAIN_DWELL_END"
	TrainArrived  RailwayEvent = "TRAIN_ARRIVE"
	TrainDeparted RailwayEvent = "TRAIN_DEPART"

	// TRK //
	TrackReserved RailwayEvent = "TRACK_RESERVE"
	TrackEntered  RailwayEvent = "TRACK_ENTER"
	TrackOccupied RailwayEvent = "TRACK_OCCUPY"
	TrackExited   RailwayEvent = "TRACK_EXIT"
	TrackReleased RailwayEvent = "TRACK_RELEASE"

	// SWT //
	SwitchSet RailwayEvent = "SWITCH_SET"
)
