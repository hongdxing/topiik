package main

func vote(cTerm int) string {
	// (lastTermV > lastTermC) || ((lastTermV == lastTermC) && (lastIndexV > lastIndexC))
	if NodeStatus.Term > uint(cTerm) {
		return RES_REJECTED
	} else {
		return RES_ACCEPTED
	}
}
