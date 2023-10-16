package game

import (
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"testing"
)

func TestWhitePawnForwardMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Pawn), "a", "2")

	srcRow, srcCol := boardRowAndCol("a", "2")
	destRow, destCol := boardRowAndCol("a", "3")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, true, "")
	utils.AssertTestCondition(t, nil, err, "Pawn forward move should be valid")
}

func TestWhitePawnForwardDoubleMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Pawn), "a", "2")

	srcRow, srcCol := boardRowAndCol("a", "2")
	destRow, destCol := boardRowAndCol("a", "4")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, true, "")
	utils.AssertTestCondition(t, nil, err, "Pawn forward double move should be valid")
}

func TestWhitePawnCaptureRightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Pawn), "a", "2")
	addFigureToBoard(board, BlackFigure(Pawn), "b", "3")

	srcRow, srcCol := boardRowAndCol("a", "2")
	destRow, destCol := boardRowAndCol("b", "3")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, true, "")
	utils.AssertTestCondition(t, nil, err, "Pawn capture right move should be valid")
}

func TestWhitePawnCaptureLeftMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Pawn), "c", "2")
	addFigureToBoard(board, BlackFigure(Pawn), "b", "3")

	srcRow, srcCol := boardRowAndCol("c", "2")
	destRow, destCol := boardRowAndCol("b", "3")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, true, "")
	utils.AssertTestCondition(t, nil, err, "Pawn capture left move should be valid")
}

func TestWhitePawnPromotionMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Pawn), "a", "7")

	srcRow, srcCol := boardRowAndCol("a", "7")
	destRow, destCol := boardRowAndCol("a", "8")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, true, WhiteFigure(Queen))
	utils.AssertTestCondition(t, nil, err, "Pawn promotion move should be valid")
}

func TestBlackPawnForwardMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(Pawn), "a", "7")

	srcRow, srcCol := boardRowAndCol("a", "7")
	destRow, destCol := boardRowAndCol("a", "6")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, false, "")
	utils.AssertTestCondition(t, nil, err, "Pawn forward move should be valid")
}

func TestBlackPawnForwardDoubleMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(Pawn), "a", "7")

	srcRow, srcCol := boardRowAndCol("a", "7")
	destRow, destCol := boardRowAndCol("a", "5")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, false, "")
	utils.AssertTestCondition(t, nil, err, "Pawn forward double move should be valid")
}

func TestBlackPawnCaptureRightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(Pawn), "a", "7")
	addFigureToBoard(board, WhiteFigure(Pawn), "b", "6")

	srcRow, srcCol := boardRowAndCol("a", "7")
	destRow, destCol := boardRowAndCol("b", "6")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, false, "")
	utils.AssertTestCondition(t, nil, err, "Pawn capture right move should be valid")
}

func TestBlackPawnCaptureLeftMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(Pawn), "c", "7")
	addFigureToBoard(board, WhiteFigure(Pawn), "b", "6")

	srcRow, srcCol := boardRowAndCol("c", "7")
	destRow, destCol := boardRowAndCol("b", "6")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, false, "")
	utils.AssertTestCondition(t, nil, err, "Pawn capture left move should be valid")
}

func TestBlackPawnPromotionMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(Pawn), "a", "2")

	srcRow, srcCol := boardRowAndCol("a", "2")
	destRow, destCol := boardRowAndCol("a", "1")
	err := validatePawnsMove(board, srcRow, srcCol, destRow, destCol, false, BlackFigure(Queen))
	utils.AssertTestCondition(t, nil, err, "Pawn promotion move should be valid")
}

func TestKnightLongMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Knight), "b", "1")

	srcRow, srcCol := boardRowAndCol("b", "1")
	destRow, destCol := boardRowAndCol("a", "3")
	err := validateKnightsMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Knight long move should be valid")
}

func TestKnightShortMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Knight), "b", "1")

	srcRow, srcCol := boardRowAndCol("b", "1")
	destRow, destCol := boardRowAndCol("d", "2")
	err := validateKnightsMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Knight short move should be valid")
}

func TestBishopRightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Bishop), "c", "1")

	srcRow, srcCol := boardRowAndCol("c", "1")
	destRow, destCol := boardRowAndCol("g", "5")
	err := validateBishopsMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Bishop right move should be valid")
}

func TestBishopLeftMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Bishop), "c", "1")

	srcRow, srcCol := boardRowAndCol("c", "1")
	destRow, destCol := boardRowAndCol("a", "3")
	err := validateBishopsMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Bishop left move should be valid")
}

func TestRookUpMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Rook), "a", "1")

	srcRow, srcCol := boardRowAndCol("a", "1")
	destRow, destCol := boardRowAndCol("a", "8")
	err := validateRooksMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Rook up move should be valid")
}

func TestRookRightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Rook), "a", "1")

	srcRow, srcCol := boardRowAndCol("a", "1")
	destRow, destCol := boardRowAndCol("h", "1")
	err := validateRooksMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Rook right move should be valid")
}

func TestKingRightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")

	srcRow, srcCol := boardRowAndCol("e", "1")
	destRow, destCol := boardRowAndCol("f", "1")
	err := validateKingsMove(board, srcRow, srcCol, destRow, destCol, true)
	utils.AssertTestCondition(t, nil, err, "King right move should be valid")
}

func TestKingLeftMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")

	srcRow, srcCol := boardRowAndCol("e", "1")
	destRow, destCol := boardRowAndCol("d", "1")
	err := validateKingsMove(board, srcRow, srcCol, destRow, destCol, true)
	utils.AssertTestCondition(t, nil, err, "King left move should be valid")
}

func TestKingUpMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")

	srcRow, srcCol := boardRowAndCol("e", "1")
	destRow, destCol := boardRowAndCol("e", "2")
	err := validateKingsMove(board, srcRow, srcCol, destRow, destCol, true)
	utils.AssertTestCondition(t, nil, err, "King up move should be valid")
}

func TestKingDownMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "4")

	srcRow, srcCol := boardRowAndCol("e", "4")
	destRow, destCol := boardRowAndCol("e", "3")
	err := validateKingsMove(board, srcRow, srcCol, destRow, destCol, true)
	utils.AssertTestCondition(t, nil, err, "King down move should be valid")
}

func TestKingDiagonalMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "4")

	srcRow, srcCol := boardRowAndCol("e", "4")
	destRow, destCol := boardRowAndCol("f", "5")
	err := validateKingsMove(board, srcRow, srcCol, destRow, destCol, true)
	utils.AssertTestCondition(t, nil, err, "King diagonal move should be valid")
}

func TestQueenStraightMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Queen), "d", "1")

	srcRow, srcCol := boardRowAndCol("d", "1")
	destRow, destCol := boardRowAndCol("d", "7")
	err := validateQueensMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Queen straight move should be valid")
}

func TestQueenDiagonalMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(Queen), "d", "1")

	srcRow, srcCol := boardRowAndCol("d", "1")
	destRow, destCol := boardRowAndCol("h", "5")
	err := validateQueensMove(board, srcRow, srcCol, destRow, destCol)
	utils.AssertTestCondition(t, nil, err, "Queen diagonal move should be valid")
}

func TestWhiteKingSideCastlingMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")
	addFigureToBoard(board, WhiteFigure(Rook), "h", "1")

	err := validateCastlingMove(board, true, true, &[]Move{})
	utils.AssertTestCondition(t, nil, err, "King side castling move should be valid")
}

func TestWhiteQueenSideCastlingMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")
	addFigureToBoard(board, WhiteFigure(Rook), "a", "1")

	err := validateCastlingMove(board, false, true, &[]Move{})
	utils.AssertTestCondition(t, nil, err, "Queen side castling move should be valid")
}

func TestBlackKingSideCastlingMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(King), "e", "8")
	addFigureToBoard(board, BlackFigure(Rook), "h", "8")

	err := validateCastlingMove(board, true, false, &[]Move{})
	utils.AssertTestCondition(t, nil, err, "King side castling move should be valid")
}

func TestBlackQueenSideCastlingMove(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(King), "e", "8")
	addFigureToBoard(board, BlackFigure(Rook), "a", "8")

	err := validateCastlingMove(board, false, false, &[]Move{})
	utils.AssertTestCondition(t, nil, err, "Queen side castling move should be valid")
}

func TestWhiteKingCheck(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, WhiteFigure(King), "e", "1")
	addFigureToBoard(board, BlackFigure(Bishop), "h", "4")

	isCheck := IsKingCheck(board, true)
	utils.AssertTestCondition(t, true, isCheck, "The white king should be in check")
}

func TestBlackKingCheck(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(King), "e", "8")
	addFigureToBoard(board, WhiteFigure(Knight), "d", "6")

	isCheck := IsKingCheck(board, false)
	utils.AssertTestCondition(t, true, isCheck, "The black king should be in check")
}

func TestGameWon(t *testing.T) {
	board := makeEmptyBoard()
	addFigureToBoard(board, BlackFigure(King), "g", "8")
	addFigureToBoard(board, BlackFigure(Rook), "f", "8")
	addFigureToBoard(board, BlackFigure(Pawn), "g", "7")
	addFigureToBoard(board, WhiteFigure(Queen), "h", "7")
	addFigureToBoard(board, WhiteFigure(Pawn), "g", "6")
	addFigureToBoard(board, WhiteFigure(King), "g", "1")

	win := IsGameWon(board, true)
	utils.AssertTestCondition(t, true, win, "The white player should have won the game")
}

func addFigureToBoard(board *Board, figure string, file string, rank string) {
	col := BoardFileToColumn(file)
	row := BoardRankToRow(rank)
	board[row][col] = figure
}

func boardRowAndCol(file string, rank string) (int, int) {
	return BoardRankToRow(rank), BoardFileToColumn(file)
}

func makeEmptyBoard() *Board {
	board := Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			board[i][j] = Empty
		}
	}
	return &board
}
