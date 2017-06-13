package room

const (
	roomPhaseWaiting = 0
	roomPhaseShuffle = 1
	roomPhaseDealing = 2
	roomPhasePlaying = 3
	roomPhaseSettle  = 4
)

type table struct {
	users [4]uint64
	phase int
}

func newTable() *table {
	return new(table)
}

func (t *table) sitDown(uid uint64) {
	for idx, val := range t.users {
		if 0 == val {
			t.users[idx] = uid
			break
		}
	}
}

func (t *table) standUP(uid uint64) {
	for idx, val := range t.users {
		if uid == val {
			t.users[idx] = 0
			break
		}
	}
}

func (t *table) play() {

}
