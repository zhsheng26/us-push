package lib

type (
	Queue struct {
		start, end *pot
		length     int
	}
	pot struct {
		value interface{}
		next  *pot
	}
)

// Create a new queue
func NewQueue() *Queue {
	return &Queue{nil, nil, 0}
}

// Take the next item off the front of the queue
func (queue *Queue) Dequeue() interface{} {
	if queue.length == 0 {
		return nil
	}
	n := queue.start
	if queue.length == 1 {
		queue.start = nil
		queue.end = nil
	} else {
		queue.start = queue.start.next
	}
	queue.length--
	return n.value
}

// Put an item on the end of a queue
func (queue *Queue) Enqueue(value interface{}) {
	pot := &pot{value, nil}
	if queue.length == 0 {
		queue.start = pot
		queue.end = pot
	} else {
		queue.end.next = pot
		queue.end = pot
	}
	queue.length++
}

// Return the number of items in the queue
func (queue *Queue) Len() int {
	return queue.length
}

// Return the first item in the queue without removing it
func (queue *Queue) Peek() interface{} {
	if queue.length == 0 {
		return nil
	}
	return queue.start.value
}
