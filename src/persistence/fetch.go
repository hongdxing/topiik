/*
* author: duan hongxing
* date: 4 Aug, 2024
* desc:
* 	Fetch binary log for respose to Sync request from follower
 */

package persistence

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"topiik/internal/util"
)

/*
 * TODO: use index to enhance
 *
 */

var cache = make(map[string]fetchingCache)

/*
* Parameters:
*	- followerId
*	- lastSeq: the last seq of follower
*
 */
func Fetch(followerId string, lastSeq int64) (res []byte) {
	l.Info().Msgf("persistence::Fetch %s %v", followerId, lastSeq)
	var file *os.File
	var seq int64 = 0
	var pos int64 = 0
	if c, ok := cache[followerId]; ok {
		file = c.F
		seq = c.seq
		pos = c.pos
	} else {
		filePath := getCurrentLogFile()
		exist, err := util.PathExists(filePath)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return
		}
		if !exist {
			return
		}
		file, err = os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			l.Err(err).Msg(err.Error())
			return
		}
		cacheVal := fetchingCache{
			F:   file,
			seq: 0,
			pos: 0,
		}
		cache[followerId] = cacheVal
	}
	//defer file.Close()

	// locate start position
	if seq == lastSeq {
		file.Seek(pos, io.SeekStart)
	} else {
		seq = 0
		pos = 0
		file.Seek(0, io.SeekStart) //set to start
		for {
			if seq < lastSeq {
				item, err := travelLog(file, &seq)
				if err != nil && err != io.EOF {
					l.Err(err).Msg(err.Error())
					return
				}

				pos += int64(len(item))
			} else {
				break
			}
		}
	}

	return res
}

/*
* Get a single log item
*
 */
func travelLog(file *os.File, seq *int64) (res []byte, err error) {
	//var seq int64
	var length int32

	// 1: read seq
	buf := make([]byte, 8)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}

	byteBuf := bytes.NewBuffer(buf)
	err = binary.Read(byteBuf, binary.LittleEndian, seq)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return
	}
	res = append(res, buf...)

	// 2: read length
	buf = make([]byte, 4)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	byteBuf = bytes.NewBuffer(buf)
	binary.Read(byteBuf, binary.LittleEndian, &length)
	res = append(res, buf...)

	// 3: read msg
	buf = make([]byte, length)
	_, err = io.ReadFull(file, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	res = append(res, buf...)
	return res, nil
}
