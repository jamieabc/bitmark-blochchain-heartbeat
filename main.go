package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sdk "github.com/bitmark-inc/bitmark-sdk-go"
)

const (
	minActionInterval = time.Duration(5) * time.Second
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

	sdk.Init(newSdkConfig(config))

	accounts, err := restoreAccountFromRecoveryPhrase(config.RecoveryPhrases)
	if nil != err {
		fmt.Printf("restore accoutn error: %s", err.Error())
		return
	}

	fmt.Println("Start heartbeat...")

	fp := newFinancePlanner(config, minActionInterval)
	actionInterval := fp.actionInterval()
	fmt.Printf("action duration: %v\n", actionInterval)
	runBackground(actionInterval, accounts)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	fmt.Printf("received signal: %v\n", sig)
	fmt.Println("Terminating...")
}
