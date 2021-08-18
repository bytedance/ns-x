package base

type Packet interface {
	Size() int
}

type RawPacket []byte

func (p RawPacket) Size() int {
	return len(p)
}
