//author: Duan Hongxing
//date: 24 Aug, 2024

package node


const (
	PTN_STS_ACTIVE  = "ACTIVE"  // Ready to use
	PTN_STS_NEW     = "NEW"     // New partition, but without Workers and Slots yet, pending Reshard
	PTN_STS_REMOVED = "REMOVED" // Mark as Removed, but still in use, pending Reshard
)
