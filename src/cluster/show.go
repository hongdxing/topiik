//©2024 www.topiik.com
//Author: Duan HongXing
//Date: 30 Aug, 2024

package cluster

import (
	"encoding/json"
	"topiik/internal/datatype"
)

// Show cluster information
func Show() string {
	cluData := datatype.ClusterData{}

	for _, ptn := range partitionInfo.PtnMap {
		ptnData := datatype.PartitionData{
			Id: ptn.Id,
		}
		ptnData.SlotFrom = ptn.SlotFrom
		ptnData.SlotTo = ptn.SlotTo
		for _, nd := range ptn.NodeSet {
			ptnData.Nodes = append(ptnData.Nodes, datatype.NodeData{Id: nd.Id, Address: nd.Addr})
		}
		cluData.Partitions = append(cluData.Partitions, ptnData)
	}

	rslt, err := json.MarshalIndent(cluData, "", "    ")
	if err != nil {
		return err.Error()
	}
	return string(rslt)
}
