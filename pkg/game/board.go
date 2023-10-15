package game

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	CaptureSign          = "x"
	KingCheckSign        = "+"
	CheckmateSign        = "#"
	DrawOfferMove        = "="
	DrawOfferRejectMove  = "!"
	KingSideCastligMove  = "0-0"
	QueenSideCastligMove = "0-0-0"
)

const (
	Rook   = "R"
	Knight = "N"
	Bishop = "B"
	King   = "K"
	Queen  = "Q"
	Pawn   = "P"
	Empty  = "0"
)

func MakeStartingBoard() string {
	firstLayer := fmt.Sprintf("%s%s%s%s%s%s%s%s", Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook)

	secondLayer := ""
	for i := 0; i < 8; i += 1 {
		secondLayer += Pawn
	}

	emptyLayers := ""
	for i := 0; i < 32; i += 1 {
		emptyLayers += Empty
	}

	return fmt.Sprintf("%s%s%s%s%s",
		BlackFigure(firstLayer), BlackFigure(secondLayer), emptyLayers, WhiteFigure(secondLayer), WhiteFigure(firstLayer))
}

func WhiteFigure(char string) string {
	return strings.ToUpper(char)
}

func BlackFigure(char string) string {
	return strings.ToLower(char)
}

func ColoredFigure(char string, isWhite bool) string {
	if isWhite {
		return WhiteFigure(char)
	} else {
		return BlackFigure(char)
	}
}

func IsValidFigure(char string) bool {
	lc := strings.ToUpper(char)
	if lc == Rook || lc == Knight || lc == Bishop || lc == King || lc == Queen || lc == Pawn || lc == Empty {
		return true
	} else {
		return false
	}
}

func IsFigureType(char string, fig string) bool {
	return strings.ToUpper(char) == strings.ToUpper(fig)
}

func IsPlayersFigure(char string, isWhite bool) bool {
	return char == ColoredFigure(char, isWhite) && char != Empty
}

func BoardFileToColumn(char string) int {
	ascii := int(rune(char[0]))
	return ascii - 97
}

func BoardColumnToFile(col int) string {
	return fmt.Sprintf("%c", col+97)
}

func BoardRankToRow(char string) int {
	index, _ := strconv.Atoi(char)
	return 8 - index
}

func BoardRowToRank(row int) string {
	return fmt.Sprintf("%d", 8-row)
}
