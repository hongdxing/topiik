package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

// Encode
func Encode(message string) ([]byte, error) {
	// Lenght of message, int32(4 bytes)
	var length = int32(len(message))
	var buffer = new(bytes.Buffer)
	// Write message HEADER
	err := binary.Write(buffer, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// Write message BODY
	err = binary.Write(buffer, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Decode
/*
func Decode(reader *bufio.Reader) (string, error) {
	// Read message HEADER(int32 4 bytes)
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}

	// Readable data in Buffer
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// Read message
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}*/

func Decode(reader *bufio.Reader) ([]byte, error) {
	// Read message HEADER(int32 4 bytes)
	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	// Readable data in Buffer
	if int32(reader.Buffered()) < length+4 {
		return nil, err
	}

	// Read message
	buf := make([]byte, int(4+length))
	_, err = reader.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
