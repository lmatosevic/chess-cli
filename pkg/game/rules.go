package game

import (
	"errors"
	"math"
	"strings"
)

func ValidateMove(board *Board, move *Move, isWhite bool, moveHistory *[]Move) error {
	if move.IsKingSideCastling || move.IsQueenSideCastling {
		return validateCastlingMove(board, move.IsKingSideCastling, isWhite, moveHistory)
	}

	destRow, destCol, figureRow, figureCol := destAndFigurePositions(board, move, isWhite)
	if figureRow == -1 || figureCol == -1 {
		return errors.New("cannot uniquely identify figure on board")
	}
	if destRow == figureRow && destCol == figureCol {
		return errors.New("destination must be different than the current figure position")
	}
	if !IsPlayersFigure(board[figureRow][figureCol], isWhite) {
		return errors.New("not players figure")
	}

	err := validateFigureMove(board, move, figureRow, figureCol, destRow, destCol, isWhite)

	if err == nil && board[destRow][destCol] != Empty {
		if IsPlayersFigure(board[destRow][destCol], isWhite) {
			return errors.New("cannot capture own figure")
		}
		if IsFigureType(board[destRow][destCol], King) {
			return errors.New("cannot capture king figure")
		}
	}

	return err
}

func ExecuteMove(board *Board, move *Move, isWhite bool) {
	if move.IsKingSideCastling || move.IsQueenSideCastling {
		kingsCol := 4
		kingsRow := 0
		if isWhite {
			kingsRow = 7
		}

		kingFigure := board[kingsRow][kingsCol]
		if move.IsKingSideCastling {
			rookCol := 7
			board[kingsRow][kingsCol+2] = kingFigure
			board[kingsRow][kingsCol] = Empty
			board[kingsRow][rookCol-2] = board[kingsRow][rookCol]
			board[kingsRow][rookCol] = Empty
		} else {
			rookCol := 0
			board[kingsRow][kingsCol-2] = kingFigure
			board[kingsRow][kingsCol] = Empty
			board[kingsRow][rookCol+3] = board[kingsRow][rookCol]
			board[kingsRow][rookCol] = Empty
		}

		return
	}

	destRow, destCol, figureRow, figureCol := destAndFigurePositions(board, move, isWhite)

	if move.FigureRank == "" {
		move.FigureRank = BoardRowToRank(figureRow)
	}
	if move.FigureFile == "" {
		move.FigureFile = BoardColumnToFile(figureCol)
	}

	if board[destRow][destCol] != Empty {
		move.IsCapture = true
	}

	move.Figure = ColoredFigure(move.Figure, isWhite)

	if move.PromotedToFigure != "" && IsFigureType(move.Figure, Pawn) && (destRow == 0 || destRow == 7) {
		board[destRow][destCol] = ColoredFigure(move.PromotedToFigure, isWhite)
	} else {
		board[destRow][destCol] = board[figureRow][figureCol]
	}
	board[figureRow][figureCol] = Empty

	move.IsKingCheck = IsKingCheck(board, !isWhite)
}

func IsGameWon(board *Board, isWhite bool) bool {
	if !IsKingCheck(board, !isWhite) {
		return false
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			opponentFigure := board[i][j]
			if !IsPlayersFigure(opponentFigure, !isWhite) {
				continue
			}

			for k := 0; k < 8; k++ {
				for m := 0; m < 8; m++ {
					if IsPlayersFigure(board[k][m], !isWhite) {
						continue
					}

					move := Move{Figure: opponentFigure, PromotedToFigure: Queen}
					err := validateFigureMove(board, &move, i, j, k, m, !isWhite)
					if err == nil && !willKingBeInCheck(board, i, j, k, m, !isWhite) {
						return false
					}
				}
			}
		}
	}

	return true
}

func IsKingCheck(board *Board, isWhite bool) bool {
	kingRow, kingCol := findFigureRowAndColumn(board, King, "", "", isWhite)
	if kingRow == -1 || kingCol == -1 {
		return false
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			figure := board[i][j]
			if IsPlayersFigure(figure, !isWhite) {
				move := Move{Figure: figure, PromotedToFigure: Queen}
				err := validateFigureMove(board, &move, i, j, kingRow, kingCol, !isWhite)
				if err == nil {
					return true
				}
			}
		}
	}

	return false
}

func validateFigureMove(board *Board, move *Move, figureRow int, figureCol int, destRow int, destCol int, isWhite bool) error {
	figure := strings.ToUpper(move.Figure)

	var err error
	switch figure {
	case Pawn:
		err = validatePawnsMove(board, figureRow, figureCol, destRow, destCol, isWhite, move.PromotedToFigure)
	case Knight:
		err = validateKnightsMove(board, figureRow, figureCol, destRow, destCol)
	case Bishop:
		err = validateBishopsMove(board, figureRow, figureCol, destRow, destCol)
	case Rook:
		err = validateRooksMove(board, figureRow, figureCol, destRow, destCol)
	case King:
		err = validateKingsMove(board, figureRow, figureCol, destRow, destCol, isWhite)
	case Queen:
		err = validateQueensMove(board, figureRow, figureCol, destRow, destCol)
	default:
		err = errors.New("unknown figure")
	}
	return err
}

func validatePawnsMove(board *Board, figureRow int, figureCol int, destRow int, destCol int, isWhite bool, promoteToFigure string) error {
	if isWhite {
		if destRow == 0 && promoteToFigure == "" {
			return errors.New("move to this position must be a promotion")
		}

		if destRow-figureRow > 0 {
			return errors.New("pawn can only move forward")
		}

		if (figureRow == 6 && destRow-figureRow < -2) || (figureRow < 6 && destRow-figureRow > -1) {
			return errors.New("pawn can only move 1 tile forward (or 2 for its first move)")
		}

		if destCol != figureCol {
			if math.Abs(float64(destCol-figureCol)) != 1 || destRow-figureRow != -1 {
				return errors.New("pawn can capture only diagonally adjacent tiles")
			}

			if board[destRow][destCol] == Empty {
				return errors.New("pawn can only move diagonally to capture other players figure")
			}
		} else {
			if board[destRow][destCol] != Empty || (destRow-figureRow == -2 && board[destRow+1][destCol] != Empty) {
				return errors.New("pawn can only move forward through empty tiles")
			}
		}
	} else {
		if destRow == 7 && promoteToFigure == "" {
			return errors.New("move to this position must be a promotion")
		}

		if destRow-figureRow < 0 {
			return errors.New("pawn can only move forward")
		}

		if (figureRow == 1 && destRow-figureRow > 2) || (figureRow > 1 && destRow-figureRow > 1) {
			return errors.New("pawn can only move 1 tile forward (or 2 for its first move)")
		}

		if destCol != figureCol {
			if math.Abs(float64(destCol-figureCol)) != 1 || destRow-figureRow != 1 {
				return errors.New("pawn can capture only diagonally adjacent tiles")
			}

			if board[destRow][destCol] == Empty {
				return errors.New("pawn can only move diagonally to capture other players figure")
			}
		} else {
			if board[destRow][destCol] != Empty || (destRow-figureRow == 2 && board[destRow-1][destCol] != Empty) {
				return errors.New("pawn can only move forward through empty tiles")
			}
		}
	}

	return nil
}

func validateKnightsMove(board *Board, figureRow int, figureCol int, destRow int, destCol int) error {
	rowDiff, colDiff := rowAndColDiffs(figureRow, figureCol, destRow, destCol)

	if (rowDiff != 1 && rowDiff != 2) || (colDiff != 1 && colDiff != 2) || (rowDiff == colDiff) {
		return errors.New("knight can only move in L shape")
	}

	return nil
}

func validateBishopsMove(board *Board, figureRow int, figureCol int, destRow int, destCol int) error {
	rowDiff, colDiff := rowAndColDiffs(figureRow, figureCol, destRow, destCol)

	if rowDiff != colDiff {
		return errors.New("bishop can only move diagonally")
	}

	rowDir, colDir := rowAndColDirections(figureRow, figureCol, destRow, destCol)

	for i := 1; i < rowDiff; i++ {
		if board[figureRow+i*rowDir][figureCol+i*colDir] != Empty {
			return errors.New("bishop can only move through empty tiles")
		}
	}

	return nil
}

func validateRooksMove(board *Board, figureRow int, figureCol int, destRow int, destCol int) error {
	rowDiff, colDiff := rowAndColDiffs(figureRow, figureCol, destRow, destCol)

	if rowDiff > 0 && colDiff > 0 {
		return errors.New("rook can only move in straight lines")
	}

	rowDir, colDir := rowAndColDirections(figureRow, figureCol, destRow, destCol)

	distance := rowDiff
	if rowDiff == 0 {
		distance = colDiff
	}

	for i := 1; i < distance; i++ {
		if board[figureRow+i*rowDir][figureCol+i*colDir] != Empty {
			return errors.New("rook can only move through empty tiles")
		}
	}

	return nil
}

func validateKingsMove(board *Board, figureRow int, figureCol int, destRow int, destCol int, isWhite bool) error {
	rowDiff, colDiff := rowAndColDiffs(figureRow, figureCol, destRow, destCol)

	if rowDiff > 1 || colDiff > 1 {
		return errors.New("king can only move for 1 tile in any direction")
	}

	if willKingBeInCheck(board, figureRow, figureCol, destRow, destCol, isWhite) {
		return errors.New("cannot move king to position where it is under check")
	}

	return nil
}

func validateQueensMove(board *Board, figureRow int, figureCol int, destRow int, destCol int) error {
	rowDiff, colDiff := rowAndColDiffs(figureRow, figureCol, destRow, destCol)

	if rowDiff != colDiff && (rowDiff != 0 && colDiff != 0) {
		return errors.New("queen can move diagonally and in straight lines")
	}

	rowDir, colDir := rowAndColDirections(figureRow, figureCol, destRow, destCol)

	distance := int(math.Max(float64(rowDiff), float64(colDiff)))

	for i := 1; i < distance; i++ {
		if board[figureRow+i*rowDir][figureCol+i*colDir] != Empty {
			return errors.New("queen can only move through empty tiles")
		}
	}

	return nil
}

func validateCastlingMove(board *Board, isKingSide bool, isWhite bool, moveHistory *[]Move) error {
	for _, m := range *moveHistory {
		if m.IsKingSideCastling || m.IsQueenSideCastling {
			return errors.New("castling has already been done")
		}
		if IsFigureType(m.Figure, King) {
			return errors.New("cannot castle anymore because king has been moved")
		}
		if IsFigureType(m.Figure, Rook) && ((isKingSide && m.FigureFile == "h") || (!isKingSide && m.FigureFile == "a")) {
			return errors.New("cannot castle anymore because rook has been moved")
		}
	}

	kingsCol := 4
	kingsRow := 0
	if isWhite {
		kingsRow = 7
	}

	colsToCheck := []int{5, 6}
	if !isKingSide {
		colsToCheck = []int{1, 2, 3}
	}

	willBeCheck := false
	for i := 0; i < len(colsToCheck); i++ {
		if board[kingsRow][colsToCheck[i]] != Empty {
			return errors.New("cannot castle while there are figures between the king and the rook")
		}
		if !willBeCheck && willKingBeInCheck(board, kingsRow, kingsCol, kingsRow, colsToCheck[i], isWhite) {
			willBeCheck = true
		}
	}

	if willBeCheck {
		return errors.New("cannot castle because at least one tile between king and rook are in check")
	}

	if IsKingCheck(board, isWhite) {
		return errors.New("cannot castle while king is in check")
	}

	return nil
}

func willKingBeInCheck(board *Board, figureRow int, figureCol int, destRow int, destCol int, isWhite bool) bool {
	tempBoard := Board{}
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			tempBoard[i][j] = board[i][j]
		}
	}

	tempBoard[destRow][destCol] = tempBoard[figureRow][figureCol]
	tempBoard[figureRow][figureCol] = Empty

	return IsKingCheck(&tempBoard, isWhite)
}

func destAndFigurePositions(board *Board, move *Move, isWhite bool) (int, int, int, int) {
	destRow, destCol := BoardRankToRow(move.DestinationRank), BoardFileToColumn(move.DestinationFile)
	figureRow, figureCol := findFigureRowAndColumn(board, move.Figure, move.FigureFile, move.FigureRank, isWhite)
	return destRow, destCol, figureRow, figureCol
}

func rowAndColDiffs(figureRow int, figureCol int, destRow int, destCol int) (int, int) {
	return int(math.Abs(float64(figureRow - destRow))), int(math.Abs(float64(figureCol - destCol)))
}

func rowAndColDirections(figureRow int, figureCol int, destRow int, destCol int) (int, int) {
	rowDirection := 0
	if figureRow-destRow < 0 {
		rowDirection = 1
	} else if figureRow-destRow > 0 {
		rowDirection = -1
	}

	colDirection := 0
	if figureCol-destCol < 0 {
		colDirection = 1
	} else if figureCol-destCol > 0 {
		colDirection = -1
	}

	return rowDirection, colDirection
}

func findFigureRowAndColumn(board *Board, figure string, file string, rank string, isWhite bool) (int, int) {
	figureRow, figureCol := -1, -1
	if file != "" {
		figureCol = BoardFileToColumn(file)
	}
	if rank != "" {
		figureRow = BoardRankToRow(rank)
	}

	playerFigure := ColoredFigure(figure, isWhite)

	if figureCol > -1 && figureRow > -1 {
		boardFigure := board[figureRow][figureCol]
		if boardFigure == playerFigure {
			return figureRow, figureCol
		}
		return -1, -1
	}

	if figureCol > -1 {
		for i := 0; i < 8; i++ {
			boardFigure := board[i][figureCol]
			if boardFigure == playerFigure {
				return i, figureCol
			}
		}
		return -1, -1
	}

	if figureRow > -1 {
		for i := 0; i < 8; i++ {
			boardFigure := board[figureRow][i]
			if boardFigure == playerFigure {
				return figureRow, i
			}
		}
		return -1, -1
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			boardFigure := board[i][j]
			if boardFigure == playerFigure {
				if figureRow != -1 && figureCol != -1 {
					return -1, -1
				}
				figureRow = i
				figureCol = j
			}
		}
	}

	return figureRow, figureCol
}
