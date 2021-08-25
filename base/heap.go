package base

// EventHeap is used to sort events according to the time of events
type EventHeap struct {
	Storage []Event
}

func (q *EventHeap) IsEmpty() bool {
	return q.Len() == 0
}

func (q *EventHeap) Less(i, j int) bool {
	ti := q.Storage[i].Time()
	tj := q.Storage[j].Time()
	return ti.Before(tj)
}

func (q *EventHeap) Len() int {
	return len(q.Storage)
}

func (q *EventHeap) Swap(i, j int) {
	q.Storage[i], q.Storage[j] = q.Storage[j], q.Storage[i]
}

func (q *EventHeap) Push(x interface{}) {
	q.Storage = append(q.Storage, x.(Event))
}

func (q *EventHeap) Pop() interface{} {
	x := q.Storage[q.Len()-1]
	q.Storage = q.Storage[:q.Len()-1]
	return x
}

func (q *EventHeap) Peek() Event {
	if q.IsEmpty() {
		return nil
	}
	return q.Storage[0]
}
