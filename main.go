package main

import (
	"fmt"
	"github.com/vmihailenco/redis"
	"os"
	"strconv"
	"time"
)

var (
	client *redis.Client
	max_pool int
	cur_pool int
)

func main() {
	redisurl := os.Getenv("REDIS_HOST")
	redispass := os.Getenv("REDIS_PASSWORD")
	
	poele_max, err := strconv.ParseInt(os.Getenv("POELE_MAX"), 10, 0)
	if err != nil { panic(err) }
	max_pool = int(poele_max)
	cur_pool = 0
	
	client = redis.NewTCPClient(redisurl, redispass, -1)
	defer client.Close()
	defer fmt.Println("Exiting...")
	
	fmt.Println("Waiting on queue...")
	for {
		if cur_pool < max_pool {
			jobReq := client.BRPopLPush("box.processing", "box.todo", 0)
			if err := jobReq.Err(); err != nil {
				fmt.Printf("!! Redis: %s\n", err)
				continue
			}
			
			cur_pool += 1
			fmt.Printf("[%i/%i] Processing job: %s\n",
			  cur_pool, max_pool, jobReq.Val())
			go process(jobReq.Val())
		}
	}
}

func process(jobId string) {
  time.Sleep(1)
	cur_pool -= 1
}