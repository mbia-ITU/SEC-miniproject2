package utilities

import (
	"fmt"
	"math/rand"
)

func PortToString(port int) string {
	return fmt.Sprintf(":%d", port)
}

func CreateShares(p int, data int, amount int) []int {
	var shares []int
	var totalShares int

	for i := 0; i < amount-1; i++ {
		share := rand.Intn(p-1) + 1
		shares = append(shares, share)
		totalShares += share
	}

	shares = append(shares, data-totalShares)

	return shares
}
