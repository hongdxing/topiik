//Author: Duan Hongxing
//Date: 6 Jul, 2024

package util

import (
	"errors"
	"os"
	"path/filepath"
	"topiik/internal/consts"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
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

func GetDataFullPath() string {
	return GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH
}
