package module

const (
	MOD_ROUTER = "router"
	MOD_GATE   = "gate"
	MOD_USER   = "user"
	MOD_ROOM   = "room"
)

const (
	MsgTypeUnknow     = "unknow"
	MsgTypeClient     = "client"
	MsgTypeDisconnect = "disconnect"
	MsgTypeGetUserID  = "getUserID"
	MsgTypeEntryRoom  = "entryRoom"
)

const (
	PayloadKeyClientData = "clientData"
	PayloadKeyConnectID  = "connectID"
	PayloadKeyUserID     = "userID"
)
