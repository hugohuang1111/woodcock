package room

import (
	"crypto/sha512"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math/rand"
	"sort"
	"time"

	"github.com/golang/glog"
	"github.com/hugohuang1111/woodcock/module"
	"github.com/hugohuang1111/woodcock/router"
)

const (
	roomPhaseWaiting      = 0
	roomPhaseShuffle      = 1
	roomPhaseDealing      = 2
	roomPhaseMakeAAbandon = 3
	roomPhasePlaying      = 4
	roomPhaseSettle       = 5
)

const (
	suitDot       = 1
	suitBamboo    = 2
	suitCharacter = 3
)

const (
	meldPong = "pong"
	meldKong = "kong"
	meldChow = "chow"
)

type tile struct {
	Suit int `json:"suit"` //dot, bamboo, character
	Rank int `json:"rank"`
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
	AbandonTile  [4]int       `json:"abandonTile"`  //缺的牌
}

type room struct {
	id              uint64
	users           [4]uint64
	phase           int
	robotEntryTimer *time.Timer
	invitRobotID    uint64
	updateChan      chan bool
	userTiles       tiles
	banker          int
	current         int
	abandonSuits    [4]int
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

func (r *room) abandonSuit(userID uint64, suit int) error {
	if roomPhaseMakeAAbandon != r.phase {
		return errors.New("phase not right")
	}

	for i, id := range r.users {
		if id == userID {
			r.abandonSuits[i] = suit
		}
	}

	r.goToPlaying()

	return nil
}

func (r *room) forceAbandonSuit() {
	if roomPhaseMakeAAbandon != r.phase {
		return
	}
	for i, suit := range r.abandonSuits {
		if 0 == suit {
			var dot int
			var bamboo int
			var character int
			for _, tile := range r.userTiles.HoldTiles[i] {
				switch tile.Suit {
				case suitDot:
					dot++
				case suitBamboo:
					bamboo++
				case suitCharacter:
					character++
				}
			}
			if dot < bamboo && dot < character {
				r.abandonSuits[i] = suitDot
			} else if bamboo < character {
				r.abandonSuits[i] = suitBamboo
			} else {
				r.abandonSuits[i] = suitCharacter
			}
		}
	}
}

func (r *room) goToPlaying() {
	if roomPhaseMakeAAbandon != r.phase {
		return
	}
	var ok = true
	for _, suit := range r.abandonSuits {
		if 0 == suit {
			ok = false
		}
	}
	if ok {
		r.sortHoldTiles()
		r.phase = roomPhasePlaying
		r.updateChan <- true
	}
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
				r.phase = roomPhaseDealing
				r.updateChan <- true
			}
		case roomPhaseShuffle:
		case roomPhaseDealing:
			r.shuffleAndDealingTiles()
			time.AfterFunc(1*time.Second, func() {
				r.phase = roomPhaseMakeAAbandon
				r.updateChan <- true
			})
		case roomPhaseMakeAAbandon:
			time.AfterFunc(1*time.Second, func() {
				r.autoAbandonSuit()
				r.goToPlaying()
			})
			time.AfterFunc(10*time.Second, func() {
				r.forceAbandonSuit()
				r.goToPlaying()
			})
		case roomPhasePlaying:
		case roomPhaseSettle:
		default:
		}

		r.broadcastScene()
	}
}

func (r *room) shuffleAndDealingTiles() {
	rawCards := rand.Perm(108)
	walls := make([]map[string]int, 0, 4)
	t := map[string]int{
		"start":  0,
		"length": 26,
	}
	walls = append(walls, t)
	t = map[string]int{
		"start":  26,
		"length": 28,
	}
	walls = append(walls, t)
	t = map[string]int{
		"start":  54,
		"length": 26,
	}
	walls = append(walls, t)
	t = map[string]int{
		"start":  80,
		"length": 28,
	}
	walls = append(walls, t)

	int64ToBytes := func(i int64) []byte {
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(i))
		return buf
	}
	bytesToInt64 := func(buf []byte) int64 {
		return int64(binary.BigEndian.Uint64(buf))
	}
	sumBytes := sha512.Sum512(int64ToBytes(time.Now().UnixNano()))
	rand.Seed(bytesToInt64(sumBytes[1:9]))
	die1 := rand.Intn(6) + 1
	rand.Seed(bytesToInt64(sumBytes[10:18]))
	die2 := rand.Intn(6) + 1
	grabStartChair := (die1 + die2 + r.banker - 1) % 4
	grabStartPos := die1
	if grabStartPos > die2 {
		grabStartPos = die2
	}
	grabStartPos *= 2

	start, _ := walls[grabStartChair]["start"]
	end := start + walls[grabStartChair]["length"]
	end--
	start = end - grabStartPos

	grabTiles := func(start, len int) (arr []int, pos int) {
		arr = make([]int, 0, len)
		p := start
		for i := 0; i < len; i++ {
			p += 108
			p %= 108
			arr = append(arr, rawCards[p])
			rawCards[p] = -1
			p--
		}

		pos = p
		return
	}
	i2t := func(i int) tile {
		var s int
		switch {
		case i < 37:
			s = suitDot
		case i < 73:
			s = suitBamboo
		case i < 109:
			s = suitCharacter
		default:
			glog.Error("i2t, shoule be here wrong i", i)
		}
		i %= 9
		i++
		return tile{Rank: i, Suit: s}
	}
	var grabArr []int
	for i := 0; i < 4; i++ {
		tiles := make([]tile, 0, 14)
		if i == r.banker {
			grabArr, start = grabTiles(start, 14)
		} else {
			grabArr, start = grabTiles(start, 13)
		}
		for _, v := range grabArr {
			tiles = append(tiles, i2t(v))
		}
		r.userTiles.HoldTiles[i] = tiles
	}

	for i, wall := range walls {
		start, _ := wall["start"]
		length, _ := wall["length"]
		wallTiles := make([]tile, 0, 28)
		for i := 0; i < length; i++ {
			v := rawCards[start+i]
			if v >= 0 {
				wallTiles = append(wallTiles, i2t(v))
			}
		}
		r.userTiles.WallTiles[i] = wallTiles
	}
	r.sortHoldTiles()
}

func (r *room) sortHoldTiles() {
	resaved := r.current
	for i := range r.users {
		r.current = i
		sort.Sort(r)
	}
	r.current = resaved
}

func (r *room) broadcastScene() {
	scene := make(map[string]interface{})
	scene["version"] = 1
	scene["type"] = "room:scene"
	scene["phase"] = r.phase
	scene["users"] = r.users
	scene["banker"] = r.banker //chair id
	scene["curUser"] = r.current
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

func (r *room) Len() int {
	return len(r.userTiles.HoldTiles[r.current])
}

func (r *room) Less(i, j int) bool {
	tiles := r.userTiles.HoldTiles[r.current]
	if tiles[i].Suit == tiles[j].Suit {
		return tiles[i].Rank < tiles[j].Rank
	} else if tiles[i].Suit == r.abandonSuits[r.current] {
		return false
	} else if tiles[j].Suit == r.abandonSuits[r.current] {
		return true
	}
	return tiles[i].Suit < tiles[j].Suit

}

func (r *room) Swap(i, j int) {
	tiles := r.userTiles.HoldTiles[r.current]
	tiles[i], tiles[j] = tiles[j], tiles[i]
}
