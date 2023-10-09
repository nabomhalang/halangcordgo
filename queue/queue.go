package queue

import "sync"

func NewQueue() Queue {
	return Queue{
		queue: make([]Elements, 0),
		rw:    &sync.RWMutex{},
	}
}

func (q *Queue) IsEmpty() bool {
	q.rw.RLock()
	defer q.rw.RUnlock()

	return len(q.queue) == 0
}

func (q *Queue) GetFirst() *Elements {
	q.rw.RLock()
	defer q.rw.RUnlock()

	if len(q.queue) == 0 {
		return nil
	}

	return &q.queue[0]
}

func (q *Queue) Add(e ...Elements) {
	q.rw.Lock()
	defer q.rw.Unlock()

	q.queue = append(q.queue, e...)
}

func (q *Queue) AddPriority(e ...Elements) {
	q.rw.Lock()
	defer q.rw.Unlock()

	if len(q.queue) == 0 {
		q.queue = append(q.queue, e...)
	} else {
		q.queue = append(q.queue[:1], append(e, q.queue[1:]...)...)
	}
}

func (q *Queue) RemoveFirst() {
	q.rw.Lock()
	defer q.rw.Unlock()

	if len(q.queue) != 0 {
		q.queue = q.queue[1:]
	}
}

func (q *Queue) GetAll() []Elements {
	q.rw.RLock()
	defer q.rw.RUnlock()

	queueCopy := make([]Elements, len(q.queue))

	copy(queueCopy, q.queue)

	return queueCopy
}

func (q *Queue) Clear() {
	q.rw.Lock()
	defer q.rw.Unlock()

	q.queue = make([]Elements, 0)
}

func (q *Queue) ModifyFirst(f func(*Elements)) {
	q.rw.Lock()
	defer q.rw.Unlock()

	if len(q.queue) != 0 {
		f(&q.queue[0])
	}
}
