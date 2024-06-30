/***
**
**
**
**
**/

package shared

import (
	"topiik/internal/datatype"
)

// cluster info
var Cluster = &datatype.Cluster{}

// the kv map
var MemMap = make(map[string]*datatype.TValue)
