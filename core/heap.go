package core

type PacketHeap struct {
	storage []*SimulatedPacket
}

func (q *PacketHeap) IsEmpty() bool {
	return q.Len() == 0
}

func (q *PacketHeap) Less(i, j int) bool {
	return q.storage[i].EmitTime.Before(q.storage[j].EmitTime)
}

func (q *PacketHeap) Len() int {
	return len(q.storage)
}

func (q *PacketHeap) Swap(i, j int) {
	q.storage[i], q.storage[j] = q.storage[j], q.storage[i]
}

func (q *PacketHeap) Push(x interface{}) {
	q.storage = append(q.storage, x.(*SimulatedPacket))
}

func (q *PacketHeap) Pop() interface{} {
	x := q.storage[q.Len()-1]
	q.storage = q.storage[:q.Len()-1]
	return x
}

func (q *PacketHeap) Peek() *SimulatedPacket {
	if q.IsEmpty() {
		return nil
	}
	return q.storage[q.Len()-1]
}
