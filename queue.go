package diode_queue_bench

import (
	"log"
	"sync"
	"time"
)

// Call is rpc call
type Call struct {
	RPC  string
	Size uint
	// SimRoundTripTime is round trip time per rpc
	SimRoundTripTime time.Duration
}

// Queue is rpc queue in diode client
type Queue struct {
	// SimTicketBytes is when to send a diode ticket after transfering bytes
	SimTicketBytes uint
	// SimRpcBytes is rpc packet size
	SimRPCBytes uint
	// TotalBytes is total bytes that sent to server
	TotalBytes uint
	// Q is rpc queue
	Q              chan Call
	done           chan struct{}
	wg             sync.WaitGroup
	ticketsCount   uint
	totalBytesSent uint
	Verbose        bool
}

// Wait the queue
func (q *Queue) Wait() {
	q.wg.Wait()
}

// Close the queue
func (q *Queue) Close() {
	close(q.done)
}

// Enqueue push call into queue
func (q *Queue) Enqueue(c Call) {
	q.wg.Add(1)
	q.Q <- c
}

// Start the queue
func (q *Queue) Start() {
	for {
		select {
		case <-q.done:
			return
		case c := <-q.Q:
			q.totalBytesSent += c.Size
			q.TotalBytes += c.Size
			if q.TotalBytes >= q.SimTicketBytes {
				q.TotalBytes -= q.SimTicketBytes
				q.Enqueue(Call{
					RPC:              "ticket",
					Size:             133,
					SimRoundTripTime: 100 * time.Millisecond,
				})
				q.ticketsCount++
			}
			if q.Verbose {
				log.Printf("Receive %d bytes of %s", c.Size, c.RPC)
			}
			time.Sleep(c.SimRoundTripTime)
			q.wg.Done()
		}
	}
}
