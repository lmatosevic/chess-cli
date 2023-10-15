package game

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type Board [8][8]string

type Move struct {
	Figure              string
	FigureFile          string
	FigureRank          string
	DestinationFile     string
	DestinationRank     string
	PromotedToFigure    string
	IsCapture           bool
	IsKingSideCastling  bool
	IsQueenSideCastling bool
	IsKingCheck         bool
}

type Game struct {
	Board Board
	Moves []Move
}

var moveRegex = regexp.MustCompile(
	fmt.Sprintf("^(\\w)([a-h])?([1-8])?(%s)?([a-h])([1-8])(\\w)?([%s%s])?$", CaptureSign, KingCheckSign, CheckmateSign))

// MakeMove godoc
// The move parameter represents the moving of a single figure on board and should be in following format:
// (figure_char)(file)?(rank)?(capture_sign)?(file)(rank)(promoted_figure_char)(king_check_sign|checkmate_sign)?
//
// The x represents that this move captures opponents figure and + at the end that
// this move is king check for opponent. (e.g. Qa3, Nxf3, Bxc4+, R3xa6+)
//
// The move can also be a request for draw by containing only = (equals sign) or rejection of draw ! (exclamation mark)
//
// This function returns following values (normalized move, is game won, error)
func (g *Game) MakeMove(move string, isWhite bool) (string, bool, error) {
	if slices.Contains([]string{DrawOfferMove, DrawOfferRejectMove}, move) {
		return move, false, nil
	}

	m, err := parseMove(move)
	if err != nil {
		return "", false, err
	}

	err = ValidateMove(&g.Board, m, isWhite, &g.Moves)
	if err != nil {
		c := "white"
		if !isWhite {
			c = "black"
		}
		return "", false, errors.New(fmt.Sprintf(`Invalid move "%s" for %s player. Reason: %s`, move, c, err.Error()))
	}

	ExecuteMove(&g.Board, m, isWhite)

	moveStr := m.String()
	isWin := IsGameWon(&g.Board, isWhite)
	if isWin {
		if m.IsKingCheck {
			moveStr = strings.Replace(moveStr, KingCheckSign, CheckmateSign, 1)
		} else {
			moveStr = fmt.Sprintf("%s#", moveStr)
		}
	}

	return m.String(), IsGameWon(&g.Board, isWhite), nil
}

func (g *Game) GetTiles() string {
	tiles := ""
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			tiles = fmt.Sprintf("%s%s", tiles, g.Board[i][j])
		}
	}
	return tiles
}

func (m *Move) String() string {
	if m.IsKingSideCastling {
		return KingSideCastligMove
	}
	if m.IsKingSideCastling {
		return QueenSideCastligMove
	}

	isCapture := ""
	if m.IsCapture {
		isCapture = CaptureSign
	}

	isKingCheck := ""
	if m.IsKingCheck {
		isKingCheck = KingCheckSign
	}

	return fmt.Sprintf("%s%s%s%s%s%s%s%s", m.Figure, m.FigureFile, m.FigureRank, isCapture, m.DestinationFile,
		m.DestinationRank, m.PromotedToFigure, isKingCheck)
}

func MakeGame(tiles string, moves []string) (*Game, error) {
	board, err := parseTiles(tiles)
	if err != nil {
		return nil, err
	}

	movesList := make([]Move, 0)
	for _, m := range moves {
		if m == DrawOfferMove || m == DrawOfferRejectMove {
			continue
		}
		move, e := parseMove(m)
		if e != nil {
			return nil, e
		}
		movesList = append(movesList, *move)
	}
	return &Game{*board, movesList}, nil
}

func parseMove(move string) (*Move, error) {
	// 1 -> figure, 2 -> figure file, 3 -> figure rank, 4 -> capture char, 5 -> dest file, 6 -> dest rank,
	// 7 -> promoted figure, 8 -> king check or endgame mark
	matches := moveRegex.FindStringSubmatch(move)
	if len(matches) < 8 {
		if move == KingSideCastligMove {
			return &Move{IsKingSideCastling: true}, nil
		}

		if move == QueenSideCastligMove {
			return &Move{IsQueenSideCastling: true}, nil
		}

		return nil, errors.New(fmt.Sprintf("invalid move format: %s", move))
	}

	if !IsValidFigure(matches[1]) {
		return nil, errors.New(fmt.Sprintf("invalid figure character: %s", matches[1]))
	}

	isCapture := false
	if matches[4] == CaptureSign {
		isCapture = true
	}

	promotedToFigure := ""
	if matches[7] != "" {
		if !IsFigureType(matches[1], Pawn) {
			return nil, errors.New(fmt.Sprintf("invalid figure for promotion: %s", matches[1]))
		}
		if IsFigureType(matches[7], Pawn) || IsFigureType(matches[7], King) {
			return nil, errors.New(fmt.Sprintf("promoting pawn to the invalid figure: %s", matches[7]))
		}
		promotedToFigure = matches[7]
	}

	isKingCheck := false
	if slices.Contains([]string{KingCheckSign, CheckmateSign}, matches[8]) {
		isKingCheck = true
	}

	m := Move{Figure: matches[1], FigureFile: matches[2], FigureRank: matches[3], DestinationFile: matches[5],
		DestinationRank: matches[6], PromotedToFigure: promotedToFigure, IsCapture: isCapture, IsKingCheck: isKingCheck}

	return &m, nil
}

func parseTiles(tiles string) (*Board, error) {
	if len(tiles) != 64 {
		return nil, errors.New("required number of tiles is 64")
	}
	board := Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			char := string(tiles[j+(i*8)])
			if !IsValidFigure(char) {
				return nil, errors.New(fmt.Sprintf("invalid figure character: %s", char))
			}
			board[i][j] = char
		}
	}
	return &board, nil
}
