package diode_queue_bench

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func BenchmarkQueue(b *testing.B) {
	b.ResetTimer()
	q := Queue{
		SimTicketBytes: 4000,
		SimRPCBytes:    4122,
		TotalBytes:     0,
		Q:              make(chan Call, 1024),
		done:           make(chan struct{}),
		wg:             sync.WaitGroup{},
	}
	go func() {
		for i := 0; i < b.N; i++ {
			m := (i % 3) + 1
			rpc := fmt.Sprintf("portsend-%d", i)
			q.Enqueue(Call{
				RPC:  rpc,
				Size: 4122 * uint(m),
			})
		}
		q.Wait()
		q.Close()
		log.Printf("total bytes sent: %d\ntotal ticket send: %d\ntotal portsend rpc send: 100\n", q.totalBytesSent, q.ticketsCount)
	}()
	q.Start()
}

func BenchmarkQueue100(b *testing.B) {
	b.ResetTimer()
	q := Queue{
		SimTicketBytes: 4000,
		SimRPCBytes:    4122,
		TotalBytes:     0,
		Q:              make(chan Call, 1024),
		done:           make(chan struct{}),
		wg:             sync.WaitGroup{},
	}
	go func() {
		for i := 0; i < 100; i++ {
			m := (i % 3) + 1
			rpc := fmt.Sprintf("portsend-%d", i)
			q.Enqueue(Call{
				RPC:              rpc,
				Size:             4122 * uint(m),
				SimRoundTripTime: 150 * time.Millisecond,
			})
		}
		q.Wait()
		q.Close()
		log.Printf("total bytes sent: %d\ntotal ticket send: %d\ntotal portsend rpc send: 100\n", q.totalBytesSent, q.ticketsCount)
	}()
	q.Start()
}
