package cmd

import (
	"bytes"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/twpayne/go-kml"
	"github.com/twpayne/go-kml/icon"
	"github.com/twpayne/go-waypoint"
)

var customVercofly2020Cmd = &cobra.Command{
	Use:   "vercofly-2020",
	Args:  cobra.NoArgs,
	Short: "Generate a KML file for the Vercofly 2020",
	RunE:  config.runCustomVercofly2020CmdE,
}

func init() {
	customCmd.AddCommand(customVercofly2020Cmd)
}

func (c *Config) runCustomVercofly2020CmdE(cmd *cobra.Command, args []string) error {
	waypointData, err := os.ReadFile("data/vercofly-2020/Vercofly2020_COMP.wpt")
	if err != nil {
		return err
	}
	waypoints, _, err := waypoint.Read(bytes.NewReader(waypointData))
	if err != nil {
		return err
	}
	waypoints = append(waypoints,
		&waypoint.T{
			ID:        "Bonus_FOOD_3",
			Latitude:  46.34470,
			Longitude: 7.86187,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_1",
			Latitude:  46.05061,
			Longitude: 7.48015,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_2",
			Latitude:  46.09672,
			Longitude: 7.54214,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "Bonus_SPECIAL_",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "",
			Latitude:  0,
			Longitude: 0,
		},
		&waypoint.T{
			ID:        "",
			Latitude:  0,
			Longitude: 0,
		},
	)

	renames := map[string]string{
		"START":                       "Start",
		"Cabane_de_Tracuit":           "Tracuit",
		"Cabane_Arpitettaz":           "Arpitettaz",
		"Cabane_du_Grand_Mountet":     "Grand Mountet",
		"Cabane_des_Becs_de_Bosson":   "Becs de Bosson",
		"Cabane_des_Aiguilles_Rouges": "Aiguilles Rouges",
		"Cabane_de_la_Tsa":            "Tsa",
		"Cabane_de_la_Dent_Blanche":   "Dent Blanche",
		"Cabane_Prafleuri":            "Prafleuri",
		"Bonus_SUPAIR_1":              "Selfie",
		"Bonus_SUPAIR_2":              "Selfie",
		"Bonus_FOOD_1":                "Food 1",
		"Bonus_FOOD_2":                "Food 2",
		"FINISH":                      "Finish",
		"Bonus_FOOD_3":                "Food 3",
		"Bonus_SPECIAL_1":             "Swim",
		"Bonus_SPECIAL_2":             "Selfie",
	}
	prefixPaddles := map[string]string{
		"START":        "go",
		"Cabane":       "red-circle",
		"Bonus_SUPAIR": "blu-stars",
		"Bonus_FOOD":   "ylw-stars",
		"FINISH":       "stop",
	}

	children := make([]kml.Element, 0, 1+len(waypoints))
	children = append(children,
		kml.Name("Vercofly 2020"),
	)
	for _, waypoint := range waypoints {
		var paddle string
		for prefix, prefixPaddle := range prefixPaddles {
			if strings.HasPrefix(waypoint.ID, prefix) {
				paddle = prefixPaddle
				break
			}
		}
		placemark := kml.Placemark(
			kml.Name(renames[waypoint.ID]),
			kml.Style(
				kml.IconStyle(
					kml.HotSpot(kml.Vec2{X: 0.5, Y: 0, XUnits: kml.UnitsFraction, YUnits: kml.UnitsFraction}),
					kml.Icon(
						kml.Href(icon.PaddleHref(paddle)),
					),
				),
				kml.LabelStyle(
					kml.Scale(0.5),
				),
			),
			kml.Point(
				kml.Coordinates(kml.Coordinate{
					Lon: waypoint.Longitude,
					Lat: waypoint.Latitude,
				}),
			),
		)
		children = append(children, placemark)
	}

	sb := &strings.Builder{}
	if err := kml.KML(kml.Document(kml.Folder(children...))).WriteIndent(sb, "", "  "); err != nil {
		return err
	}
	return c.writeOutputString(sb.String())
}
