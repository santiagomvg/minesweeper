package main

import (
	"github.com/pborman/uuid"
	"net/http"
	"time"
)

const gameCookieName string = "SV_minesweeper"

func getWebGame(w http.ResponseWriter, r *http.Request, forceNew bool) *game {

	c, err := r.Cookie(gameCookieName)
	if forceNew || err != nil || c == nil {
		return createNewWebGame(w, 10, 10, 3)
	} else {

		boardLock.RLock()
		defer boardLock.RUnlock()

		bg, exists := gameBoards[c.Value]
		if exists && bg.isValid() {
			return &bg
		} else {
			return createNewWebGame(w, 10, 10, 3)
		}
	}
}

func createNewWebGame(w http.ResponseWriter, rows int, cols int, mines int) *game {

	board := generateBoard(rows, cols, mines)
	newGame := game{
		startTime: time.Now(),
		limit:     time.Minute * 10,
		board:     board,
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
