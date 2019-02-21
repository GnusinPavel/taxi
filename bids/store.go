package bids

import (
	"math/big"
	"sync"
	"sync/atomic"
)

const (
	maxBids = 50
)

var (
	maxBidsCount = big.NewInt(maxBids)

	activeBids = storage{
		RWMutex: sync.RWMutex{},
		bids:    make([]*Bid, maxBidsCount.Int64()),
	}

	archiveBids = storage{
		RWMutex: sync.RWMutex{},
		bids:    make([]*Bid, 0),
	}
)

type storage struct {
	sync.RWMutex
	bids []*Bid
}

func (s *storage) randomBid() *Bid {
	bid := s.bids[randomPosition()]
	atomic.AddInt64(&bid.Count, 1)
	return bid
}

func init() {
	for i := 0; i < int(maxBidsCount.Int64()); i++ {
		activeBids.bids[i] = generateBig()
	}
}
