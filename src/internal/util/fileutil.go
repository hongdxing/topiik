/*
* author: duan hongxing
* date: 22 Jun 2024
* desc:
 */

package util

import (
	"bufio"
	"encoding/json"
	"os"
)

func WriteBinaryFile(path string, data []byte) error {

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

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	//binary.Read(file, binary.LittleEndian, &bytes)

	return data, nil
}
