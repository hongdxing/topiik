/***
** author: duan hongxing
** data: 6 Jul 2024
** desc:
**
**/

package ccss

import "fmt"

var nodeStatus = &NodeStatus{Role: CCSS_ROLE_CO, Term: 0}
var captialMap = make(map[string]Capital)
var sailorMap = make(map[string]Sailor)
var partitionMap = make(map[string]Partition)

func InitMetadata(capitals map[string]Capital, solars map[string]Sailor, partitions map[string]Partition) {
	captialMap = capitals
	fmt.Println("capitals")
	fmt.Println(captialMap)
	sailorMap = solars
	fmt.Println("sailors")
	fmt.Println(sailorMap)
	partitionMap = partitions
	fmt.Println("partitions")
	fmt.Println(partitionMap)
}

func Map2Array[T any](theMap map[string]T) (arr []T) {
	for _, v := range theMap {
		arr = append(arr, v)
	}
	return arr
}
