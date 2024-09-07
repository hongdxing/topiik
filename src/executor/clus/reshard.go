//author: Duan HongXing
//date: 6 Sep, 2024

package clus

import "topiik/internal/datatype"

// Reshard partition
// After new-partition or remove-partition, the partitions change not take effect yet
// Only after Reshard partition executed, Topiik start to split or merge partitions
func Reshard(req datatype.Req) (result string, err error) {

	return "", nil
}
