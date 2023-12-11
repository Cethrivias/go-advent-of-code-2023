package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

var checked = []Coords{}
var tiles [][]TileType

func Solve(path string) (total int) {
	input := getInput(path)

	for _, line := range input.Map {
		println(string(line))
	}

	loop := Loop{}
	directions := getDirections(input.Map, input.Start)
	for _, direction := range directions {
		currLoop, ok := move(input.Map, direction, Loop{})
		if ok {
			total = len(currLoop) / 2
			if len(currLoop)%2 != 0 {
				total++
			}
			loop = currLoop
			break
		}
	}

	// Part 2

	tiles = make([][]TileType, len(input.Map))
	for i, row := range input.Map {
		tiles[i] = make([]TileType, len(row))
	}

	// go through tiles 1 by 1
	// skip if loop
	// mark as regular if next to border or have other regular nearby
	//    check if has enclosed nearby then mark them as regular too
	// mark everything else as enclosed

	for i := range tiles {
		for j := range tiles[i] {
			// println()
			// printMap(input.Map, loop)
			// time.Sleep(250 * time.Millisecond)
			if isLoop(input.Map, loop, Coords{i, j}) {
				tiles[i][j] = LoopTile
				continue
			}
			if isRegular(input.Map, Coords{i, j}) {
				tiles[i][j] = RegularTile
				spreadRegularity(input.Map, Coords{i, j})
				continue
			}
			tiles[i][j] = EnclosedTile
		}
	}

	// Check if possible to squeeze through the pipes
	println()
	printMap(input.Map, loop)
	time.Sleep(250 * time.Millisecond)

	for i := range tiles {
		for j, tileType := range tiles[i] {
			if tileType == EnclosedTile {
				if canEscape(input.Map, loop, Coords{i, j}) {
					spreadRegularity(input.Map, Coords{i, j})
				}
			}
		}
	}

	total = 0
	for i := range tiles {
		for j := range tiles[i] {
			if tiles[i][j] == EnclosedTile {
				total++
			}
		}
	}

	return total
}

func canEscape(m Map, l Loop, c Coords) bool {
	// fmt.Printf("Checking: %+v\n", c)
	// below
	cCurr := c.Bottom()
	if m[cCurr.i][cCurr.j] == '7' { // && tiles[cCurr.Right().i][cCurr.Right().j] == LoopTile {
		println("here")
		wallA := []Coords{cCurr}
		wallB := []Coords{cCurr.Right()}
		cCurr := cCurr.Bottom()
		cOther := cCurr.Right()
		wallA = append(wallA, cCurr)
		wallB = append(wallB, cOther)
		for {
			directions := getDirections(m, cCurr)
			directionsA := []Coords{}
			for _, direction := range directions {
				if !isLoop(m, l, direction) {
					continue
				}
				if sliceContains(wallA, direction) {
					continue
				}
				directionsA = append(directionsA, direction)
			}
			if len(directionsA) != 1 {
				log.Fatal("WTF A")
			}

			directions = getDirections(m, cOther)
			directionsB := []Coords{}
			for _, direction := range directions {
				if !isLoop(m, l, direction) {
					continue
				}
				if sliceContains(wallB, direction) {
					continue
				}
				directionsB = append(directionsB, direction)
			}
			if len(directionsA) != 1 {
				log.Fatal("WTF B")
			}

			cCurr = directionsA[0]
			wallA = append(wallA, cCurr)
			cOther = directionsB[0]
			wallB = append(wallB, cOther)

			fmt.Printf("A: %+v; B: %+v\n", directionsA[0], directionsB[0])
			if math.Abs(float64(directionsA[0].i-directionsB[0].i)) > 1 ||
				math.Abs(float64(directionsA[0].j-directionsB[0].j)) > 1 {
				fmt.Printf("Walls diverge %+v : %+v\n", directionsA, directionsB)
			}
		}
	}

	return false
}

func sliceContains(slice []Coords, c Coords) bool {
	for _, it := range slice {
		if it == c {
			return true
		}
	}
	return false
}

func isLoop(m Map, l Loop, c Coords) bool {
	if m[c.i][c.j] == 'S' {
		return true
	}
	for _, cLoop := range l {
		if c == cLoop {
			return true
		}
	}
	return false
}

func isRegular(m Map, c Coords) bool {
	for _, cCurr := range []Coords{c.Top(), c.Right(), c.Bottom(), c.Left()} {
		if !m.withinBounds(cCurr) || tiles[cCurr.i][cCurr.j] == RegularTile {
			return true
		}
	}
	return false
}

func spreadRegularity(m Map, c Coords) {
	for _, cCurr := range []Coords{c.Top(), c.Right(), c.Bottom(), c.Left()} {
		if m.withinBounds(cCurr) && tiles[cCurr.i][cCurr.j] == EnclosedTile {
			tiles[cCurr.i][cCurr.j] = RegularTile
			spreadRegularity(m, cCurr)
		}
	}
}

func move(m Map, c Coords, loop Loop) (Loop, bool) {
	if m[c.i][c.j] == 'S' {
		return loop, true
	}
	for _, cVisited := range loop {
		if c.i == cVisited.i && c.j == cVisited.j {
			return loop, false
		}
	}

	loop = append(loop, c)

	directions := getDirections(m, c)
	for _, direction := range directions {
		loop, ok := move(m, direction, loop)
		if ok {
			return loop, ok
		}
	}

	return loop, false
}

func printMap(m Map, l Loop) {
	mCopy := make(Map, len(m))

	for i := range m {
		mCopy[i] = make([]rune, len(m[i]))
		for j, tile := range m[i] {
			switch tiles[i][j] {
			case RegularTile:
				mCopy[i][j] = 'O'
			case EnclosedTile:
				mCopy[i][j] = 'I'
			default:
				mCopy[i][j] = tile
			}
		}
	}

	for _, row := range mCopy {
		println("> ", string(row))
	}
}

func getDirections(m Map, c Coords) []Coords {
	res := []Coords{}

	switch m[c.i][c.j] {
	case 'S':
		if c, ok := m.TopConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.RightConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.BottomConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.LeftConnected(c); ok {
			res = append(res, c)
		}
	case '-':
		if c, ok := m.RightConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.LeftConnected(c); ok {
			res = append(res, c)
		}
	case '|':
		if c, ok := m.TopConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.BottomConnected(c); ok {
			res = append(res, c)
		}
	case 'L':
		if c, ok := m.TopConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.RightConnected(c); ok {
			res = append(res, c)
		}
	case 'F':
		if c, ok := m.RightConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.BottomConnected(c); ok {
			res = append(res, c)
		}
	case '7':
		if c, ok := m.BottomConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.LeftConnected(c); ok {
			res = append(res, c)
		}
	case 'J':
		if c, ok := m.TopConnected(c); ok {
			res = append(res, c)
		}
		if c, ok := m.LeftConnected(c); ok {
			res = append(res, c)
		}
	default:
		fmt.Printf("Unknown symbol: '%s'", string(m[c.i][c.j]))
	}

	return res
}

func getInput(path string) *Input {
	input := Input{}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {
		text := scanner.Text()

		input.Map = append(input.Map, []rune(text))

		j := strings.Index(text, "S")
		if j != -1 {
			input.Start = Coords{i: i, j: j}
		}
		i++
	}

	return &input
}

type Input struct {
	Map   Map
	Start Coords
}

type Map [][]rune

func (m Map) TopConnected(c Coords) (Coords, bool) {
	return m.hasConnection(c.Top(), "|F7S")
}

func (m Map) RightConnected(c Coords) (Coords, bool) {
	return m.hasConnection(c.Right(), "-J7S")
}

func (m Map) BottomConnected(c Coords) (Coords, bool) {
	return m.hasConnection(c.Bottom(), "|JLS")
}

func (m Map) LeftConnected(c Coords) (Coords, bool) {
	return m.hasConnection(c.Left(), "-FLS")
}

func (m Map) hasConnection(c Coords, pipes string) (Coords, bool) {
	return c, m.withinBounds(c) && strings.Contains(pipes, string(m[c.i][c.j]))
}

func (m Map) withinBounds(c Coords) bool {
	return c.i >= 0 && c.i < len(m) && c.j >= 0 && c.j < len(m[c.i])
}

type Coords struct {
	i int
	j int
}

func (c *Coords) Top() Coords {
	return Coords{i: c.i - 1, j: c.j}
}
func (c *Coords) Right() Coords {
	return Coords{i: c.i, j: c.j + 1}
}
func (c *Coords) Bottom() Coords {
	return Coords{i: c.i + 1, j: c.j}
}
func (c *Coords) Left() Coords {
	return Coords{i: c.i, j: c.j - 1}
}

type Loop []Coords

type Tile struct {
	Type   bool
	Coords Coords
}

type TileType string

var (
	RegularTile  TileType = "regular_tile"
	EnclosedTile TileType = "enclosed_tile"
	LoopTile     TileType = "loop_tile"
)
