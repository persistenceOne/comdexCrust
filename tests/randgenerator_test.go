package tests

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var ticketCounter int64

func TestRand(t *testing.T) {
	var numbers = make(map[string]bool)
	count := 0
	start := int(time.Now().Unix())
	for {
		now := 10000000 + int(time.Now().UnixNano())%89999999
		atomic.AddInt64(&ticketCounter, 1)
		ticketCounter := int(ticketCounter)
		ticket := 10000 + ticketCounter%89999
		number := 10000 + rand.Intn(89999)
		randN := "TR" + strconv.Itoa(ticket) + strconv.Itoa(now) + strconv.Itoa(number)
		//fmt.Println(string(randN))
		if _, ok := numbers[randN]; ok {
			fmt.Println(string(randN))
			require.Nil(t, count)
		} else {
			count++
			end := int(time.Now().Unix())
			fmt.Println(count / (end - start + 1))
			numbers[string(randN)] = true
		}
	}
}
