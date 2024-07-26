/*
* author: duan hongxing
* date: 22 Jun 2024
* desc:
 */

package util

import (
	"encoding/binary"
	"os"
)

func WriteBinaryFile(path string, data []byte) (err error) {
	var file *os.File
	exist, err := PathExists(path)
	if err != nil {
		return err
	}
	if !exist {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	err = binary.Write(file, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func ReadBinaryFile(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, err
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)
	binary.Read(file, binary.LittleEndian, &bytes)

	/*bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return nil, err
	}*/

	return data, nil
}
