package base

// Queue is a ring queue implemented in array/slice
// O(1) time complexity for enqueue/dequeue/random access
// O(n) space complexity, where n is the maximum count of elements in the queue at the same time
type Queue struct {
	head, tail, length int
	storage            []interface{}
}

// NewQueue create a queue with the initial capacity hint
func NewQueue(length int) *Queue {
	storage := make([]interface{}, length+1)
	for len(storage) < cap(storage) {
		storage = append(storage, nil)
	}
	return &Queue{head: 0, tail: 0, length: len(storage), storage: storage}
}

func (q *Queue) expand() {
	for i := 0; i < q.tail; i++ {
		q.storage = append(q.storage, q.storage[i])
		q.storage[i] = nil
		q.length++
	}
	q.tail = q.length
	q.storage = append(q.storage, nil)
	q.length++
	for q.length < cap(q.storage) {
		q.storage = append(q.storage, nil)
		q.length++
	}
}

// IsEmpty is used to determine whether the queue is empty
func (q *Queue) IsEmpty() bool {
	return q.Length() == 0
}

// Length of the queue
func (q *Queue) Length() int {
	result := q.tail - q.head
	for result < 0 {
		result += q.length
	}
	return result
}

// Enqueue insert the given element to the end of queue
func (q *Queue) Enqueue(data interface{}) {
	q.storage[q.tail] = data
	q.tail++
	if q.tail >= q.length {
		q.tail = 0
	}
	if q.tail == q.head {
		q.expand()
	}
}

// Dequeue remove and return element at the head of the queue, panic if empty
func (q *Queue) Dequeue() interface{} {
	if q.head == q.tail {
		panic("queue is empty")
	}
	if q.head >= q.length {
		q.head = 0
	}
	result := q.storage[q.head]
	q.storage[q.head] = nil
	q.head++
	return result
}

// At return the element of the given index
func (q *Queue) At(index int) interface{} {
	if index >= q.Length() {
		panic("index overflow")
	}
	index += q.head
	for index >= q.length {
		index -= q.length
	}
	return q.storage[index]
}

// Do iterate the queue with the given action
func (q *Queue) Do(action func(interface{})) {
	for i := q.head; i != q.tail; i++ {
		if i >= q.length {
			i = -1
		}
		action(q.storage[i])
	}
}
