//author: Duan Hongxing
//date: 13 Aug, 2024

package persistence

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"topiik/internal/consts"
	"topiik/internal/util"
)

// binlog files, the keys are partition id
var blfs = make(map[string]*os.File)

// binlog sequences, the keys are partition id
var blseq = make(map[string]int64)

const binlogFileMaxSize = 1073741824 // 1G

// Parameters:
//   - binlogs: one or more logs
func ReceiveBinlog(data []byte) (seq int64, err error) {
	ptnId := string(data[:consts.PTN_ID_LEN])
	binlogs := data[consts.PTN_ID_LEN:]
	fmt.Println(ptnId)
	/*
		if flrActiveBLF == nil {
			path := getPtnActiveBLF(ptnId)
			flrActiveBLF, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0664)
			if err != nil {
				l.Err(err).Msgf("persistence::ReceiveBinlog %s", err.Error())
			}
		}*/
	f, err := getPtnActiveBLF(ptnId)
	if err != nil {
		return seq, err
	}

	//var seq int64
	reader := bytes.NewReader(binlogs)
	for {
		//buf, err := parseOne(reader, &seq)
		// binlog
		bl, err := receiveOne(reader)
		if err != nil {
			if err != io.EOF {
				l.Err(err).Msg(err.Error())
			}
			break
		}
		// prepend seq
		blseq[ptnId] += 1
		bbuf := new(bytes.Buffer)
		binary.Write(bbuf, binary.LittleEndian, blseq[ptnId])
		buf := bbuf.Bytes()
		buf = append(buf, bl...)

		err = binary.Write(f, binary.LittleEndian, buf)
		if err != nil {
			// ??? need deduct???
			//binlogSeq-- // parse ok, but write failed, return pre seq
			break
		}
	}
	return blseq[ptnId], nil
}

func receiveOne(r io.Reader) (res []byte, err error) {
	var length int32

	// 1: read length
	buf := make([]byte, 4)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	bbuf := bytes.NewBuffer(buf)
	binary.Read(bbuf, binary.LittleEndian, &length)
	res = append(res, buf...)

	// 2: read msg
	buf = make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		if err != io.EOF {
			l.Err(err).Msg(err.Error())
		}
		return nil, err
	}
	res = append(res, buf...)
	return res, nil
}

var flocker sync.Mutex

// get partition active binlog file
func getPtnActiveBLF(ptnId string) (*os.File, error) {
	flocker.Lock()
	defer flocker.Unlock()
	if f, ok := blfs[ptnId]; ok {
		fi, err := f.Stat()
		if err != nil {
			l.Warn().Msgf("persistence::getPtnActiveBLF f.Stat() failed")
			if f != nil {
				f.Close()
			}
		} else {
			if fi.Size() > binlogFileMaxSize {
				f.Close()
			} else {
				return f, nil
			}
		}
	}
	// if binlog sequence not set yet(brand new partition)
	if _, ok := blseq[ptnId]; !ok {
		blseq[ptnId] = 0
	}
	// sequence pre padding 0s, totla lenght 20
	// eg: ptnid-00000000000000000999.bin
	fname := fmt.Sprintf("%s-%020d.bin", ptnId, blseq[ptnId])
	path := util.GetMainPath() + consts.SLASH + consts.DATA_DIR + consts.SLASH + fname
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		l.Err(err).Msgf("persistence::getPtnActiveBLF create binlog file failed: %s", err.Error())
		return nil, err
	}
	//cache file
	blfs[ptnId] = f
	return f, nil
}
