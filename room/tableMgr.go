package room

var (
	tableMap map[uint64]*room
)

func init() {
	tableMap = make(map[uint64]*room)
}

func getOrCreateTable(rid uint64) *room {
	t, exist := tableMap[rid]
	if exist {
		return t
	}
	t = newRoom(rid)
	tableMap[rid] = t

	return t
}
