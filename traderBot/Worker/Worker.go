package Worker

import (
	"../Analyze"
	"github.com/toorop/go-bittrex"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	ID               string                    `json:"id"`
	InTradeStrategy  *Analyze.AnalyzerInTrade  `json:"in_trade_strategy"`
	OutTradeStrategy *Analyze.AnalyzerOutTrade `json:"out_trade_strategy"`
	TradeStrategy    string                    `json:"trade_strategy"`
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

		if worker.TradeStrategy == "Fast" {
			go worker.FastTradeBuy()
			go worker.FastTradeSell()
		}

		if worker.TradeStrategy == "Slow" {
			// TODO более детальный анализ данных, но стратегия торгов таже
			//go worker.FastTradeBuy()
			//go worker.FastTradeSell()
		}

		if worker.TradeStrategy == "VerySlow" {
			go worker.VerySlowTradeBuy()
			go worker.FastTradeSell()
		}

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
	Sell            bool            `json:"sell"`       // уже продана
	SellRate        float64         `json:"sell_rate"`  // цена за которую продал бот
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
