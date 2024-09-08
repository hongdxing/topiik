//author: Duan Hongxing
//date: 22 Jun, 2024

package util

import (
	"encoding/binary"
	"os"
)

func WriteBinaryFile(path string, data []byte) (err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	err = binary.Write(f, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	return nil
}

func ReadBinaryFile(path string) ([]byte, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0664)
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
	err = binary.Read(file, binary.LittleEndian, &bytes)
	/*bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	*/
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
