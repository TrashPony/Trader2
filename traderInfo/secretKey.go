package traderInfo

import "github.com/toorop/go-bittrex"

const (
	API_KEY    = "83beb1021eb142dc93c1493059632a16"
	API_SECRET = "8de2bbaa02a44323acb225966469e52b"
)

var bittrexConnect *bittrex.Bittrex

func GetBittrex() *bittrex.Bittrex {

	if bittrexConnect == nil {
		bittrexConnect = bittrex.New(API_KEY, API_SECRET)
	}

	return bittrexConnect
}
