package gate

type connect interface {
	setID(id uint64)
	ID() uint64
	send(data []byte)
	onRecv(data []byte)
	run()
	close()
}
