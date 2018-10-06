package Worker

import (
	"../Analyze"
	"crypto/rand"
	"fmt"
	"github.com/toorop/go-bittrex"
	"io"
	"strconv"
	"sync"
	"time"
)

var PoolWorker = make(map[string]*Worker)

type Worker struct {
	ID               string                    `json:"id"`
	InTradeStrategy  *Analyze.AnalyzerInTrade  `json:"in_trade_strategy"`
	OutTradeStrategy *Analyze.AnalyzerOutTrade `json:"out_trade_strategy"`
	StartBTCCash     float64                   `json:"start_btc_cash"`
	AvailableBTCCash float64                   `json:"available__btc_cash"`
	BuyActiveMarket  string                    `json:"active_markets"`
	BuyOrder         *bittrex.Order2           `json:"buy_order"`
	AltBalances      map[string]*Alt           `json:"alt_balances"`
	Log              []LogWorker               `json:"log"`
	Fee              float64                   `json:"fee"`
	mx               sync.Mutex
}

func (worker *Worker) Run(fee float64) bool {
	if worker.StartBTCCash >= 0.00075 && worker.InTradeStrategy != nil && worker.OutTradeStrategy != nil {

		worker.Fee = fee
		worker.AvailableBTCCash = worker.StartBTCCash
		worker.AltBalances = make(map[string]*Alt)
		worker.Log = make([]LogWorker, 0)

		go worker.TradeBuyBot()
		go worker.TradeSellBot()

		worker.ID = newUUID()
		PoolWorker[worker.ID] = worker

		worker.AddLog("я родился :)")
		return true
	} else {
		return false
	}
}

type Alt struct {
	AltName         string          `json:"alt_name"`
	Balance         float64         `json:"balance"`
	BuyPrice        float64         `json:"buy_price"`         // цена за которую купил эту пачку
	ProfitPrice     float64         `json:"profit_price"`      // мин цена продажи что бы выйти в 0
	GrowProfitPrice float64         `json:"grow_profit_price"` // нарастающий профит
	TopAsc          float64         `json:"top_asc"`
	SellOrder       *bittrex.Order2 `json:"sell_order"` // ордер в пачке если он есть
}

func (worker *Worker) AddAlt(altName string, balance, buyPrice, profitPrice float64) {

	worker.mx.Lock()
	defer worker.mx.Unlock()

	alt := &Alt{AltName: altName, Balance: balance, BuyPrice: buyPrice, ProfitPrice: profitPrice, GrowProfitPrice: profitPrice}
	// немного ебаный ключь мапы ¯\_(ツ)_/¯
	worker.AltBalances[altName+":"+strconv.FormatFloat(balance, 'f', 6, 64)+":"+strconv.FormatFloat(buyPrice, 'f', 6, 64)] = alt
}

func (worker *Worker) RemoveAlt(alt *Alt) {

	worker.mx.Lock()
	defer worker.mx.Unlock()

	key := alt.AltName + ":" + strconv.FormatFloat(alt.Balance, 'f', 6, 64) + ":" + strconv.FormatFloat(alt.BuyPrice, 'f', 6, 64)
	delete(worker.AltBalances, key)
}

func (worker *Worker) AddLog(log string) {
	logWorker := LogWorker{Time: time.Now(), Log: log}

	worker.Log = append(worker.Log, logWorker)
}

type LogWorker struct {
	Time time.Time `json:"time"`
	Log  string    `json:"log"`
}

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
