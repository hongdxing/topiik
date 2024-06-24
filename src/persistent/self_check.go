/***
** author: duan hongxing
** date: 24 Jun 2024
** desc:
**
**/

package persistent

import (
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"
)

func SelfCheck() (err error) {
	fmt.Printf("Self check start\n")
	dataDir := "data"
	nodeFile := path.Join(dataDir, string(os.PathSeparator), "node")
	var hasDataDir bool
	hasDataDir, err = exists("data")
	if err != nil {
		return err
	}

	if !hasDataDir {
		fmt.Println("Creating data dir...")
		err = os.Mkdir(dataDir, os.ModeDir)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Creating node file...")
		var file *os.File
		file, err = os.Create(nodeFile)
		if err != nil {
			// if node file create failed, remove data dir too
			os.Remove(dataDir)
		}
		nodeId, err := uuid.NewUUID()
		if err != nil {
			// remove both nodeFile and dataDir
			os.Remove(nodeFile)
			os.Remove(dataDir)
		}
		fmt.Println(nodeId)
		file.WriteString(nodeId.String())
		file.Close()
	} else {
		//
	}

	if err != nil {
		return err
	}
	fmt.Printf("Self check done\n")
	return nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
