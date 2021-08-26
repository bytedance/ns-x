package base

// EventHeap is used to sort events according to the time of events
type EventHeap struct {
	storage []Event
}

func NewEventHeap(events ...Event) *EventHeap {
	return &EventHeap{storage: events}
}

func (q *EventHeap) IsEmpty() bool {
	return q.Len() == 0
}

func (q *EventHeap) Less(i, j int) bool {
	ti := q.storage[i].Time()
	tj := q.storage[j].Time()
	return ti.Before(tj)
}

func (q *EventHeap) Len() int {
	return len(q.storage)
}

func (q *EventHeap) Swap(i, j int) {
	q.storage[i], q.storage[j] = q.storage[j], q.storage[i]
}

func (q *EventHeap) Push(x interface{}) {
	q.storage = append(q.storage, x.(Event))
}

func (q *EventHeap) Pop() interface{} {
	x := q.storage[q.Len()-1]
	q.storage = q.storage[:q.Len()-1]
	return x
}

func (q *EventHeap) Peek() Event {
	if q.IsEmpty() {
		return nil
	}
	return q.storage[0]
}
