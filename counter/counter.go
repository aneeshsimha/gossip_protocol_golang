package counter

type Counter struct {
	countNum uint64
	Count    chan uint64
	// when deleting, call close(Counter.Count)
}

func New() *Counter {
	c := &Counter{
		countNum: 1, // start at 1
		Count:    make(chan uint64),
	}
	go c.increment()
	return c
}

func (c *Counter) increment() {
	for {
		c.Count <- c.countNum
		c.countNum += 1
	}
}
