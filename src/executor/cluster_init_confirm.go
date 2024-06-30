/***
** author: duan hongxing
** date: 30 Jun 2024
** desc:
**	Init cluster command implement
**
**/

package executor

import (
	"errors"
	"fmt"
	"topiik/cluster"
)

func clusterInitConfirm() (err error) {
	if cluster.IsNodeEmpty() {
		return nil
	}
	fmt.Println("=====================")
	return errors.New(cluster.CLUSTER_INIT_FAILED)
}
