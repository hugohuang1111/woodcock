package room

var (
	tableMap map[uint64]*table
)

func init() {
	tableMap = make(map[uint64]*table)
}

func getOrCreateTable(tid uint64) *table {
	t, exist := tableMap[tid]
	if exist {
		return t
	}
	t = newTable()
	tableMap[tid] = t

	return t
}
