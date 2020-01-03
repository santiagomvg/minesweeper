function MineSW() {
    //default callbacks
    this._errorCallback = function(e) { alert(e); };
    this._renderCallback = function(board) { alert(board); };
}

MineSW.prototype.setRenderCallback = function(cb) {

};

MineSW.prototype.setErrorCallback = function(cb) {

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
        data: { row: rowCount, col: colCount },
        success: function(data) { self._processBoard(data); }
    });
};

MineSW.prototype._doAction = function(action, row, col) {

    var self = this;
    $.ajax({
        type: "POST",
        url: "/action",
        data: { action: action, row: row, col: col },
        success: function(data) { self._processBoard(data); }
    });
};

MineSW.prototype._processBoard = function(data) {
    this._renderCallback(data);
};

MineSW.prototype._gameOver  = function() {

};