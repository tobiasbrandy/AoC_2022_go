package day16

import (
	"fmt"
	"github.com/gammazero/deque"
	"github.com/tobiasbrandy/AoC_2022_go/internal/hashext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/priorityq"
	"regexp"
	"strings"

	"github.com/tobiasbrandy/AoC_2022_go/internal/stringer"

	"github.com/tobiasbrandy/AoC_2022_go/internal/set"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/regext"
)

type Valve struct {
	name     string
	flowRate int
	tunnels  []string
}

func (v *Valve) String() string {
	return stringer.String(v)
}

type ValvePair struct {
	left, right string
}

func (vp ValvePair) Rotate() ValvePair {
	return ValvePair{vp.right, vp.left}
}

var valvesDataParser = regexp.MustCompile(
	`Valve (?P<name>[A-Z]+) has flow rate=(?P<flowRate>\d+); tunnels? leads? to valves? (?P<tunnels>[A-Z]+(, [A-Z]+)*)`)

func parseValve(data string) *Valve {
	valveData := regext.NamedCaptureGroups(valvesDataParser, data)
	name := valveData["name"]
	flowRate := parse.Int(valveData["flowRate"])
	tunnels := strings.Split(valveData["tunnels"], ", ")
	return &Valve{
		name:     name,
		flowRate: flowRate,
		tunnels:  tunnels,
	}
}

type TraversalState struct {
	valve string
	cost  int
}

func (ts TraversalState) Less(o TraversalState) bool {
	return o.cost < ts.cost
}

type State struct {
	valve    string
	opened   set.Set[string]
	timeLeft int
	pressure int
}

func (s State) String() string {
	return stringer.String(s)
}

func parseValves(inputPath string) (valves map[string]*Valve, closed set.Set[string]) {
	valves = make(map[string]*Valve)
	closed = set.Set[string]{}

	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		valve := parseValve(line)
		valveName := valve.name

		valves[valveName] = valve
		if valve.flowRate > 0 {
			closed.Add(valveName)
		}
	})

	return valves, closed
}

func minDist(valves map[string]*Valve, src, dst string) int {
	if src == dst {
		return 0
	}

	pq := priorityq.New[TraversalState]()
	visited := set.Set[string]{}

	pq.Push(TraversalState{
		valve: src,
		cost:  0,
	})
	visited.Add(src)

	for pq.Len() > 0 {
		state := pq.Pop()
		if state.valve == dst {
			return state.cost
		}

		for _, v := range valves[state.valve].tunnels {
			if !visited.Contains(v) {
				pq.Push(TraversalState{
					valve: v,
					cost:  state.cost + 1,
				})
				visited.Add(v)
			}
		}
	}

	errexit.HandleMainError(fmt.Errorf("no path between %v and %v", src, dst))
	return -1
}

func closedValvesMinDistances(valves map[string]*Valve, closed set.Set[string]) map[ValvePair]int {
	// Could be computed with Floyd-Warshall, maybe more efficient
	minDists := make(map[ValvePair]int)

	for vl := range closed {
		for vr := range closed {
			pair := ValvePair{left: vl, right: vr}
			if _, ok := minDists[pair]; !ok {
				dist := minDist(valves, vl, vr)
				minDists[pair] = dist
				minDists[pair.Rotate()] = dist
			}
		}
	}

	return minDists
}

func Part1(inputPath string) any {
	const (
		initValve = "AA"
		totalTime = 30
	)

	valves, closed := parseValves(inputPath)
	closedC := len(closed)

	closed.Add(initValve)
	minDists := closedValvesMinDistances(valves, closed)
	closed.Remove(initValve)

	states := deque.New[State]() // BFS
	states.PushBack(State{
		valve:    initValve,
		opened:   set.Set[string]{},
		timeLeft: totalTime,
		pressure: 0,
	})

	maxPressure := 0
	for states.Len() > 0 {
		state := states.PopFront() // BFS
		fmt.Println(state, states.Len())

		if maxPressure < state.pressure {
			maxPressure = state.pressure
		}

		if state.timeLeft == 0 || state.opened.Len() == closedC {
			// Leaf state
			continue
		}

		// New state for every valve we can open
		for v := range closed.Diff(state.opened) {
			timeLeft := state.timeLeft - minDists[ValvePair{state.valve, v}] - 1
			if timeLeft < 0 {
				continue
			}

			newOpened := state.opened.Copy()
			newOpened.Add(v)
			states.PushBack(State{
				valve:    v,
				opened:   newOpened,
				timeLeft: timeLeft,
				pressure: state.pressure + valves[v].flowRate*timeLeft,
			})
		}
	}

	return maxPressure
}

func Part2(inputPath string) any {
	// Quite bad :(
	const (
		initValve = "AA"
		totalTime = 26
	)

	valves, closed := parseValves(inputPath)
	closedC := len(closed)

	closed.Add(initValve)
	minDists := closedValvesMinDistances(valves, closed)
	closed.Remove(initValve)

	states := deque.New[State]() // BFS
	states.PushBack(State{
		valve:    initValve,
		opened:   set.Set[string]{},
		timeLeft: totalTime,
		pressure: 0,
	})

	paths := hashext.NewHashMap[set.Set[string], int]()
	for states.Len() > 0 {
		state := states.PopFront() // BFS
		fmt.Println(state, states.Len())

		h := paths.Hash(state.opened)
		if v := paths.Vals[h]; v < state.pressure {
			paths.SetH(h, state.opened, state.pressure)
		}

		if state.timeLeft == 0 || state.opened.Len() == closedC {
			// Leaf state
			continue
		}

		// New state for every valve we can open
		for v := range closed.Diff(state.opened) {
			timeLeft := state.timeLeft - minDists[ValvePair{state.valve, v}] - 1
			if timeLeft < 0 {
				continue
			}

			newOpened := state.opened.Copy()
			newOpened.Add(v)
			states.PushBack(State{
				valve:    v,
				opened:   newOpened,
				timeLeft: timeLeft,
				pressure: state.pressure + valves[v].flowRate*timeLeft,
			})
		}
	}

	maxPressure := 0
	for leftH, leftP := range paths.Keys {
		for rightH, rightP := range paths.Keys {
			if !leftP.Disjoint(rightP) {
				continue
			}

			pressure := paths.Vals[leftH] + paths.Vals[rightH]
			if maxPressure < pressure {
				maxPressure = pressure
				fmt.Println(maxPressure)
			}
		}
	}

	return maxPressure
}
