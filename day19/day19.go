package day19

import (
	"fmt"
	"github.com/gammazero/deque"
	"github.com/tobiasbrandy/AoC_2022_go/internal/mathext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/parse"
	"github.com/tobiasbrandy/AoC_2022_go/internal/regext"
	"github.com/tobiasbrandy/AoC_2022_go/internal/stringer"
	"regexp"

	"github.com/tobiasbrandy/AoC_2022_go/internal/errexit"
	"github.com/tobiasbrandy/AoC_2022_go/internal/fileline"
)

type Blueprint struct {
	oreOreCost,
	clayOreCost,
	obsidianOreCost,
	obsidianClayCost,
	geodeOreCost,
	geodeObsidianCost int
}

func (b Blueprint) String() string {
	return stringer.String(b)
}

var blueprintParser = regexp.MustCompile(`Blueprint \d+:` +
	` Each ore robot costs (?P<oreOreCost>\d+) ore.` +
	` Each clay robot costs (?P<clayOreCost>\d+) ore.` +
	` Each obsidian robot costs (?P<obsidianOreCost>\d+) ore and (?P<obsidianClayCost>\d+) clay.` +
	` Each geode robot costs (?P<geodeOreCost>\d+) ore and (?P<geodeObsidianCost>\d+) obsidian.` +
	``)

func parseBlueprint(data string) Blueprint {
	blueprintInfo := regext.NamedCaptureGroups(blueprintParser, data)
	return Blueprint{
		oreOreCost:        parse.Int(blueprintInfo["oreOreCost"]),
		clayOreCost:       parse.Int(blueprintInfo["clayOreCost"]),
		obsidianOreCost:   parse.Int(blueprintInfo["obsidianOreCost"]),
		obsidianClayCost:  parse.Int(blueprintInfo["obsidianClayCost"]),
		geodeOreCost:      parse.Int(blueprintInfo["geodeOreCost"]),
		geodeObsidianCost: parse.Int(blueprintInfo["geodeObsidianCost"]),
	}
}

type State struct {
	ore, clay, obsidian, geode,
	oreR, clayR, obsidianR, geodeR,
	timeLeft int
}

func (s State) String() string {
	return stringer.String(s)
}

// blueprintMaxGeodes finds the max geodes that could be generated using DFS with pruning optimizations.
//
// Pruning optimizations:
// - Pruning a branch if its best possible score (crafting one geode robot per minute until the end) is worse than the current best branch.
// - Building a geode robot immediately whenever they are available.
// - Forbidding the crafting of robots of a type if we are already generating the maximum needed of that material per minute for any recipe.
// - Only choose to wait if it allows building a new type of robot.
//
// Other potential pruning optimizations:
// - Ending a branch early if there is no possible way to generate enough obsidian to craft any more geode robots.
// - Pruning a branch if it tries to create a robot of a certain type, when it could have been created in the previous state too, but it decided to do nothing instead.
func blueprintMaxGeodes(b Blueprint, initTime, initOreR int) int {
	maxGeode := 0 // This is what we are looking for

	maxOreCost := mathext.IntMax(b.oreOreCost, b.clayOreCost, b.obsidianOreCost, b.geodeOreCost)
	maxClayCost := b.obsidianClayCost
	maxObsidianCost := b.obsidianClayCost

	var states deque.Deque[State]
	states.PushBack(State{
		oreR:     initOreR,
		timeLeft: initTime,
	})

	for states.Len() > 0 {
		state := states.PopBack() // DFS

		// Mine ores
		ore := state.ore + state.oreR
		clay := state.clay + state.clayR
		obsidian := state.obsidian + state.obsidianR
		geode := state.geode + state.geodeR
		timeLeft := state.timeLeft - 1

		if maxGeode < geode {
			maxGeode = geode
		}

		maxGeodesPossible := geode + timeLeft*state.geodeR + timeLeft*(timeLeft+1)/2
		if timeLeft == 0 || maxGeodesPossible <= maxGeode {
			// End state
			continue
		}

		// Geode Robot
		// Max priority: if we can build it, we do it. We don't analyze other options. Tree pruning.
		if state.ore >= b.geodeOreCost && state.obsidian >= b.geodeObsidianCost {
			states.PushBack(State{
				ore:       ore - b.obsidianOreCost,
				clay:      clay,
				obsidian:  obsidian - b.geodeObsidianCost,
				geode:     geode,
				oreR:      state.oreR,
				clayR:     state.clayR,
				obsidianR: state.obsidianR,
				geodeR:    state.geodeR + 1,
				timeLeft:  timeLeft,
			})
			continue
		}

		// Only wait if waiting allows unlocking a new robot
		waitingUnlocks := false

		// Obsidian Robot
		if state.ore >= b.obsidianOreCost && state.clay >= b.obsidianClayCost {
			if state.obsidianR < maxObsidianCost {
				states.PushBack(State{
					ore:       ore - b.obsidianOreCost,
					clay:      clay - b.obsidianClayCost,
					obsidian:  obsidian,
					geode:     geode,
					oreR:      state.oreR,
					clayR:     state.clayR,
					obsidianR: state.obsidianR + 1,
					geodeR:    state.geodeR,
					timeLeft:  timeLeft,
				})
			}
		} else if state.clayR > 0 {
			waitingUnlocks = true
		}

		// Clay Robot
		if state.ore >= b.clayOreCost {
			if state.clayR < maxClayCost {
				states.PushBack(State{
					ore:       ore - b.clayOreCost,
					clay:      clay,
					obsidian:  obsidian,
					geode:     geode,
					oreR:      state.oreR,
					clayR:     state.clayR + 1,
					obsidianR: state.obsidianR,
					geodeR:    state.geodeR,
					timeLeft:  timeLeft,
				})
			}
		} else {
			waitingUnlocks = true
		}

		// Ore Robot
		if state.ore >= b.oreOreCost {
			if state.oreR < maxOreCost {
				states.PushBack(State{
					ore:       ore - b.oreOreCost,
					clay:      clay,
					obsidian:  obsidian,
					geode:     geode,
					oreR:      state.oreR + 1,
					clayR:     state.clayR,
					obsidianR: state.obsidianR,
					geodeR:    state.geodeR,
					timeLeft:  timeLeft,
				})
			}
		} else {
			waitingUnlocks = true
		}

		// Wait
		if waitingUnlocks {
			states.PushBack(State{
				ore:       ore,
				clay:      clay,
				obsidian:  obsidian,
				geode:     geode,
				oreR:      state.oreR,
				clayR:     state.clayR,
				obsidianR: state.obsidianR,
				geodeR:    state.geodeR,
				timeLeft:  timeLeft,
			})
		}
	}

	return maxGeode
}

func Part1(inputPath string) any {
	const (
		initTime int = 24
		initOreR int = 1
	)

	total := 0
	id := 1
	fileline.ForEach(inputPath, errexit.HandleScanError, func(line string) {
		maxGeode := blueprintMaxGeodes(parseBlueprint(line), initTime, initOreR)
		fmt.Println(maxGeode)
		total += id * maxGeode
		id++
	})

	return total
}

func Part2(inputPath string) any {
	const (
		initTime       int = 32
		initOreR       int = 1
		blueprintCount int = 3
	)

	total := 1
	fileline.ForEachN(inputPath, errexit.HandleScanError, blueprintCount, func(line string) {
		maxGeode := blueprintMaxGeodes(parseBlueprint(line), initTime, initOreR)
		fmt.Println(maxGeode)
		total *= maxGeode
	})

	return total
}
