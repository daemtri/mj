package data

import (
	"config"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"sync"
)

var gen *ServerIDGen
var lock sync.Mutex

var lockroomid sync.Mutex

var roomIDGen *RoomIDGen

func InitIDGen() {
	gen = &ServerIDGen{ServerID: strconv.Itoa(config.Opts().Server_id)}
	gen.Insert()

	roomIDGen = &RoomIDGen{ServerID: uint64(config.Opts().Server_id)}
	roomIDGen.Insert()
}

type RoomIDGen struct {
	ServerID   uint64 `bson:"_id"`
	LastRoomID uint64 `bson:"LastRoomID"`
}

func (this *RoomIDGen) Insert() error {
	count, _ := C(_GEN_ROOM_ID).Find(bson.M{"_id": this.ServerID}).Count()
	if count == 0 {
		this.LastRoomID = 1000
		return C(_GEN_ROOM_ID).Insert(this)
	}
	return nil
}

func (this *RoomIDGen) Get() (uint64, error) {
	lockroomid.Lock()
	defer lockroomid.Unlock()
	this.ServerID = uint64(config.Opts().Server_id)
	err := C(_GEN_ROOM_ID).UpdateId(this.ServerID, bson.M{"$inc": bson.M{"LastRoomID": 1}})
	if err != nil {
		return 0, err
	}
	err = C(_GEN_ROOM_ID).FindId(this.ServerID).One(this)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strconv.FormatUint(this.ServerID, 10)+strconv.FormatUint(this.LastRoomID, 10), 10, 64)
}

func GenRoomID() (uint64, error) {
	return roomIDGen.Get()
}

type ServerIDGen struct {
	ServerID   string `bson:"_id"`
	LastUserID uint64 `bson:"LastUserID"`
}

func GenUserID() (string, error) {
	return gen.Get()
}

func (s *ServerIDGen) Exists() bool {
	count, _ := C(_GEN_USER_ID).Find(bson.M{"_id": s.ServerID}).Count()
	return count != 0
}

func (s *ServerIDGen) Insert() error {
	count, _ := C(_GEN_USER_ID).Find(bson.M{"_id": s.ServerID}).Count()
	if count == 0 {
		s.LastUserID = 6000
		return C(_GEN_USER_ID).Insert(s)
	}
	return nil
}

func (s *ServerIDGen) Get() (string, error) {
	lock.Lock()
	defer lock.Unlock()
	s.ServerID = strconv.Itoa(config.Opts().Server_id)
	err := C(_GEN_USER_ID).UpdateId(s.ServerID, bson.M{"$inc": bson.M{"LastUserID": 1}})
	if err != nil {
		return "", err
	}
	err = C(_GEN_USER_ID).FindId(s.ServerID).One(s)
	if err != nil {
		return "", err
	}

	return s.ServerID + strconv.FormatUint(s.LastUserID, 10), err
}
