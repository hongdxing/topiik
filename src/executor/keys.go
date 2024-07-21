/***
** author: duan hongxing
** date: 29 Jun 2024
** desc:
**	return keys
**
**/
package executor

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"topiik/internal/consts"
	"topiik/memo"
)

/***
** Return keys
** Parameters:
**	- args: the arguments, command line that CMD(KEYS) stripped
** Return:
**	-
** Synctax: KEYS pattern
**	- pattern is a string to search keys, use astrisk(*) for pattern search
**/
func keys(pieces []string) (result []string, err error) {
	if len(pieces) != 1 {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	pattern := pieces[0]
	if !strings.HasPrefix(pattern, "*") { // exactly match from beginning
		pattern = "^" + pattern
	}
	if !strings.HasSuffix(pattern, "*") { // exactly match from endding
		pattern = pattern + "$"
	}
	fmt.Println(strings.ReplaceAll(pattern, "*", ".*"))
	reg, err := regexp.Compile(strings.ReplaceAll(pattern, "*", ".*"))
	if err != nil {
		return nil, errors.New(RES_SYNTAX_ERROR)
	}
	keys := make([]string, 0, len(memo.MemMap))
	for k := range memo.MemMap {
		// Need to exclude internal using KEYs
		if reg.MatchString(k) && !strings.HasPrefix(k, consts.RESEVERD_PREFIX) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}
