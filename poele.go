package poele

import (
	"fmt"
	"github.com/vmihailenco/redis"
	"log"
	"runtime"
	"strings"
	"time"
)

type Poële struct {
	rds         *redis.Client
	prefix      string
	queue       chan string
	concurrency int
	process     func(string) interface{}
}

func New(rds *redis.Client, pfx string, proc func(string) interface{}) Poële {
	ncpu := runtime.NumCPU()
	ch := make(chan string, ncpu)
	return Poële{rds, pfx, ch, ncpu, proc}
}

func (p Poële) Serve() {
	log.Printf(">>> Opening the kitchen with ‘%s’ as menu du jour.\n", p.prefix)
	runtime.GOMAXPROCS(p.concurrency)
	for i := 0; i < p.concurrency; i++ {
		go p.doWork(i)
	}

	go p.giveWork()
}

func (p Poële) giveWork() {
	from := fmt.Sprintf("poele.%s.todo", p.prefix)
	to := fmt.Sprintf("poele.%s.processing", p.prefix)

	for {
		item := p.rds.BRPopLPush(from, to, 0)
		if err := item.Err(); err != nil {
			log.Printf("!!! Redis: %s\n", err)
			continue
		}

		val := item.Val()
		log.Printf(">>> Crèpes au bacon, order %s, coming up!\n", val)
		p.queue <- val
	}
}

func (p Poële) doWork(id int) {
	log.Printf(">>> Firing up gas Nº%d and getting a clean pan…\n", id)

	for {
		item := <-p.queue
		if item == "" {
			break
		}

		log.Printf("[%d] Cooking: %s\n", id, item)
		then := time.Now()

		ret := p.process(item)

		dur1 := uDuration(time.Since(then))

		p.rds.LPush(fmt.Sprintf("poele.%s.done", p.prefix), item)
		p.rds.LRem(fmt.Sprintf("poele.%s.processing", p.prefix), -1, item)

		dur2 := uDuration(time.Since(then))

		log.Printf("[%d] Done in (%s/%s): %v\n", id, dur1, dur2, ret)
	}
}

func uDuration(d time.Duration) string {
	return strings.Replace(d.String(), "u", "µ", 1)
}
