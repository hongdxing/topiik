//author: Duan HongXing
//date: 2 Aug, 2024

package cluster

import (
	"topiik/internal/logger"
)

var l = logger.Get()

var ctlUpdCh chan struct{}
var ptnUpdCh chan struct{}

func init() {
	l.Info().Msg("Init cluster package")
	ptnUpdCh = make(chan struct{}, 2)
	ctlUpdCh = make(chan struct{}, 2)
}
