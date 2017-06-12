package module

//Message message pass by modules
type Message struct {
	Sender    string
	Recver    string
	Payload   []byte
	Userid    uint64
	Type      string //client
	ConnectID uint64 //connect id
}
