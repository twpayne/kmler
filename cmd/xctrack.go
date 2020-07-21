package cmd

import (
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/sphere"
	"github.com/twpayne/go-xctrack"
)

func (c *Config) makeXCTrackTaskRoute(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	task, err := xctrack.ParseTask(data)
	if err != nil {
		return err
	}

	taskFolder := c.makeXCTrackTaskFolder(filename, task)

	sb := &strings.Builder{}
	if err := kml.KML(kml.Document(taskFolder)).WriteIndent(sb, "", "  "); err != nil {
		return err
	}
	return c.writeOutputString(sb.String())
}

func (c *Config) makeXCTrackTaskRoutePlacemark(task *xctrack.Task) kml.Element {
	if len(task.Turnpoints) < 2 {
		panic("too few turnpoints")
	}

	s := sphere.FAI

	coordinates := make([]kml.Coordinate, 0, len(task.Turnpoints))
	for _, turnpoint := range task.Turnpoints {
		coordinate := kml.Coordinate{
			Lon: turnpoint.Waypoint.Lon,
			Lat: turnpoint.Waypoint.Lat,
		}
		coordinates = append(coordinates, coordinate)
	}

	geometries := make([]kml.Element, 0, 2*len(task.Turnpoints))
	for i, turnpoint := range task.Turnpoints {
		if i == 0 {
			continue
		}
		lineString := kml.LineString(
			kml.Tessellate(true),
			kml.Coordinates(
				s.Offset(coordinates[i-1], float64(task.Turnpoints[i-1].Radius), s.InitialBearingTo(coordinates[i-1], coordinates[i])),
				s.Offset(coordinates[i], float64(task.Turnpoints[i].Radius), s.InitialBearingTo(coordinates[i], coordinates[i-1])),
			),
		)
		circle := kml.LineString(
			kml.Coordinates(s.Circle(coordinates[i], float64(turnpoint.Radius), 1.0)...),
		)
		geometries = append(geometries, lineString, circle)
	}
	return kml.Placemark(
		kml.Style(
			kml.LineStyle(
				kml.Color(color.RGBA{R: 0, G: 127, B: 255, A: 255}),
				kml.Width(1),
			),
		),
		kml.MultiGeometry(geometries...),
	)
}

func (c *Config) makeXCTrackTaskFolder(filename string, task *xctrack.Task) kml.Element {
	return kml.Folder(
		kml.Name(filepath.Base(filename)),
		c.makeXCTrackTaskRoutePlacemark(task),
	)
}
