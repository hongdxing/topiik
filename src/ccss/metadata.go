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
var salorMap = make(map[string]Salor)
var partitionMap = make(map[string]Partition)

func InitMetadata(capitals map[string]Capital, solars map[string]Salor, partitions map[string]Partition) {
	captialMap = capitals
	fmt.Println("capitals")
	fmt.Println(captialMap)
	salorMap = solars
	fmt.Println("salors")
	fmt.Println(salorMap)
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
