package module

// SkelectionCB skelection interface
type SkelectionCB interface {
	OnMsg(msg *Message)
}

// Skelecton module skelecton
type Skelecton struct {
	msgList chan *Message
	running bool
	cb      SkelectionCB
}

// Run run
func (s *Skelecton) Run(cb SkelectionCB) {
	s.running = true
	s.cb = cb
	s.msgList = make(chan *Message, 1024)
	go s.work()
}

// Stop stop
func (s *Skelecton) Stop() {
	s.running = false
}

// Add add module msg
func (s *Skelecton) Add(msg *Message) {
	s.msgList <- msg
}

func (s *Skelecton) work() {
	for s.running {
		msg := <-s.msgList
		s.cb.OnMsg(msg)
	}
}
