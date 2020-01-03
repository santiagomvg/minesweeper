function MineSW() {}

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
        data: { row: rowCount, col: colCount },
        success: self._processBoard()
    });
};

MineSW.prototype._doAction = function(action, row, col) {

    var self = this;
    $.ajax({
        type: "POST",
        url: "/action",
        data: { action: action, row: row, col: col },
        success: self._processBoard()
    });
};

MineSW.prototype._processBoard = function(data) {

};

MineSW.prototype._gameOver  = function() {

};