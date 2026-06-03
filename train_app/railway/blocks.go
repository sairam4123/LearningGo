package railway

type BlockSection struct {
	stnA *Station
	stnB *Station
	Id   string

	tracks []*TrackData

	// Waiting []*TrainData
}

func (bsec *BlockSection) Init(stnA *Station, stnB *Station) {
	bsec.tracks = make([]*TrackData, 0)
	bsec.stnA = stnA
	bsec.stnB = stnB
}

func (bsec *BlockSection) AddTrack(td *TrackData) {
	bsec.tracks = append(bsec.tracks, td)
}
