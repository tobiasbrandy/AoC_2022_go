package day15

import (
	"errors"
	"regexp"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/pos"
	"github.com/tobiasbrandy/AoC_2022_go/internal/regext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/set"
)

var sensorDataParser = regexp.MustCompile(
	`Sensor at x=(?P<sensorX>-?\d+), y=(?P<sensorY>-?\d+): closest beacon is at x=(?P<beaconX>-?\d+), y=(?P<beaconY>-?\d+)`)

func parseSensorData(sensorData string) (sensor, beacon pos.D2) {
	sensorInfo := regext.NamedCaptureGroups(sensorDataParser, sensorData)
	sensor = pos.New2D(parse.Int(sensorInfo["sensorX"]), parse.Int(sensorInfo["sensorY"]))
	beacon = pos.New2D(parse.Int(sensorInfo["beaconX"]), parse.Int(sensorInfo["beaconY"]))
	return sensor, beacon
}

func Part1(inputPath string) any {
	const rowY int = 2000000

	invalidXSet := set.Set[int]{}

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		sensor, beacon := parseSensorData(line)

		dist := sensor.Distance1(beacon)
		remainder := dist - mathext.IntAbs(sensor.Y-rowY)
		if remainder >= 0 {
			for x := sensor.X - remainder; x <= sensor.X+remainder; x++ {
				invalidXSet.Add(x)
			}
			if beacon.Y == rowY {
				invalidXSet.Remove(beacon.X)
			}
		}
	})

	return invalidXSet.Len()
}

type SensorArea struct {
	center pos.D2
	dist   int
}

func (sr SensorArea) Contains(pos pos.D2) bool {
	return sr.center.Distance1(pos) <= sr.dist
}

func (sr SensorArea) Perimeter(out chan pos.D2) {
	defer close(out)

	right := pos.New2D(sr.center.X+sr.dist+1, sr.center.Y)
	down := pos.New2D(sr.center.X, sr.center.Y+sr.dist+1)
	left := pos.New2D(sr.center.X-sr.dist-1, sr.center.Y)
	up := pos.New2D(sr.center.X, sr.center.Y-sr.dist-1)

	// Right -> Down
	for per := right; per != down; per.X, per.Y = per.X-1, per.Y+1 {
		out <- per
	}

	// Down -> Left
	for per := down; per != left; per.X, per.Y = per.X-1, per.Y-1 {
		out <- per
	}

	// Left -> Up
	for per := left; per != up; per.X, per.Y = per.X+1, per.Y-1 {
		out <- per
	}

	// Up -> Right
	for per := up; per != right; per.X, per.Y = per.X+1, per.Y+1 {
		out <- per
	}
}

func isInAnyArea(areas []SensorArea, pos pos.D2) bool {
	for _, area := range areas {
		if area.Contains(pos) {
			return true
		}
	}

	return false
}

func emptyPos(areas []SensorArea, lowerLim, upperLim int) pos.D2 {
	for _, area := range areas {
		perimeter := make(chan pos.D2)
		go area.Perimeter(perimeter)
		for p := range perimeter {
			inBounds := p.X >= lowerLim && p.X <= upperLim && p.Y >= lowerLim && p.Y <= upperLim
			if inBounds && !isInAnyArea(areas, p) {
				return p
			}
		}
	}

	errexit.HandleMainError(errors.New("no empty position found"))
	return pos.D2{}
}

func Part2(inputPath string) any {
	const (
		lowerLim int = 0
		upperLim int = 4000000

		freqX int = 4000000
		freqY int = 1
	)

	var areas []SensorArea

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		sensor, beacon := parseSensorData(line)

		dist := sensor.Distance1(beacon)
		areas = append(areas, SensorArea{center: sensor, dist: dist})
	})

	beaconPos := emptyPos(areas, lowerLim, upperLim)
	freq := freqX*beaconPos.X + freqY*beaconPos.Y

	return freq
}
