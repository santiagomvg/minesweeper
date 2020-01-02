package main

import (
	"net/http"
	"time"
)

func getWebGame(r *http.Request) *game {

	c, err := r.Cookie("SV_minesweeper")
	if err != nil || c == nil {
		board := generateBoard(10, 10, 1)
		return &game{
			startTime: time.Now(),
			limit:     time.Minute * 10,
			board:     board,
		}
	} else {

		boardLock.RLock()
		defer boardLock.RUnlock()

		bg, exists := gameBoards[c.Value]
		if exists {
			return &bg
		} else {
			board := generateBoard(10, 10, 1)
			return &game{
				startTime: time.Now(),
				limit:     time.Minute * 10,
				board:     board,
			}
		}

	}
}
