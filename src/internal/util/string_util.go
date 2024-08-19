package util

import (
	"errors"
	"math/rand"
	"regexp"
	"strconv"
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
		if r == '"' && pre != backSlash {
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
			if len(piece) > 1 && piece[len(piece)-1] != doubleQuote { //the last must double quote too
				return false
			} else if len(piece) > 2 && piece[len(piece)-2] == backSlash { // but if the last double quote escaped, then wrong
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

var letterRunes = []rune("AaBb0CcDd1EeFf2GgHh3iIJj4KkLl5MmNn6OoPp7QqRr8SsTt9UuVvWwXxYyZz")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

/*
* Return: [host, port, port2], err
 */
func SplitAddress(address string) ([]string, error) {
	reg := regexp.MustCompile(`(.*)((?::))((?:[0-9]+))$`)
	pieces := reg.FindStringSubmatch(address)
	if len(pieces) != 4 {
		return nil, errors.New("Invalid Listen format: " + address)
	}
	iPort, err := strconv.Atoi(pieces[3])
	if err != nil {
		return nil, errors.New("Invalid Listen format: " + address)
	}
	iPort += 10000
	return []string{pieces[1], pieces[3], strconv.Itoa(iPort)}, nil
}
