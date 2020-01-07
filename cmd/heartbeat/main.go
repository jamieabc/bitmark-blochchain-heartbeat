package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bitmark-inc/logger"

	bitmarkSDK "github.com/bitmark-inc/bitmark-sdk-go"
	financePlanner "github.com/jamieabc/bitmark-blochchain-heartbeat/internal/finance_planner"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/internal/periodic"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/parser"
	"github.com/jamieabc/bitmark-blochchain-heartbeat/pkg/sdk"
)

const (
	bookFilename   = "books.txt"
	configFilename = "sdk.conf"
)

var books []string

func init() {
	file, err := os.Open(bookFilename)
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
	p, err := parser.NewParser(configFilename)
	if nil != err {
		fmt.Printf("new parser with error: %s", err)
		return
	}

	config, err := p.Parse()
	if err != nil {
		fmt.Printf("parse config with error: %s", err)
		return
	}

	err = logger.Initialise(config.Logging)
	if nil != err {
		fmt.Printf("initialise logger with error: %s\n", err)
		return
	}

	bitmarkSDK.Init(sdk.NewSDKConfig(config, books))

	accounts, err := sdk.RestoreAccountFromRecoveryPhrase(config.RecoveryPhrases)
	if nil != err {
		fmt.Printf("restore accoutn error: %s", err.Error())
		return
	}

	fmt.Println("Start heartbeat...")

	fp := financePlanner.NewFinancePlanner(config)
	duration := fp.ActionDuration()
	fmt.Printf("action duration: %v\n", duration)
	shutdownChan := make(chan struct{})
	per := periodic.NewPeriodic(duration, accounts, shutdownChan, config)
	per.Do()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	shutdownChan <- struct{}{}
	fmt.Printf("received signal: %v\n", sig)
	fmt.Println("Terminating...")
}
