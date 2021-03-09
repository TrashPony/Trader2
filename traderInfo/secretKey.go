package traderInfo

import "github.com/toorop/go-bittrex"

const (
	API_KEY    = "52fd7aa114c84ef2bd013e5ad9af288f"
	API_SECRET = "93801485474e497f832271bac32eea15"
)

var bittrexConnect *bittrex.Bittrex

func GetBittrex() *bittrex.Bittrex {

	if bittrexConnect == nil {
		bittrexConnect = bittrex.New(API_KEY, API_SECRET)
	}

	return bittrexConnect
}
