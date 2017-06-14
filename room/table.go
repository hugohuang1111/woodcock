package room

import (
	"time"

	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

const (
	roomPhaseWaiting = 0
	roomPhaseShuffle = 1
	roomPhaseDealing = 2
	roomPhasePlaying = 3
	roomPhaseSettle  = 4
)

type table struct {
	users           [4]uint64
	phase           int
	robotEntryTimer *time.Timer
	invitRobotID    uint64
	updateChan      chan bool
}

func newTable() *table {
	t := new(table)
	t.invitRobotID = 9000
	t.updateChan = make(chan bool, 1024)
	go t.play()
	return t
}

func (t *table) sitDown(uid uint64) {
	for idx, val := range t.users {
		if 0 == val {
			t.users[idx] = uid
			t.updateChan <- true
			break
		}
	}
}

func (t *table) standUp(uid uint64) {
	for idx, val := range t.users {
		if uid == val {
			t.users[idx] = 0
			break
		}
	}
}

func (t *table) play() {
	for {
		<-t.updateChan
		switch t.phase {
		case roomPhaseWaiting:
			if nil != t.robotEntryTimer {
				t.robotEntryTimer.Stop()
			}
			t.robotEntryTimer = time.AfterFunc(5*time.Second, func() {
				msg := new(module.Message)
				msg.Sender = module.MOD_ROOM
				msg.Recver = module.MOD_ROOM
				msg.Userid = t.invitRobotID
				msg.Type = module.MOD_MSG_TYPE_ENTRY_ROOM
				router.Route(msg)
				t.invitRobotID++
				if t.invitRobotID > 9010 {
					t.invitRobotID = 9000
				}
			})
		case roomPhaseShuffle:
		case roomPhaseDealing:
		case roomPhasePlaying:
		case roomPhaseSettle:
		default:
		}
	}

}
