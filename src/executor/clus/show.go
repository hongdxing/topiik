/*
* ©2024 www.topiik.com
* Author: Duan HongXing
* Date: 30 Aug, 2024
 */

package clus

import (
	"topiik/cluster"
	"topiik/internal/datatype"
)

func Show(req datatype.Req) (rslt string) {

	return cluster.Show()
}
