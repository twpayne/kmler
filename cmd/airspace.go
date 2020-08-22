package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-openaip"
)

var airspaceCmd = &cobra.Command{
	Use:   "airspace",
	Args:  cobra.MinimumNArgs(1),
	Short: "Convert airspace to KML",
	RunE:  config.runAirspaceCmdE,
}

func init() {
	rootCmd.AddCommand(airspaceCmd)
}

func (c *Config) runAirspaceCmdE(cmd *cobra.Command, args []string) error {
	airspaces := make([]*openaip.OpenAIP, 0, len(args))
	for _, arg := range args {
		airspace, err := parseAirspace(arg)
		if err != nil {
			return err
		}
		airspaces = append(airspaces, airspace)
	}

	airspaceFolders := make([]kml.Element, 0, len(airspaces))
	for _, airspace := range airspaces {
		airspaceFolders = append(airspaceFolders, c.makeAirspaceFolder(airspace))
	}

	sb := &strings.Builder{}
	if err := kml.KML(kml.Document(airspaceFolders...)).WriteIndent(sb, "", "  "); err != nil {
		return err
	}
	return c.writeOutputString(sb.String())
}

func altMeters(alt openaip.Alt) float64 {
	switch alt.Unit {
	case "F":
		return 0.3048 * alt.Value
	case "M":
		return alt.Value
	case "FL":
		return 100 * 0.3048 * alt.Value
	default:
		panic(fmt.Sprintf("%s: unknown altitude unit", alt.Unit))
	}
}

func altLimitMeters(altLimit openaip.AltLimit) float64 {
	return altMeters(altLimit.Value) // FIXME
}

func formatAltLimit(altLimit openaip.AltLimit) string {
	switch {
	case altLimit.Reference == "GND" && altLimit.Value.Value == 0:
		return altLimit.Reference
	case altLimit.Value.Unit == "FL":
		return fmt.Sprintf("%s %d", altLimit.Value.Unit, int(altLimit.Value.Value))
	default:
		return fmt.Sprintf("%d %s %s", int(altLimit.Value.Value), altLimit.Value.Unit, altLimit.Reference)
	}
}

func parseAirspace(arg string) (*openaip.OpenAIP, error) {
	switch strings.ToLower(filepath.Ext(arg)) {
	case ".aip":
		return parseOpenAIP(arg)
	default:
		return nil, fmt.Errorf("%s: unknown extension", arg)
	}
}

func (c *Config) makeAirspaceFolder(oaip *openaip.OpenAIP) kml.Element {
	placemarks := make([]kml.Element, 0, len(oaip.Airspaces))
	for _, airspace := range oaip.Airspaces {
		polygons := make([]kml.Element, 0, len(airspace.Polygons))
		for _, polygon := range airspace.Polygons {
			polygon := kml.Polygon(
				kml.OuterBoundaryIs(
					kml.LinearRing(
						kml.CoordinatesArray(polygon.Coords...),
					),
				),
			)
			polygons = append(polygons, polygon)
		}
		var geometry kml.Element
		if len(polygons) == 1 {
			geometry = polygons[0]
		} else {
			geometry = kml.MultiGeometry(polygons...)
		}
		name := fmt.Sprintf("%s (%s-%s)", airspace.Name, formatAltLimit(airspace.AltLimitBottom), formatAltLimit(airspace.AltLimitTop))
		placemark := kml.Placemark(
			kml.Name(name),
			// kml.MinAltitude(altLimitMeters(airspace.AltLimitBottom)),
			// kml.MaxAltitude(altLimitMeters(airspace.AltLimitTop)),
			geometry,
		)
		placemarks = append(placemarks, placemark)
	}
	return kml.Folder(
		placemarks...,
	)
}

func parseOpenAIP(arg string) (*openaip.OpenAIP, error) {
	f, err := os.Open(arg)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return openaip.Read(f)
}
