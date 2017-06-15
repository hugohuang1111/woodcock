package room

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"
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

const (
	suitDot       = "dot"
	suitBamboo    = "bamboo"
	suitCharacter = "character"
)

const (
	meldPong = "pong"
	meldKong = "kong"
	meldChow = "chow"
)

type tile struct {
	Suit string `json:"suit"` //dot, bamboo, character
	Rank int    `json:"rank"`
}

type meldTile struct {
	Suit  string `json:"suit"` //pong, kong
	Tiles []tile `json:"tiles"`
}

type tiles struct {
	DiscardTiles [][]tile     `json:"discardTiles"` //已出的牌
	HoldTiles    [][]tile     `json:"holdTiles"`    //持有的牌,手上的牌
	MeldTiles    [][]meldTile `json:"meldTiles"`    //碰杠的牌
	WallTiles    [][]tile     `json:"wallTiles"`    //牌墙
}

type room struct {
	id              uint64
	users           [4]uint64
	phase           int
	robotEntryTimer *time.Timer
	invitRobotID    uint64
	updateChan      chan bool
	userTiles       tiles
}

func newRoom(id uint64) *room {
	r := new(room)
	r.id = id
	r.invitRobotID = 9000
	r.updateChan = make(chan bool, 10)

	r.userTiles.DiscardTiles = make([][]tile, 4)
	r.userTiles.HoldTiles = make([][]tile, 4)
	r.userTiles.WallTiles = make([][]tile, 4)
	r.userTiles.MeldTiles = make([][]meldTile, 4)

	for i := 0; i < 4; i++ {
		r.userTiles.DiscardTiles[i] = []tile{}
		r.userTiles.HoldTiles[i] = []tile{}
		r.userTiles.WallTiles[i] = []tile{}
		r.userTiles.MeldTiles[i] = []meldTile{}
	}

	go r.play()
	return r
}

func (r *room) sitDown(uid uint64) bool {
	for idx, val := range r.users {
		if 0 == val || uid == val {
			r.users[idx] = uid
			r.updateChan <- true
			return true
		}
	}

	return false
}

func (r *room) standUp(uid uint64) {
	for idx, val := range r.users {
		if uid == val {
			r.users[idx] = 0
			break
		}
	}
}

func (r *room) userCount() int {
	var i int
	for _, uid := range r.users {
		if 0 != uid {
			i++
		}
	}

	return i
}

func (r *room) play() {
	for {
		<-r.updateChan
		switch r.phase {
		case roomPhaseWaiting:
			if nil != r.robotEntryTimer {
				r.robotEntryTimer.Stop()
			}
			if 4 != r.userCount() {
				r.robotEntryTimer = time.AfterFunc(5*time.Second, func() {
					msg := new(module.Message)
					msg.Sender = module.MOD_ROOM
					msg.Recver = module.MOD_ROOM
					msg.Type = module.MsgTypeEntryRoom
					msg.Payload = make(map[string]interface{})
					msg.Payload[module.PayloadKeyUserID] = r.invitRobotID
					msg.Payload[module.PayloadKeyRoomID] = r.id
					router.Route(msg)
					r.invitRobotID++
					if r.invitRobotID > 9010 {
						r.invitRobotID = 9000
					}
				})
			} else {
				r.phase = roomPhaseShuffle
				r.updateChan <- true
			}
		case roomPhaseShuffle:
			r.shuffleTiles()
			time.AfterFunc(5*time.Second, func() {
				r.phase = roomPhaseDealing
				r.updateChan <- true
			})
		case roomPhaseDealing:
			r.dealingTiles()
			time.AfterFunc(5*time.Second, func() {
				r.phase = roomPhasePlaying
				r.updateChan <- true
			})
		case roomPhasePlaying:
		case roomPhaseSettle:
		default:
		}

		r.broadcastScene()
	}
}

func (r *room) shuffleTiles() {

}

func (r *room) dealingTiles() {

}

func (r *room) broadcastScene() {
	scene := make(map[string]interface{})
	scene["version"] = 1
	scene["type"] = "room:scene"
	scene["phase"] = r.phase
	scene["users"] = r.users
	scene["banker"] = -1 //chair id
	scene["tiles"] = r.userTiles

	sceneString, err := json.Marshal(scene)
	if nil != err {
		glog.Error("marshal user tiles to string failed:", err)
		return
	}

	msg := new(module.Message)
	msg.Sender = module.MOD_ROOM
	msg.Recver = module.MOD_GATE
	msg.Type = module.MsgTypeClient
	msg.Payload = make(map[string]interface{})
	msg.Payload[module.PayloadKeyClientData] = sceneString

	for _, uid := range r.users {
		//if uid is 0, should be not exist, ok is false
		if info, ok := userConnMap[uid]; ok && 0 != info.connID {
			msg.Payload[module.PayloadKeyConnectID] = info.connID
			router.Route(msg)
		}
	}
}
