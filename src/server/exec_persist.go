//author: Duan Hongxing
//date: 15 Sep, 2024

package server

import (
	"topiik/persistence"
	"topiik/resp"
)

func persist(msg []byte) (rslt string, err error) {
	err = persistence.Append(msg)
	if err != nil {
		return resp.RES_NIL, err
	}
	return resp.RES_OK, nil
}
