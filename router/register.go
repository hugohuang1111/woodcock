package router

var regTable map[uint64]uint64

func init() {
	regTable = make(map[uint64]uint64)
}

func regUID(uID, connID uint64) {
}
