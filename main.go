// Copyright (c) 2020 Justin Bellamy

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type payload struct {
	id        rune
	name      string
	iteration int
}

func main() {
	logger := make(map[int]string)

	// initialize default values

	// total lifetime for the application to run (in seconds)
	lifetime := 10800 //10800 seconds in 3 hours

	// what is the iteration counter we want to use (in seconds)
	iterator := 1

	// iteration payloads
	hr := payload{'h', "hour", 3600}  // 3600 secconds in a hour
	min := payload{'m', "minute", 60} // 60 seconds in a minute
	sec := payload{'s', "second", 1}  // 1 second

	// display console instructions
	fmt.Print(`To update the output at any time, type one of the following letters followed by the new word to change the output, then press enter.
	Type 's' to edit seconds.
	Type 'm' to edit minutes.
	Type 'h' to edit hours.
	Type 'q' and press enter to exit the program.`)

	// create a map of the payloads. index 0 has highest output priority, 1 has second highest, 2 has the third highest, etc.
	m := map[int]payload{0: hr, 1: min, 2: sec}

	// run the clock
	ch := make(chan struct{})
	go runClock(ch, lifetime, iterator, m, logger)

	go func() {
		// handle console input for updating names while application is running
		reader := bufio.NewReader(os.Stdin)
		for {
			cmdString, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			err = handleConsoleInput(cmdString, m)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}()
	ch <- struct{}{}
}

// loops through the time iterator, displays output and terminates the app when it reaches its lifetime
func runClock(ch chan struct{}, lifetime int, iterator int, m map[int]payload, logger map[int]string) {
	startTime := time.Now()
	i := 1 //start logging at iteration/second 1
	for {
		time.Sleep(time.Duration(iterator) * time.Second)
		elapsed := time.Since(startTime) //get duration
		second := int(elapsed.Seconds()) //convert duration to int
		response := calculateOutputResponse(second, m)
		logger[i] = response
		fmt.Printf("%s\n", response)
		i++
		if lifetime <= second {
			<-ch
			break
		}
	}
}

// determines which output gets displayed depending on priority (map key of 0 is highest priority).
func calculateOutputResponse(elapsed int, m map[int]payload) string {
	// cant range here because go doesn't scan map in sequential key order
	for i := 0; i < len(m); i++ {
		t := m[i]
		remainder := elapsed % t.iteration
		if remainder == 0 {
			return t.name
		}
	}
	return ""
}

// hard coded case statement. i would refactor this to scan map name fields
func handleConsoleInput(consoleInput string, m map[int]payload) error {
	consoleInput = strings.TrimSuffix(consoleInput, "\n")
	consoleInputArgs := strings.Fields(consoleInput)
	if consoleInputArgs[0] == "q" {
		os.Exit(0)
	}
	updatedName := consoleInputArgs[1]
	switch consoleInputArgs[0] {
	case "h":
		hourPayload := m[0]
		fmt.Printf("\n >>> %s changed to %s <<< \n", hourPayload.name, updatedName)
		m[0] = payload{'h', updatedName, hourPayload.iteration}
	case "m":
		minPayload := m[1]
		fmt.Printf("\n >>> %s changed to %s <<< \n", minPayload.name, updatedName)
		m[1] = payload{'m', updatedName, minPayload.iteration}
	case "s":
		secPayload := m[2]
		fmt.Printf("\n >>> %s changed to %s <<< \n", secPayload.name, updatedName)
		m[2] = payload{'s', updatedName, secPayload.iteration}
	}
	return nil
}
