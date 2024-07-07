/***
** OS utils
**
**
**
**/

package util

import (
	"os"
	"path/filepath"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetMainPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	mainPath := filepath.Dir(ex)
	return mainPath
}
