package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/rs/zerolog/log"
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

func EncodeB(msg []byte) ([]byte, error) {
	// Lenght of message, int32(4 bytes)
	var length = int32(len(msg))
	var buffer = new(bytes.Buffer)
	// Write message HEADER
	err := binary.Write(buffer, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// Write message BODY
	err = binary.Write(buffer, binary.LittleEndian, msg)
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
	lengthByte, err := reader.Peek(4)
	if err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err = binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}

	var size int = int(length) + 4

	var ready int = 0
	buf := make([]byte, length)
	ready, err = reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if ready < size {
		for {
			tmp := make([]byte, size-ready)
			i, err := reader.Read(tmp)
			if err != nil {
				return nil, err
			}
			buf = append(buf, tmp...)
			ready += i
			if ready == size {
				break
			}
		}
	}

	/*
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
	*/
	return buf, nil
}

func EncodeHeader(icmd uint8, ver uint8) (header []byte, err error) {
	var buffer = new(bytes.Buffer)
	// Write VER
	err = binary.Write(buffer, binary.LittleEndian, ver)
	if err != nil {
		return nil, err
	}
	header = append(header, buffer.Bytes()...)
	buffer.Reset()

	// Write CMD
	err = binary.Write(buffer, binary.LittleEndian, icmd)
	if err != nil {
		return nil, err
	}
	header = append(header, buffer.Bytes()...)
	return header, nil
}

func DecodeHeader(buf []byte) (icmd uint8, ver uint8, err error) {
	if len(buf) < 2 {
		return 0, 0, errors.New("SYNTAX_ERR")
	}
	byteBuf := bytes.NewBuffer([]byte{buf[0]})
	err = binary.Read(byteBuf, binary.LittleEndian, &ver)
	if err != nil {
		log.Err(err)
		return 0, 0, errors.New("SYNTAX_ERR")
	}
	byteBuf = bytes.NewBuffer([]byte{buf[1]})
	err = binary.Read(byteBuf, binary.LittleEndian, &icmd)
	if err != nil {
		log.Err(err)
		return 0, 0, errors.New("SYNTAX_ERR")
	}
	return icmd, ver, nil
}
