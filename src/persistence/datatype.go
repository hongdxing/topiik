package persistence

import "os"

type fetchingCache struct {
	F   *os.File
	seq int64
	pos int64
}
