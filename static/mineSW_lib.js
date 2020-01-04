const ZeroAdjacentMines = 0;
const OneAdjacentMines = 1;
const TwoAdjacentMines = 2;
const ThreeAdjacentMines = 3;
const FourAdjacentMines = 4;
const FiveAdjacentMines = 5;
const SixAdjacentMines = 6;
const SevenAdjacentMines = 7;
const EightAdjacentMines = 8;
const Uncleared = 9;
const UnclearedAndMarked = 10;
const Mine = 11; //for GameOver display

function MineSW() {
    //default callbacks
    this._errorCallback = function(e) { alert(e); };
    this._renderCallback = function(board) { alert(board); };
}

MineSW.prototype.setRenderCallback = function(cb) {
    this._renderCallback = cb;
};

MineSW.prototype.setErrorCallback = function(cb) {
    this._errorCallback = cb;
};

MineSW.prototype.markCell = function(row, col) {
    this._doAction('mark', row, col);
};

MineSW.prototype.clearCell = function(row, col) {
    this._doAction('clear', row, col);
};

MineSW.prototype.startGame = function(rowCount, colCount) {

    var self = this;
    $.ajax({
        type: "POST",
        url: "/newGame",
        data: JSON.stringify({ row: parseInt(rowCount), col: parseInt(colCount) }),
        contentType : 'application/json',
        success: function(data) { self._processBoard(data); }
    });
};

MineSW.prototype._doAction = function(action, row, col) {

    var self = this;
    $.ajax({
        type: "POST",
        url: "/action",
        data: JSON.stringify({ action: action, row: row, col: col }),
        contentType : 'application/json',
        success: function(data) { self._processBoard(data); }
    });
};

MineSW.prototype._processBoard = function(data) {
    this._renderCallback(JSON.parse(data));
};

MineSW.prototype._gameOver  = function() {

};