/***
** author: duan hongxing
** data: 3 Jul 2024
** desc:
**
**/

package ccss

type Capital struct {
	Id      string
	Address string
}

type Sailor struct {
	Id      string
	Address string
}

type Partition struct {
	Id            string   // Id of the partition, random 16 alphnum
	LeaderSailorId string   // The Sailor Id where the Leader Partition  located
	SailorIds      []string // The Sailors where the Partition located
}
