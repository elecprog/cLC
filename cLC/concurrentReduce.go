package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/ElecProg/lamcalc"
)

// Here we implement a system to call a function and make it stoppable on ^+C
// The function isn't really stopped, rather it's result is discarded, and we
// don't wait for it anymore either.

func concurrentReduce(term lamcalc.Term) (lamcalc.Term, error) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	defer signal.Stop(sig)

	res := make(chan lamcalc.Term, 2)

	donen := make(chan bool, 1)
	go lamcalc.ConcNorReduce(term, res, donen)

	donea := make(chan bool, 1)
	go lamcalc.ConcAorReduce(term, res, donea)

	select {
	case res := <-res:
		// Send stop signals
		donen <- true
		donea <- true

		return res, nil

	case <-sig:
		// Remove the '^C' from the terminal:
		fmt.Print("\b\b")

		// Stop computations
		donen <- true
		donea <- true

		return nil, errors.New("keyboard interrupt")
	}
}