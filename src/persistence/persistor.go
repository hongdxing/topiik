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
const ptnFolderPrefix = "ptn-"

// Parameters:
//   - binlogs: one or more logs
func ReceiveBinlog(data []byte) (seq int64, err error) {
	ptnId := string(data[:consts.PTN_ID_LEN])
	binlogs := data[consts.PTN_ID_LEN:]
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

	// make sure partition folder exists
	path := ptnFolder(ptnId)
	exists, err := util.PathExists(path)
	if err != nil {
		l.Err(err).Msgf("persistence::getPtnActiveBLF %s", err.Error())
		return nil, err
	}
	if !exists {
		err = os.Mkdir(path, 0664)
		if err != nil {
			return nil, err
		}
	}

	// if binlog seq not set in memory yet
	// 2 conditions:
	//  1) brand new partition
	//	2) the server restarted, lost memory seq, need read last log from binlog
	if _, ok := blseq[ptnId]; !ok {
		fname, err := maxPtnBLF(ptnId)
		if err != nil {
			l.Err(err).Msgf("persistence::getPtnActiveBLF %s", err.Error())
			return nil, err
		}
		// if binlog sequence not set yet(brand new partition)
		if fname == "" {
			blseq[ptnId] = 0
		} else {
			// parse file to get biggest
			fname = ptnFolder(ptnId) + fname
			f, err := os.OpenFile(fname, os.O_RDONLY, 0664)
			if err != nil {
				l.Err(err).Msgf("persistence::getPtnActiveBLF %s", err.Error())
				return nil, err
			}

			// parse the binlog one by one to get last
			// todo: use index
			var seq int64
			for {
				_, err := parseOne(f, &seq)
				if err != nil {
					if err != io.EOF {
						l.Panic().Msg(err.Error())
					}
					break
				}
			}
			blseq[ptnId] = seq
		}

	}
	// sequence pre padding 0s, totla lenght 20
	// eg: 00000000000000000999.bin
	fname := fmt.Sprintf("%020d.bin", blseq[ptnId])
	path = path + consts.SLASH + fname
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		l.Err(err).Msgf("persistence::getPtnActiveBLF create binlog file failed: %s", err.Error())
		return nil, err
	}
	//cache file
	blfs[ptnId] = f
	return f, nil
}

// get the latest binlog file by walk the partition folder
func maxPtnBLF(ptnId string) (fname string, err error) {
	dir, err := os.Open(ptnFolder(ptnId))
	if err != nil {
		l.Err(err).Msgf("persistence::maxPtnBLF %s", err.Error())
		return fname, err
	}
	fs, err := dir.Readdir(0)
	if err != nil {
		l.Err(err).Msgf("persistence::maxPtnBLF %s", err.Error())
		return fname, err
	}
	for _, f := range fs {
		if !f.IsDir() {
			if f.Name() > fname {
				fname = f.Name()
			}
		}
	}
	return fname, nil
}

// the partition folder full path with slash endding
func ptnFolder(ptnId string) string {
	return util.GetDataFullPath() + ptnFolderPrefix + ptnId + consts.SLASH
}
