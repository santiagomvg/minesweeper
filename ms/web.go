package ms

import (
	"encoding/json"
	"github.com/pborman/uuid"
	"net/http"
	"time"
)

const gameCookieName string = "SV_minesweeper"

type initData struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func getWebGame(w http.ResponseWriter, r *http.Request, forceNew bool) *game {

	var input initData
	if forceNew {
		json.NewDecoder(r.Body).Decode(&input)
	} else {
		input.Col = 10
		input.Row = 10
	}

	c, err := r.Cookie(gameCookieName)
	if forceNew || err != nil || c == nil {
		return createNewWebGame(w, input.Row, input.Col, 3)
	} else {

		boardLock.RLock()
		defer boardLock.RUnlock()

		bg, exists := gameBoards[c.Value]
		if exists && bg.isValid() {
			return &bg
		} else {
			return createNewWebGame(w, input.Row, input.Col, 3)
		}
	}
}

func createNewWebGame(w http.ResponseWriter, rows int, cols int, mines int) *game {

	board, err := generateBoard(rows, cols, mines)
	if err != nil {
		panic(err)
	}
	newGame := game{
		startTime: time.Now(),
		limit:     time.Minute * 10,
		board:     *board,
	}

	cookieValue := uuid.New()
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{
		Name:    gameCookieName,
		Value:   cookieValue,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)

	boardLock.Lock()
	defer boardLock.Unlock()
	gameBoards[cookieValue] = newGame
	return &newGame
}
