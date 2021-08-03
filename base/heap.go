package base

// PacketHeap ...
type PacketHeap struct {
	Storage []*SimulatedPacket
}

func (q *PacketHeap) IsEmpty() bool {
	return q.Len() == 0
}

func (q *PacketHeap) Less(i, j int) bool {
	return q.Storage[i].EmitTime.Before(q.Storage[j].EmitTime)
}

func (q *PacketHeap) Len() int {
	return len(q.Storage)
}

func (q *PacketHeap) Swap(i, j int) {
	q.Storage[i], q.Storage[j] = q.Storage[j], q.Storage[i]
}

func (q *PacketHeap) Push(x interface{}) {
	q.Storage = append(q.Storage, x.(*SimulatedPacket))
}

func (q *PacketHeap) Pop() interface{} {
	x := q.Storage[q.Len()-1]
	q.Storage = q.Storage[:q.Len()-1]
	return x
}

func (q *PacketHeap) Peek() *SimulatedPacket {
	if q.IsEmpty() {
		return nil
	}
	return q.Storage[0]
}
