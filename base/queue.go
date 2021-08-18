package base

// DataQueue is a ring queue of fixed size, it's usually used as a history record of data
type DataQueue struct {
	head, tail, length int
	storage            []interface{}
}

func NewDataQueue(length int) *DataQueue {
	return &DataQueue{head: 0, tail: 0, length: length + 1, storage: make([]interface{}, length+1)}
}

func (q *DataQueue) IsEmpty() bool {
	return q.Length() == 0
}

func (q *DataQueue) Length() int {
	result := q.tail - q.head
	for result < 0 {
		result += q.length
	}
	return result
}

func (q *DataQueue) Enqueue(data interface{}) {
	q.storage[q.tail] = data
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

func (q *DataQueue) Dequeue() interface{} {
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

func (q *DataQueue) At(index int) interface{} {
	if index >= q.Length() {
		panic("index is overflow")
	}
	index += q.head
	for index >= q.length {
		index -= q.length
	}
	return q.storage[index]
}

// Do iterate the queue with the given action
func (q *DataQueue) Do(action func(interface{})) {
	for i := q.head; i != q.tail; i++ {
		if i >= q.length {
			i = -1
		}
		action(q.storage[i])
	}
}
