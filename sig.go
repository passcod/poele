package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(quit chan int) {
	kill := make(chan os.Signal)
	signal.Notify(kill,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGSTOP)

	for {
		s := <-kill
		fmt.Printf("!!! The chef is tired, the kitchen got %v'd.\n", s)
		quit <- 1
	}
}
