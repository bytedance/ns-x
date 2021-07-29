package byte_ns

type PacketQueue struct {
	head, tail, length int
	storage            []*SimulatedPacket
}

func NewPacketQueue(length int) *PacketQueue {
	return &PacketQueue{head: 0, tail: 0, length: length + 1, storage: make([]*SimulatedPacket, length+1, length+1)}
}

func (q *PacketQueue) IsEmpty() bool {
	return q.Length() == 0
}

func (q *PacketQueue) Length() int {
	result := q.tail - q.head
	for result < 0 {
		result += q.length
	}
	return result
}

func (q *PacketQueue) Enqueue(packet *SimulatedPacket) {
	q.storage[q.tail] = packet
	q.tail++
	if q.tail >= q.length {
		q.tail = 0
	}
	if q.head == q.tail {
		q.head++
	}
	if q.head >= q.length {
		q.head = 0
	}
}

func (q *PacketQueue) Dequeue() *SimulatedPacket {
	if q.head == q.tail {
		panic("record is empty")
	}
	result := q.storage[q.head]
	q.head++
	if q.head >= q.length {
		q.head = 0
	}
	return result
}

func (q *PacketQueue) At(index int) *SimulatedPacket {
	if index >= q.Length() {
		panic("index is overflow")
	}
	index += q.head
	for index >= q.length {
		index -= q.length
	}
	return q.storage[index]
}

func (q *PacketQueue) Do(action func(simulatedPacket *SimulatedPacket)) {
	for i := q.head; i != q.tail; i++ {
		if i >= q.length {
			i = -1
		}
		action(q.storage[i])
	}
}
