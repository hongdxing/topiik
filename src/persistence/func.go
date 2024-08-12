package persistence

func SetPtnFlrAddr2Lst(addr2Lst []string) {
	//clear(ptnFlrAddr2Lst)
	ptnFlrAddr2Lst = nil
	ptnFlrAddr2Lst = append(ptnFlrAddr2Lst, addr2Lst...)
}
