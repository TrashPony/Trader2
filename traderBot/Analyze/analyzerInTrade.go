package Analyze

import "../../traderInfo"

/**
 * Created by trash on 09.08.2017.
 * Надо сделать еще 1 проверку 1го и второго оредара на адекватность первого+
 * -
 * Надо смотреть разностьи между продажным и покупным ордером
 * Возможно имеет смысл продать валюту сразу а не ждать пока ее кто то купит
 * И смотреть процентное соотношение между ними может продать сразу будет даже выгоднее+
 * -
 * Дописать стратегию хайпа, если валюта растет то не обязательно что она начнет прям тут же падать
 * Максимум я потеряю 2% а выйграть могу до овер дохуя+-
 * -
 * Надо собирать статистику суммы по ордерам продажи и покупки если сумма продажи сильно больше
 * то вероятнее всего курс просядет, если покупака больше то значит ее скупают и курс вероятно выростет+
 * -
 * Надо реализовать деление биткоинов между альтами, что бы бот не тратил сразу весь банк на 1 валюту
 * -
 * Если есть ордер на продажу то надо кидать на калькулятор
 * -
 * Мониторить последние 10 ордеров если там есть профитные то не сливать монету
 * -
 * Реализовать проверка на каждой итерации на рентабильность дальнейших торгов, может иметь смысл продать валюту если получить в плюс сразу+
 * -
 * avg начало и конца истории в обоих позициях
 * -
 * сделать анализ рынка в момент продажи на адекватность сброса
 * -
 */

type AnalyzerInTrade struct {
	Name string `json:"name"`
}

func (analyze *AnalyzerInTrade) Analyze(market *traderInfo.Market) (buy bool, fast bool) {

	Bid, _ := market.MarketSummary.Bid.Float64()
	Ask, _ := market.MarketSummary.Ask.Float64()
	Low, _ := market.MarketSummary.Low.Float64()
	High, _ := market.MarketSummary.High.Float64()
	FirstOrderBuy, _ := market.OrdersBuy[1].Rate.Float64()
	SecondOrderBuy, _ := market.OrdersBuy[0].Rate.Float64()
	Last, _ := market.MarketSummary.Last.Float64()

	avgLowHigh := (Low + High) / 2
	differenceAskBind := percentageCalculator(Bid, Ask)
	demand := percentageCalculator(float64(market.MarketSummary.OpenSellOrders), float64(market.MarketSummary.OpenBuyOrders))
	second := percentageCalculator(FirstOrderBuy, SecondOrderBuy)

	// boolean sumCap = market.summ25QuantityOrderBidsBook > (market.summ25QuantityOrderAsksBook * 1.5);
	// boolean historyProf = readHistoryMarket(market, params, log);

	openOrdersCheck := demand > 70                     // спрос предложение
	lastPriceCheck := Last >= avgLowHigh               // последний ордер ниже средней по суткам
	differenceAskBindCheck := differenceAskBind > 0.55 // комисия
	secondCheck := second < 0.10                       // Второй оредар на покупку

	return secondCheck && differenceAskBindCheck && openOrdersCheck && lastPriceCheck, false
}

func percentageCalculator(a, b float64) (result float64) {
	result = 100 - (a * 100 / b)
	return result
}