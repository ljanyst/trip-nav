package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/antchfx/xmlquery"
	. "github.com/ljanyst/gosh"
	log "github.com/sirupsen/logrus"
)

type Waypoint struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type List struct {
	Title     string
	Waypoints []Waypoint
}

func main() {
	var kml = flag.String("kml", "", "the input KML file")
	var tmplFile = flag.String("template", "list.tmpl", "the output template")
	var out = flag.String("out", "", "the output file")
	flag.Parse()

	SetupLogging("Info")

	if *out == "" {
		log.Fatalf("You need to specify the output file")
	}

	f, err := os.Open(*kml)
	if err != nil {
		log.Fatalf("Cannot open input KML file: %s", err)
	}
	defer f.Close()

	doc, err := xmlquery.Parse(f)
	if err != nil {
		log.Fatalf("Cannot parse input KML file: %s", err)
	}

	title := xmlquery.FindOne(doc, "/kml/Document/name/text()")
	if title == nil {
		log.Fatalf("Cannot extract map title")
	}
	log.Infof("Map title: %+v", title.Data)

	marks := xmlquery.Find(doc, "/kml/Document/Folder/Placemark")
	points := []Waypoint{}
	skipped := 0
	duplicates := 0
	previous := Waypoint{}
	for _, mark := range marks {
		coords := xmlquery.FindOne(mark, "./Point/coordinates/text()")
		if coords == nil {
			skipped++
			continue
		}
		name := xmlquery.FindOne(mark, "./name/text()")
		if name == nil {
			log.Error("Record has no name, skipping")
			skipped++
			continue
		}

		wp := Waypoint{}
		wp.Name = strings.TrimSpace(name.Data)
		coordsArr := strings.Split(coords.Data, ",")
		if wp.Latitude, err = strconv.ParseFloat(strings.TrimSpace(coordsArr[1]), 64); err != nil {
			log.Errorf("Cannot parse latitude as double: %s", err)
		}
		if wp.Longitude, err = strconv.ParseFloat(strings.TrimSpace(coordsArr[0]), 64); err != nil {
			log.Errorf("Cannot parse longitude as double: %s", err)
		}
		if previous.Name == wp.Name {
			duplicates++
			continue
		}
		previous = wp
		points = append(points, wp)

		log.Infof("(%f, %f) %s", wp.Latitude, wp.Longitude, wp.Name)

	}

	log.Infof("Skipped points: %d", skipped)
	log.Infof("Duplicate points: %d", duplicates)
	log.Infof("Valid points: %d", len(points))

	tmplData, err := ioutil.ReadFile(*tmplFile)
	if err != nil {
		log.Fatalf("Cannot read the template file: %s", err)
	}

	list := List{
		title.Data,
		points,
	}

	tmpl, err := template.New("list").Parse(string(tmplData))
	if err != nil {
		log.Fatalf("Cannot parse the tamplate: %s", err)
	}

	outFile, err := os.Create(*out)
	if err != nil {
		log.Fatalf("Cannot create the output file: %s", err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, list)
	if err != nil {
		log.Fatalf("Cannot write the output file: %s", err)
	}

	log.Infof("Written the waypoints to: %s", *out)
}
