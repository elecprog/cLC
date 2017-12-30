package LamCalc

import (
	"errors"
)

// AorReduce reduces a lambda expression using applicative order
func (lx Appl) AorReduce() (Term, error) {
	nw := lx.aorReduceOnce()

	for c := 1; nw.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		nw = nw.aorReduceOnce()
	}

	return nw.etaReduce(), nil
}

// aorReduceOnce reduces a lambda application once
func (lx Appl) aorReduceOnce() Term {
	if !lx[1].canReduce() {
		switch fst := lx[0].(type) {
		case Abst:
			return fst.betaReduce(lx[1])

		default:
			return Appl{lx[0].aorReduceOnce(), lx[1]}
		}
	}

	return Appl{lx[0], lx[1].aorReduceOnce()}
}

// AorReduce reduces a lambda abstraction using applicative order
func (la Abst) AorReduce() (Term, error) {
	nw := la.aorReduceOnce()

	for c := 1; nw.canReduce(); c++ {
		if c == MaxReductions {
			return nil, errors.New("exeeded maximum amount of reductions")
		}

		nw = nw.aorReduceOnce()
	}

	return nw.etaReduce(), nil
}

// aorReduceOnce reduces a lambda abstraction once
func (la Abst) aorReduceOnce() Term {
	return Abst{la[0].aorReduceOnce()}
}

// AorReduce returns the variable itself
func (lv Var) AorReduce() (Term, error) {
	return lv, nil
}

// aorReduceOnce reduces a lambda variable once
func (lv Var) aorReduceOnce() Term {
	return lv
}