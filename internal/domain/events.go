package domain

const (
	StatusAlive = iota
	StatusSuspect
	StatusDead
)

type Event struct {
	ID        string
	Topic     string
	Payload   []byte
	Timestamp int64
}

type Node struct {
	ID          string
	Address     string
	Incarnation uint32
	Status      uint8
}
