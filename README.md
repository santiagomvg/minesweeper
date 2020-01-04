# minesweeper GO
Deviget coding challenge by Santiago Vidal

Online game: http://mine-sw79.herokuapp.com/

Minesweeper web game in the Go programming language.
Persistence for boards is managed in-memory with a browser cookie. Many users can play at the same time.
There's no user accounts implementation.

The board is represented in a two-dimensional array in json format.<br/>
Every cell defines two values. A constant (what to display) and a boolean (if there's a mine in that cell).

## Backend APIs
### POST /action
Triggers an action on the board. It can be a "clear cell" action or a "mark cell" one.<br/>
Returns a json array with the new state of the board. 
```json
{
  "action": "mark",
  "row": 1,
  "col": 2
}
```
```json
{
  "action": "clear",
  "row": 1,
  "col": 2
}
```

### POST /newGame
Resets the current game (cookie based) with the desired columns and rows count.<br/>
Returns a json array with the board.
```json
{
  "row": 15,
  "col": 15
}
```