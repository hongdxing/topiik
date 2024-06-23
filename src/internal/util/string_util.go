package util

import (
	"errors"
	"strings"
)

const (
	space       = ' '
	doubleQuote = '"'
	backSlash   = '\\'
)

func SplitCommandLine(str string) (pieces []string, err error) {
	quoted := false
	pre := space
	pieces = strings.FieldsFunc(str, func(r rune) bool {
		if r == '"' && pre != '\\' {
			quoted = !quoted
		}
		pre = r
		return !quoted && r == space
	})
	isValid := ValidateCommandLinePieces(&pieces)
	if isValid {
		return pieces, nil
	}
	return nil, errors.New("INVALID_QUOTATION")
}

func ValidateCommandLinePieces(pieces *[]string) bool {
	var quoted bool
	for i, piece := range *pieces {
		quoted = false
		// check if piece is double quoted
		piece = strings.TrimSpace(piece)
		if piece[0] == doubleQuote {
			if piece[len(piece)-1] != doubleQuote { //the last must double quote too
				return false
			} else if piece[len(piece)-2] == backSlash { // but if the last double quote escaped, then wrong
				return false
			}
			quoted = true
		}

		// check if double quote in the middle of piece
		idx := strings.Index(piece, `"`)
		if idx > 0 && idx < len(piece)-1 {
			if piece[idx-1] != backSlash {
				return false
			}
		}

		if quoted {
			(*pieces)[i] = strings.TrimRight(strings.TrimLeft(strings.TrimSpace(piece), `"`), `"`)
		}
	}
	return true
}
