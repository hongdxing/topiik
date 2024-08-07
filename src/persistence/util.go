/***
** author: duan hongxing
** date: 29 Jun 2024
** desc:
**	Util functions
**
**/

package persistence

import "os"

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
