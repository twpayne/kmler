package cmd

import (
	"os"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/igc"
)

type igcTrack struct {
	filename string
	igcT     *igc.T
}

func (t *igcTrack) Filename() string {
	return t.filename
}

func (t *igcTrack) Geom() geom.T {
	return t.igcT.LineString
}

func parseIGC(filename string) (track, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	igcT, err := igc.Read(f)
	if err != nil {
		return nil, err
	}

	return &igcTrack{
		filename: filename,
		igcT:     igcT,
	}, nil
}
