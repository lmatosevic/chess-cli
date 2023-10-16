package game

import (
	"github.com/lmatosevic/chess-cli/pkg/utils"
	"testing"
)

func TestTilesParser(t *testing.T) {
	_, err := parseTiles(MakeStartingBoard())
	utils.AssertTestCondition(t, nil, err, "Tiles should be parsed without error")
}

func TestRegularMoveParser(t *testing.T) {
	_, err := parseMove("Pa2a4")
	utils.AssertTestCondition(t, nil, err, "Move should be parsed without error")
}

func TestRegularShortFileMoveParser(t *testing.T) {
	_, err := parseMove("Paa4")
	utils.AssertTestCondition(t, nil, err, "Move should be parsed without error")
}

func TestRegularShortRankMoveParser(t *testing.T) {
	_, err := parseMove("P2a4")
	utils.AssertTestCondition(t, nil, err, "Move should be parsed without error")
}

func TestCaptureMoveParser(t *testing.T) {
	_, err := parseMove("Nb1xc3")
	utils.AssertTestCondition(t, nil, err, "Capture move should be parsed without error")
}

func TestPromotionMoveParser(t *testing.T) {
	_, err := parseMove("Pa7a8Q")
	utils.AssertTestCondition(t, nil, err, "Promotion move should be parsed without error")
}

func TestKingCheckMoveParser(t *testing.T) {
	_, err := parseMove("Ra2a6+")
	utils.AssertTestCondition(t, nil, err, "King check move should be parsed without error")
}

func TestKingSideCastlingMoveParser(t *testing.T) {
	_, err := parseMove("0-0")
	utils.AssertTestCondition(t, nil, err, "King side castling move should be parsed without error")
}

func TestQueenSideCastlingMoveParser(t *testing.T) {
	_, err := parseMove("0-0-0")
	utils.AssertTestCondition(t, nil, err, "Queen side castling move should be parsed without error")
}

func TestMakeGame(t *testing.T) {
	_, err := MakeGame(MakeStartingBoard(), []string{"Pa2a4", "ng8f6", "Pc2c3"})
	utils.AssertTestCondition(t, nil, err, "Game should be made without error")
}
