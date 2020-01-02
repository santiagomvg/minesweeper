package ms

import (
	"encoding/json"
	"errors"
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
	gameOver  bool
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

func (g game) clearCell(row int, col int, userAction bool) {

	cell := &g.board[row][col]
	if cell.flags == Uncleared && !cell.hasMine {

		mines := g.getSurroundingMines(row, col)
		cellFlags := &g.board[row][col].flags
		switch mines {
		case 0:
			*cellFlags = ZeroAdjacentMines //no adjacent mines. let's keep cleaning
			g.clearCell(row+1, col, false)
			g.clearCell(row+1, col+1, false)
			g.clearCell(row, col+1, false)
			g.clearCell(row-1, col+1, false)
			g.clearCell(row-1, col, false)
			g.clearCell(row-1, col-1, false)
			g.clearCell(row, col-1, false)
			g.clearCell(row+1, col-1, false)

		case 1:
			*cellFlags = OneAdjacentMines
		case 2:
			*cellFlags = TwoAdjacentMines
		case 3:
			*cellFlags = ThreeAdjacentMines
		case 4:
			*cellFlags = FourAdjacentMines
		case 5:
			*cellFlags = FiveAdjacentMines
		case 6:
			*cellFlags = SixAdjacentMines
		case 7:
			*cellFlags = SevenAdjacentMines
		case 8:
			*cellFlags = EightAdjacentMines
		}

	} else if userAction && cell.hasMine {
		g.endGame()
	}
}

func (g game) getSurroundingMines(row int, col int) int {

	count := 0
	for y := row - 1; y <= row+1; y++ {
		if y >= 0 && y < len(g.board) {
			for x := col - 1; x <= col+1; x++ {
				if x >= 0 && x < len(g.board[y]) {
					if g.board[x][y].hasMine {
						count++
					}
				}
			}
		}
	}
	return count
}

func (g game) endGame() {
	g.gameOver = true
	for _, row := range g.board {
		for _, col := range row {
			if col.hasMine {
				col.flags = Mine //display flag
			}
		}
	}
}

func (g game) markCell(row int, col int) {
	c := g.board[row][col]
	if c.flags == Uncleared {
		c.flags = UnclearedAndMarked
	} else if c.flags == UnclearedAndMarked {
		c.flags = Uncleared
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
const Mine cellAttr = 11 //for GameOver display

//this is the board sent to the client. it has no mines information to avoid cheating
type clientBoard struct {
	expiresIn time.Duration
	board     gameGrid
}

func HandleGameAction(w http.ResponseWriter, r *http.Request) error {

	g := getWebGame(w, r, false)

	var data inputData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return err
	}

	if err := data.action.validate(); err != nil {
		return err
	}

	switch data.action {
	case "mark":
		g.markCell(data.row, data.col)
	case "clear":
		g.clearCell(data.row, data.col, true)
	}
	g.stream(w)
	return nil
}

func HandleRestartAction(w http.ResponseWriter, r *http.Request) error {
	getWebGame(w, r, true).stream(w)
	return nil
}

func generateBoard(rows int, cols int, mines int) (*gameGrid, error) {

	if rows <= 0 || cols <= 0 {
		return nil, errors.New("invalid grid dimensions")
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
	ret := gameGrid(grid)
	return &ret, nil
}