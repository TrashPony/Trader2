package Worker

import (
	"crypto/rand"
	"fmt"
	"io"
)

var PoolWorker = make(map[string]*Worker)

func newUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	uuid[8] = uuid[8]&^0xc0 | 0x80
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func GetPoolWorker() map[string]*Worker {
	return PoolWorker
}

func CheckEmploymentMarket(marketName string) bool {
	for _, worker := range PoolWorker {
		if worker.BuyActiveMarket == marketName {
			return true
		}
	}

	return false
}
