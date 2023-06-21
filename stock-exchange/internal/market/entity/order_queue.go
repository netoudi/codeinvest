package entity

type OrderQueue struct {
	Orders []*Order
}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{
		Orders: []*Order{},
	}
}

func (q *OrderQueue) Less(i int, j int) bool {
	return q.Orders[i].Price < q.Orders[j].Price
}

func (q *OrderQueue) Swap(i int, j int) {
	q.Orders[i], q.Orders[j] = q.Orders[j], q.Orders[i]
}

func (q *OrderQueue) Len() int {
	return len(q.Orders)
}

func (q *OrderQueue) Push(x interface{}) {
	q.Orders = append(q.Orders, x.(*Order))
}

func (q *OrderQueue) Pop() interface{} {
	old := q.Orders
	n := len(old)
	item := old[n-1]
	q.Orders = old[0 : n-1]
	return item
}
