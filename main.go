package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"pkg/text/template"
	"strings"

	"github.com/kellydunn/golang-geo"
)

type Result struct {
	Address string
	Point   *geo.Point
}

var geocoder geo.Geocoder

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	geocoder = new(geo.GoogleGeocoder)

	query := strings.Join(args, " ")

	point, err := geocoder.Geocode(query)
	if err != nil {
		printErr(err)
	}

	addr, err := geocoder.ReverseGeocode(point)
	if err != nil {
		printErr(err)
	}

	result := Result{
		Address: addr,
		Point:   point,
	}

	tmpl(os.Stdout, resultTemplate, result)
}

const usageTemplate = `geocode is a command-line tool to geocode addresses

Usage:

  geocode [address]

`

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, nil)
}

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(2)
}

const resultTemplate = `Address: {{.Address}}
Coordinates: {{.Point.Lat}}, {{.Point.Lng}}

`
