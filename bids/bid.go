package bids

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	charOffset = 97
	maxLetters = 26
)

var (
	numberOfLetters = big.NewInt(maxLetters)
)

type Bid struct {
	Name  string
	Count int64
}

// GetRandom returns a random Bid from active bids
func GetRandom() *Bid {
	return activeBids.randomBid()
}

// GetStatistics returns both of active and archive bids
func GetStatistics() []*Bid {
	total := make([]*Bid, maxBidsCount.Int64())

	activeBids.RLock()
	copy(total, activeBids.bids)
	activeBids.RUnlock()

	archiveBids.RLock()
	total = append(total, archiveBids.bids...)
	archiveBids.RUnlock()

	return total
}

// CreateNewBid creates a new bid and move a replaced one to archived
func CreateNewBid() {
	bid := generateBig()
	pos := randomPosition()

	activeBids.Lock()
	archive := activeBids.bids[pos]
	activeBids.bids[pos] = bid
	activeBids.Unlock()

	archiveBids.Lock()
	archiveBids.bids = append(archiveBids.bids, archive)
	archiveBids.Unlock()
}

func generateBig() *Bid {
	return &Bid{
		Name:  randomSymbol() + randomSymbol(),
		Count: 0,
	}
}

func randomSymbol() string {
	n, err := rand.Int(rand.Reader, numberOfLetters)
	if err != nil {
		log.Printf("Can't generate random number: %s", err)
	}
	return string(n.Int64() + charOffset)
}

func randomPosition() int {
	n, err := rand.Int(rand.Reader, maxBidsCount)
	if err != nil {
		log.Printf("Can't generate random position: %s", err)
	}
	return int(n.Int64())
}
