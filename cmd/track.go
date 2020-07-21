package cmd

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/icon"
)

var trackCmd = &cobra.Command{
	Use:   "track",
	Args:  cobra.MinimumNArgs(1),
	Short: "Convert a track to KML",
	RunE:  config.runTrackCmdE,
}

type track interface {
	Filename() string
	Geom() geom.T
}

func init() {
	rootCmd.AddCommand(trackCmd)
}

func (c *Config) runTrackCmdE(cmd *cobra.Command, args []string) error {
	tracks := make([]track, 0, len(args))
	for _, arg := range args {
		track, err := parseTrack(arg)
		if err != nil {
			return err
		}
		tracks = append(tracks, track)
	}

	trackFolders := make([]kml.Element, 0, len(tracks))
	for _, track := range tracks {
		trackFolders = append(trackFolders, c.makeTrackFolder(track))
	}

	sb := &strings.Builder{}
	if err := kml.GxKML(kml.Document(trackFolders...)).WriteIndent(sb, "", "  "); err != nil {
		return err
	}
	return c.writeOutputString(sb.String())
}

func (c *Config) makeAltitudeMarksFolder(t track, children []kml.Element) kml.Element {
	altitudeMarkStyle := kml.SharedStyle("altitudeMark",
		kml.IconStyle(
			kml.HotSpot(kml.Vec2{ // FIXME HotSpot does not seem to be respected
				X:      0.5,
				Y:      0.5,
				XUnits: kml.UnitsFraction,
				YUnits: kml.UnitsFraction,
			}),
			kml.Icon(
				kml.Href(icon.PaletteHref(4, 24)),
			),
		),
		kml.LabelStyle(
			kml.Scale(0.5),
		),
	)
	ls, ok := t.Geom().(*geom.LineString)
	if !ok {
		panic(fmt.Sprintf("%T: unsupported geometry type", t.Geom()))
	}
	n := ls.NumCoords()
	zs := make([]float64, 0, n)
	for i := 0; i < n; i++ {
		zs = append(zs, ls.Coord(i)[2])
	}
	indexes := salient(zs, 100)
	placemarks := make([]kml.Element, 0, len(indexes))
	for _, index := range indexes {
		coord := ls.Coord(index)
		placemark := kml.Placemark(
			kml.Name(fmt.Sprintf("%dm", int(coord[2]+0.5))),
			kml.StyleURL(altitudeMarkStyle.URL()),
			kml.Point(
				kml.AltitudeMode(kml.AltitudeModeAbsolute),
				kml.Coordinates(kml.Coordinate{
					Lon: coord[0],
					Lat: coord[1],
					Alt: coord[2],
				}),
			),
		)
		placemarks = append(placemarks, placemark)
	}
	return kml.Folder(append(append(children, altitudeMarkStyle), placemarks...)...)
}

func (c *Config) makeTrackFolder(t track) kml.Element {
	var (
		noIconStyle = kml.IconStyle(
			kml.Icon(
				kml.Href(icon.PaletteHref(2, 15)),
			),
		)
		noLabelStyle = kml.LabelStyle(
			kml.Scale(0),
		)
		trackLineStyle = kml.LineStyle(
			kml.Color(color.RGBA{R: 255, G: 0, B: 0, A: 255}),
			kml.Width(2),
		)
		shadowLineStyle = kml.LineStyle(
			kml.Color(color.RGBA{R: 0, G: 0, B: 0, A: 127}),
			kml.Width(1),
		)
	)

	return kml.Folder(
		kml.Name(filepath.Base(t.Filename())),
		c.makeTrackPlacemark(t,
			[]kml.Element{
				kml.Name("Track"),
				kml.Style(
					noIconStyle,
					noLabelStyle,
					trackLineStyle,
				),
			},
			[]kml.Element{
				kml.AltitudeMode(kml.AltitudeModeAbsolute),
			},
		),
		c.makeTrackPlacemark(t,
			[]kml.Element{
				kml.Name("Shadow"),
				kml.Style(
					noIconStyle,
					noLabelStyle,
					shadowLineStyle,
				),
			},
			[]kml.Element{
				kml.AltitudeMode(kml.AltitudeModeClampToGround),
			},
		),
		c.makeAltitudeMarksFolder(t,
			[]kml.Element{
				kml.Name("Altitude marks"),
				kml.Style(
					kml.ListStyle( // FIXME this style does not seem to be respected
						kml.ListItemType(kml.ListItemTypeCheckHideChildren),
					),
				),
			},
		),
	)
}

func (c *Config) makeTrackPlacemark(t track, children, geometryChildren []kml.Element) kml.Element {
	if t.Geom().Layout() < geom.XYZM {
		panic("unsupported geometry layout")
	}
	switch g := t.Geom().(type) {
	case *geom.LineString:
		children = append(children,
			c.makeLineStringGxTrack(g, geometryChildren),
		)
	default:
		panic(fmt.Sprintf("%T: unsupported geometry", g))
	}
	return kml.Placemark(children...)
}

func (c *Config) makeLineStringGxTrack(ls *geom.LineString, children []kml.Element) kml.Element {
	n := ls.NumCoords()
	gxCoordsAndWhens := make([]kml.Element, 0, 2*n)
	for i := 0; i < n; i++ {
		c := ls.Coord(i)
		gxCoordsAndWhens = append(gxCoordsAndWhens,
			kml.GxCoord(kml.Coordinate{
				Lon: c[0],
				Lat: c[1],
				Alt: c[2],
			}),
			kml.When(
				time.Unix(int64(c[3]), 0),
			),
		)
	}
	return kml.GxTrack(append(children, gxCoordsAndWhens...)...)
}

func parseTrack(filename string) (track, error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".igc":
		return parseIGC(filename)
	default:
		return nil, fmt.Errorf("%s: unknown extension", filename)
	}
}
