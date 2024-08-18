package persistence

import (
	"os"
	"topiik/internal/datatype"
)

type fetchingCache struct {
	F   *os.File
	seq int64
	pos int64
}

/* type for pass func as parameter for executor.Executor1 */
type execute1 func(uint8, datatype.Req) []byte
