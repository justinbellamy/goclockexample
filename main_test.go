// Copyright (c) 2020 Justin Bellamy

package main

import (
	"reflect"
	"testing"
	"time"
)

func TestRunClockUnchanged(t *testing.T) {
	// tests basic mechanics of the clock function

	// define three types of payloads to use
	// seconds = 1 sec iterator
	// minutes = 2 sec iterator (change to 60 for testing per minute)
	// hours = 4 sec iterator (change to 3600 for testing 3 hours)
	sec := payload{'s', "second", 1}
	min := payload{'m', "minute", 2}
	hr := payload{'h', "hour", 4}

	// // create map of the payloads
	p := map[int]payload{0: hr, 1: min, 2: sec}

	iterator := 1
	lifetime := 5

	// map of expected values. can scale this to unlimited time
	m1 := make(map[int]string)
	m1[1] = "second"
	m1[2] = "minute"
	m1[3] = "second"
	m1[4] = "hour"
	m1[5] = "second"

	logger := make(map[int]string)
	ch := make(chan struct{})
	go runClock(ch, lifetime, iterator, p, logger)
	ch <- struct{}{}
	eq := reflect.DeepEqual(m1, logger)
	if !eq {
		t.Error("Test failed for run clock.")
	}
}

func TestRunClockChanged(t *testing.T) {
	// tests basic mechanics of the clock function which output is altered after program starts

	// define three types of payloads to use
	// seconds = 1 sec iterator
	// minutes = 2 sec iterator (change to 60 for testing per minute)
	// hours = 4 sec iterator (change to 3600 for testing 3 hours)
	sec := payload{'s', "second", 1}
	min := payload{'m', "minute", 2}
	hr := payload{'h', "hour", 4}

	// // create map of the payloads
	p := map[int]payload{0: hr, 1: min, 2: sec}

	iterator := 1
	lifetime := 5

	// map of expected values. can scale this to unlimited time
	m2 := make(map[int]string)
	m2[1] = "second"
	m2[2] = "minute"
	m2[3] = "second"
	m2[4] = "alert" // we change this value after the application starts to simulate console change
	m2[5] = "second"

	logger := make(map[int]string)
	ch := make(chan struct{})
	go runClock(ch, lifetime, iterator, p, logger)
	//simulate changing output via console after 1 second while app is running...
	go func() {
		time.Sleep(1 * time.Second)
		p[0] = payload{'h', "alert", 4}
	}()
	ch <- struct{}{}

	eq := reflect.DeepEqual(m2, logger)
	if !eq {
		t.Error("Test failed for run clock.")
	}
}

func TestCalculateOutputResponse(t *testing.T) {
	sec := payload{'s', "second", 1}
	min := payload{'m', "minute", 60}
	hr := payload{'h', "hour", 3600}

	// create map of the payloads
	m3 := map[int]payload{0: hr, 1: min, 2: sec}

	secExpected := "second"
	secActual := calculateOutputResponse(1, m3)
	if secActual != secExpected {
		t.Error("Test failed for seconds iterator.")
	}

	minExpected := "minute"
	minActual := calculateOutputResponse(60, m3)
	if minActual != minExpected {
		t.Error("Test failed for minutes iterator.")
	}

	hrExpected := "hour"
	hrActual := calculateOutputResponse(3600, m3)
	if hrActual != hrExpected {
		t.Error("Test failed for hours iterator.")
	}
}
