package main

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

//every game board has a browser's key cookie associated. This is the in-memory persistence and its synchronized by boardLock
type gameGrid [][]cell
type cellAttr byte

var gameBoards map[string]game
var boardLock sync.RWMutex

type cell struct {
	hasMine bool
	flags   cellAttr
}
type game struct {
	startTime time.Time
	limit     time.Duration
	board     gameGrid
}

func (g game) isValid() bool {
	return g.startTime.Add(g.limit).After(time.Now())
}

func (g game) stream(out io.Writer) {
	clBoard := clientBoard{board: g.board, expiresIn: g.limit}
	if err := json.NewEncoder(out).Encode(clBoard); err != nil {
		panic(err)
	}
}

//every cell board has a byte witch identifies its content. It's defined by the following constants
//I could've used a bitwise approach but this seems cleaner to the client.
const ZeroAdjacentMines cellAttr = 0
const OneAdjacentMines cellAttr = 1
const TwoAdjacentMines cellAttr = 2
const ThreeAdjacentMines cellAttr = 3
const FourAdjacentMines cellAttr = 4
const FiveAdjacentMines cellAttr = 5
const SixAdjacentMines cellAttr = 6
const SevenAdjacentMines cellAttr = 7
const EightAdjacentMines cellAttr = 8
const Uncleared cellAttr = 9
const UnclearedAndMarked cellAttr = 10

//this is the board sent to the client. it has no mines information to avoid cheating
type clientBoard struct {
	expiresIn time.Duration
	board     gameGrid
}

func handleGameAction(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func handleRestartAction(w http.ResponseWriter, r *http.Request) error {
	getWebGame(w, r, true).stream(w)
	return nil
}

func generateBoard(rows int, cols int, mines int) gameGrid {

	if rows <= 0 || cols <= 0 {
		panic("grid dimensions must be positive.")
	}
	grid := make([][]cell, rows)
	for r := 0; r < rows; r++ {
		grid[r] = make([]cell, cols)
		for c := 0; c < cols; c++ {
			grid[r][c].flags = Uncleared
		}
	}
	min := int(math.Round(float64(rows*cols) * 0.1))
	max := int(math.Round(float64(rows*cols) * 0.2))
	mineCount := min + rand.Intn(max-min+1)
	rm := mineCount
	for rm > 0 {
		x, y := rand.Intn(rows), rand.Intn(cols)
		if !grid[x][y].hasMine {
			rm--
			grid[x][y].hasMine = true
		}
	}
	return gameGrid(grid)
}
