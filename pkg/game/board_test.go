package game

import (
	"gitlab.com/lmatosevic/chess-cli/pkg/utils"
	"testing"
)

func TestBoardSetup(t *testing.T) {
	board := MakeStartingBoard()
	expectedBoard := "rnbqkbnrpppppppp00000000000000000000000000000000PPPPPPPPRNBQKBNR"
	utils.AssertTestCondition(t, expectedBoard, board, "Invalid board setup")
}

func TestWhiteFigureString(t *testing.T) {
	utils.AssertTestCondition(t, "P", WhiteFigure(Pawn), "Invalid pawn white figure")
	utils.AssertTestCondition(t, "N", WhiteFigure(Knight), "Invalid knight white figure")
	utils.AssertTestCondition(t, "B", WhiteFigure(Bishop), "Invalid bishop white figure")
	utils.AssertTestCondition(t, "R", WhiteFigure(Rook), "Invalid rook white figure")
	utils.AssertTestCondition(t, "K", WhiteFigure(King), "Invalid king white figure")
	utils.AssertTestCondition(t, "Q", WhiteFigure(Queen), "Invalid queen white figure")
}

func TestBlackFigureString(t *testing.T) {
	utils.AssertTestCondition(t, "p", BlackFigure(Pawn), "Invalid pawn black figure")
	utils.AssertTestCondition(t, "n", BlackFigure(Knight), "Invalid knight black figure")
	utils.AssertTestCondition(t, "b", BlackFigure(Bishop), "Invalid bishop black figure")
	utils.AssertTestCondition(t, "r", BlackFigure(Rook), "Invalid rook black figure")
	utils.AssertTestCondition(t, "k", BlackFigure(King), "Invalid king black figure")
	utils.AssertTestCondition(t, "q", BlackFigure(Queen), "Invalid queen black figure")
}

func TestFigureValidity(t *testing.T) {
	utils.AssertTestCondition(t, true, IsValidFigure(Rook), "The figure R should be valid")
	utils.AssertTestCondition(t, false, IsValidFigure("X"), "The figure X should be invalid")
}

func TestFigureType(t *testing.T) {
	utils.AssertTestCondition(t, true, IsFigureType(Pawn, BlackFigure(Pawn)), "The pawn figure should be of same type")
	utils.AssertTestCondition(t, false, IsFigureType(Pawn, BlackFigure(Rook)), "The pawn figure should not be of same type")
}

func TestPlayersFigure(t *testing.T) {
	utils.AssertTestCondition(t, true, IsPlayersFigure(WhiteFigure(Pawn), true), "The white pawn figure should be from white player")
	utils.AssertTestCondition(t, false, IsPlayersFigure(BlackFigure(Pawn), true), "The black pawn figure should not be from white player")
}

func TestBoardRankToRowMapping(t *testing.T) {
	utils.AssertTestCondition(t, 3, BoardRankToRow("5"), "Board rank 5 should map to row 3")
	utils.AssertTestCondition(t, 0, BoardRankToRow("8"), "Board rank 8 should map to row 0")
	utils.AssertTestCondition(t, 7, BoardRankToRow("1"), "Board rank 1 should map to row 7")
}

func TestBoardRowToRankMapping(t *testing.T) {
	utils.AssertTestCondition(t, "5", BoardRowToRank(3), "Board row 3 should map to row 5")
	utils.AssertTestCondition(t, "8", BoardRowToRank(0), "Board row 0 should map to row 8")
	utils.AssertTestCondition(t, "1", BoardRowToRank(7), "Board row 7 should map to row 1")
}

func TestBoardFileToColumnMapping(t *testing.T) {
	utils.AssertTestCondition(t, 3, BoardFileToColumn("d"), "Board file d should map to column 3")
	utils.AssertTestCondition(t, 0, BoardFileToColumn("a"), "Board file a should map to column 0")
	utils.AssertTestCondition(t, 7, BoardFileToColumn("h"), "Board file h should map to column 7")
}

func TestBoardColumnToFileMapping(t *testing.T) {
	utils.AssertTestCondition(t, "d", BoardColumnToFile(3), "Board column 3 should map to file d")
	utils.AssertTestCondition(t, "a", BoardColumnToFile(0), "Board column 0 should map to file a")
	utils.AssertTestCondition(t, "h", BoardColumnToFile(7), "Board column 7should map to file h")
}
