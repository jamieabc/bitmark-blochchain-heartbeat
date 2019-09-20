package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitmark-inc/logger"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
)

var books []string

func init() {
	file, err := os.Open("books.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		books = append(books, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	config, err := newConfig()
	if nil != err {
		fmt.Printf("parse config: %s", err.Error())
		return
	}

	err = logger.Initialise(config.Logging)
	if nil != err {
		fmt.Printf("initialise logger with error: %s\n", err)
		return
	}

	sdk.Init(newSdkConfig(config))

	accounts, err := restoreAccountFromRecoveryPhrase(config.RecoveryPhrases)
	if nil != err {
		fmt.Printf("restore accoutn error: %s", err.Error())
		return
	}

	fmt.Println("Start heartbeat...")

	fp := newFinancePlanner(config)
	duration := fp.actionDuration()
	fmt.Printf("action duration: %v\n", duration)
	shutdownChan := make(chan struct{})
	doPeriodicTasks(taskInfo{
		duration:     duration,
		accounts:     accounts,
		shutdownChan: shutdownChan,
		config:       config,
	})

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	shutdownChan <- struct{}{}
	fmt.Printf("received signal: %v\n", sig)
	fmt.Println("Terminating...")
}
