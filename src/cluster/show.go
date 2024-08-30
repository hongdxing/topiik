/*
* ©2024 www.topiik.com
* Author: Duan HongXing
* Date: 30 Aug, 2024
 */

package cluster

import (
	"fmt"
	"strings"
)

func Show() string {
	var sb strings.Builder

	sb.WriteString("Partitions:\n")
	var i = 1
	for _, ptn := range partitionInfo.PtnMap {
		sb.WriteString(fmt.Sprintf("  %s\n", ptn.Id))
		sb.WriteString("    slots:")
		for _, slot := range ptn.Slots {
			sb.WriteString(fmt.Sprintf("[%v, %v]", slot.From, slot.To))
		}
		sb.WriteString("\n")
		sb.WriteString("    nodes:\n")
		var j = 1
		for _, nd := range ptn.NodeSet {
			sb.WriteString(fmt.Sprintf("      %s %s\n", nd.Id, nd.Addr))
			j++
		}
		i++
	}
	return sb.String()
}
