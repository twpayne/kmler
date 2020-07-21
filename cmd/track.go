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
