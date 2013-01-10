package main

import (
	"fmt"
	"github.com/vmihailenco/redis"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var client *redis.Client

func main() {
	quit := make(chan int, 1)
	go handleSignals(quit)

	redisurl := os.Getenv("REDIS_HOST")
	redispass := os.Getenv("REDIS_PASSWORD")

	ncpu := runtime.NumCPU()
	runtime.GOMAXPROCS(ncpu)

	client = redis.NewTCPClient(redisurl, redispass, -1)
	defer client.Close()

	queue := make(chan string, ncpu)
	for i := 0; i < ncpu; i++ {
		go doWork(i, queue)
	}

	go giveWork(queue)

	<-quit
	for i := 0; i < ncpu; i++ {
		queue <- ""
	}
}

func giveWork(queue chan string) {
	from := "box.todo"
	to := "box.processing"

	for {
		item := client.BRPopLPush(from, to, 0)
		if err := item.Err(); err != nil {
			fmt.Printf("!!! Redis: %s\n", err)
			continue
		}

		val := item.Val()
		fmt.Printf(">>> Crèpes au bacon, order %s, coming up!\n", val)
		queue <- val
	}
}

func doWork(id int, queue chan string) {
	fmt.Printf(">>> Firing up gas Nº%d and getting a clean pan…\n", id)

	for {
		item := <-queue
		if item == "" {
			break
		}

		fmt.Printf("[%d] Cooking: %s\n", id, item)
		then := time.Now()

		ret := process(item)

		dur1 := uDuration(time.Since(then))

		client.LPush("box.done", item)
		client.LRem("box.processing", -1, item)

		dur2 := uDuration(time.Since(then))

		fmt.Printf("[%d] Done in (%s/%s): %v\n", id, dur1, dur2, ret)
	}
}

func uDuration(d time.Duration) string {
	return strings.Replace(d.String(), "u", "µ", 1)
}

func process(item string) int64 {
	n, _ := strconv.ParseInt(item, 10, 64)
	return fib(n)
}

func fib(n int64) int64 {
	if n < 2 {
		return n
	}

	return fib(n-1) + fib(n-2)
}
