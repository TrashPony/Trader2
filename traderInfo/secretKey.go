package traderInfo

import "github.com/toorop/go-bittrex"

const (
	API_KEY    = "1"
	API_SECRET = "1"
)

var bittrexConnect *bittrex.Bittrex

func GetBittrex() *bittrex.Bittrex {

	if bittrexConnect == nil {
		bittrexConnect = bittrex.New(API_KEY, API_SECRET)
	}

	return bittrexConnect
}
