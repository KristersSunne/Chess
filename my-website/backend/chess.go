package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Piece int

const (
    Empty Piece = iota
    Pawn
    Knight
    Bishop
    Rook
    Queen
    King
)

type Color int
const (
    White Color = iota
    Black
)

type Square struct {
    Piece Piece
    Color Color
}

type Board [8][8]Square

func (p Piece) String() string {
	switch p {
	case Pawn:
		return "P"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	case King:
		return "K"
	default:
		return "."
	}
}

func NewBoard() Board {
	var b Board

	for i := 0; i < 8; i++ {
		b[1][i] = Square{Pawn, White}
		b[6][i] = Square{Pawn, Black}
	}

	b[0][0] = Square{Rook, White}
	b[0][7] = Square{Rook, White}
	b[7][0] = Square{Rook, Black}
	b[7][7] = Square{Rook, Black}

	b[0][1] = Square{Knight, White}
	b[0][6] = Square{Knight, White}
	b[7][1] = Square{Knight, Black}
	b[7][6] = Square{Knight, Black}

	b[0][2] = Square{Bishop, White}
	b[0][5] = Square{Bishop, White}
	b[7][2] = Square{Bishop, Black}
	b[7][5] = Square{Bishop, Black}

	b[0][3] = Square{Queen, White}
	b[7][3] = Square{Queen, Black}

	b[0][4] = Square{King, White}
	b[7][4] = Square{King, Black}

	return b
}

func printBoard(b Board) {
    var i = 0
    for i < len(b) {
        var j = 0
        for j < len(b[0]) {
            fmt.Print("[", b[i][j].Piece, ",", b[i][j].Color, "]")
            j++
        }
        fmt.Println(" ")
        i++
    }
    fmt.Println(" ")
}

func movePiece(b Board, xFrom, yFrom, xTo, yTo int) Board {
	selected := b[xFrom][yFrom]

	b[xFrom][yFrom] = Square{} 
	b[xTo][yTo] = selected

	return b
}

var reader = bufio.NewReader(os.Stdin)

func input() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func sameColor(b Board, fx, fy, tx, ty int) bool {
	return b[tx][ty].Piece != Empty &&
		b[fx][fy].Color == b[tx][ty].Color
}

func inBounds(x, y int) bool {
	return x >= 0 && x < 8 && y >= 0 && y < 8
}

func isValidMove(b Board, fromX, fromY, toX, toY int) bool {
	sq := b[fromX][fromY]

	if sq.Piece == Empty {
		return false
	}

	switch sq.Piece {
	case Pawn:
		return validPawnMove(b, fromX, fromY, toX, toY)
	case Knight:
		return validKnightMove(b, fromX, fromY, toX, toY)
	case Bishop:
		return validBishopMove(b, fromX, fromY, toX, toY)
	case Rook:
		return validRookMove(b, fromX, fromY, toX, toY)
	case Queen:
		return validQueenMove(b, fromX, fromY, toX, toY)
	case King:
		return validKingMove(b, fromX, fromY, toX, toY)
	}

	return false
}

func validBishopMove(b Board, fx, fy, tx, ty int) bool {
	if !inBounds(tx, ty) {
		return false
	}

	dx := fx - tx
	if dx < 0 {
		dx = -dx
	}
	dy := fy - ty
	if dy < 0 {
		dy = -dy
	}

	if dx != dy {
		return false
	}

	if sameColor(b, fx, fy, tx, ty) {
		return false
	}

	stepX := 1
	if tx < fx {
		stepX = -1
	}
	stepY := 1
	if ty < fy {
		stepY = -1
	}

	x := fx + stepX
	y := fy + stepY

	for x != tx && y != ty {
		if b[x][y].Piece != Empty {
			return false
		}
		x += stepX
		y += stepY
	}

	return true
}

func validKingMove(b Board, fx, fy, tx, ty int) bool {
	if !inBounds(tx, ty) {
		return false
	}

	dx := fx - tx
	if dx < 0 {
		dx = -dx
	}
	dy := fy - ty
	if dy < 0 {
		dy = -dy
	}

	if dx > 1 || dy > 1 {
		return false
	}

	if sameColor(b, fx, fy, tx, ty) {
		return false
	}

	next := movePiece(b, fx, fy, tx, ty)

	color := b[fx][fy].Color
	if isKingInCheck(next, color) {
		return false
	}

	return true
}

func validKnightMove(b Board, fx, fy, tx, ty int) bool {
	if !inBounds(tx, ty) {
		return false
	}

	dx := fx - tx
	if dx < 0 {
		dx = -dx
	}
	dy := fy - ty
	if dy < 0 {
		dy = -dy
	}

	if !((dx == 2 && dy == 1) || (dx == 1 && dy == 2)) {
		return false
	}

	return !sameColor(b, fx, fy, tx, ty)
}

func validPawnMove(b Board, fx, fy, tx, ty int) bool {
	if !inBounds(tx, ty) {
		return false
	}

	pawn := b[fx][fy]
	dir := 1
	startRow := 1

	if pawn.Color == Black {
		dir = -1
		startRow = 6
	}

	if fy == ty && b[tx][ty].Piece == Empty {
		if tx == fx+dir {
			return true
		}
		if fx == startRow && tx == fx+2*dir && b[fx+dir][fy].Piece == Empty {
			return true
		}
	}

	if tx == fx+dir && (ty == fy+1 || ty == fy-1) {
		return b[tx][ty].Piece != Empty &&
			b[tx][ty].Color != pawn.Color
	}

	return false
}

func validQueenMove(b Board, fx, fy, tx, ty int) bool {
	return validRookMove(b, fx, fy, tx, ty) ||
		validBishopMove(b, fx, fy, tx, ty)
}


func validRookMove(b Board, fx, fy, tx, ty int) bool {
	if !inBounds(tx, ty) {
		return false
	}

	if fx != tx && fy != ty {
		return false
	}

	if sameColor(b, fx, fy, tx, ty) {
		return false
	}

	dx := 0
	dy := 0

	if fx < tx {
		dx = 1
	} else if fx > tx {
		dx = -1
	}

	if fy < ty {
		dy = 1
	} else if fy > ty {
		dy = -1
	}

	x := fx + dx
	y := fy + dy

	for x != tx || y != ty {
		if b[x][y].Piece != Empty {
			return false
		}
		x += dx
		y += dy
	}

	return true
}

type AttackMap [8][8]bool

func markAttacks(b Board, fx, fy int, attacked *AttackMap) {
	switch b[fx][fy].Piece {
	case Pawn:
		markPawnAttacks(b, fx, fy, attacked)
	case Knight:
		markKnightAttacks(fx, fy, attacked)
	case Bishop:
		markSlidingAttacks(b, fx, fy, attacked, bishopDirs)
	case Rook:
		markSlidingAttacks(b, fx, fy, attacked, rookDirs)
	case Queen:
		markSlidingAttacks(b, fx, fy, attacked, queenDirs)
	case King:
		markKingAttacks(fx, fy, attacked)
	}
}

var rookDirs = [][2]int{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
}

var bishopDirs = [][2]int{
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

var queenDirs = append(rookDirs, bishopDirs...)

func markSlidingAttacks(b Board, fx, fy int, attacked *AttackMap, dirs [][2]int) {
	for _, d := range dirs {
		x := fx + d[0]
		y := fy + d[1]

		for inBounds(x, y) {
			attacked[x][y] = true

			if b[x][y].Piece != Empty {
				break
			}

			x += d[0]
			y += d[1]
		}
	}
}

var knightOffsets = [][2]int{
	{2, 1}, {1, 2}, {-1, 2}, {-2, 1},
	{-2, -1}, {-1, -2}, {1, -2}, {2, -1},
}

func markKnightAttacks(fx, fy int, attacked *AttackMap) {
	for _, o := range knightOffsets {
		x := fx + o[0]
		y := fy + o[1]
		if inBounds(x, y) {
			attacked[x][y] = true
		}
	}
}

func markPawnAttacks(b Board, fx, fy int, attacked *AttackMap) {
	dir := 1
	if b[fx][fy].Color == Black {
		dir = -1
	}

	for _, dy := range []int{-1, 1} {
		x := fx + dir
		y := fy + dy
		if inBounds(x, y) {
			attacked[x][y] = true
		}
	}
}

func markKingAttacks(fx, fy int, attacked *AttackMap) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			x := fx + dx
			y := fy + dy
			if inBounds(x, y) {
				attacked[x][y] = true
			}
		}
	}
}

type kindLocation [8][8]Square

func findKing(b Board, c Color) (int, int) {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if b[x][y].Piece == King && b[x][y].Color == c {
				return x, y
			}
		}
	}
	panic("king not found")
}

func isKingInCheck(b Board, color Color) bool {
	kx, ky := findKing(b, color)

	enemy := White
	if color == White {
		enemy = Black
	}

	attacked := buildAttackMap(b, enemy)
	return attacked[kx][ky]
}

func buildAttackMap(b Board, by Color) AttackMap {
	var attacked AttackMap

	for x:= 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if b[x][y].Piece != Empty && b[x][y].Color == by {
				markAttacks(b, x, y, &attacked)
			}
		}
	}
	return attacked
}

func printAttackMap(a AttackMap) {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if a[x][y] {
				fmt.Print("[x]")
			} else {
				fmt.Print("[ ]")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
    b := NewBoard()
    printBoard(b)
    i := 0
    for i < 5 {
        x_from_str := input()
        x_from, err := strconv.Atoi(x_from_str)
        y_from_str := input()
        y_from, err := strconv.Atoi(y_from_str)
        x_to_str := input()
        x_to, err := strconv.Atoi(x_to_str)
        y_to_str := input()
        y_to, err := strconv.Atoi(y_to_str)

        if err != nil {
            fmt.Println("Invalid number")
            return
        }

        if isValidMove(b, x_from, y_from, x_to, y_to) {
            b = movePiece(b, x_from, y_from, x_to, y_to)
        } else {
            fmt.Println("Invalid move")
        }

        printBoard(b)
		a0 := buildAttackMap(b, White)
		a1 := buildAttackMap(b, Black)
		printAttackMap(a0)
		printAttackMap(a1)
    }
}
