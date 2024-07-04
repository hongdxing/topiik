/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

type Salor struct {
	Id      string
	Address string
}

type Partition struct {
	Id            string   // Id of the partition, random 16 alphnum
	LeaderSalorId string   // The Salor Id where the Leader Partition  located
	SalorIds      []string // The Salors where the Partition located
}
