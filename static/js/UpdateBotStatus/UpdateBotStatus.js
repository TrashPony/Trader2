function UpdateStatus(jsonMessage) {
    console.log(jsonMessage);
    let botsData = JSON.parse(jsonMessage);

    for (let bot in botsData.Workers) {
        if (botsData.Workers.hasOwnProperty(bot)) {
            let botBox = document.getElementById(bot);
            if (botBox) {
                Update(bot, botsData.Workers[bot]);
            } else {
                CreateBotBox(bot);
                Update(bot, botsData.Workers[bot]);
            }
        }
    }
}

function Update(id, bot) {
    UpdateHeaderInfo(id, bot);
    UpdateIcon(id, bot);
    UpdateCashTable(id, bot);
    UpdateAltTable(id, bot);
    UpdateBuyStatus(id, bot);
    UpdateLog(id, bot);
}

function UpdateHeaderInfo(id, bot) {
    let tdName = document.getElementById("name" + id);
    tdName.innerHTML = id.split("-")[0]; // а то уиды длинные не влезают)

    let tdEfficiency = document.getElementById("Efficiency" + id);

    let botBTCAvailable = bot.available__btc_cash;

    for (let alt in bot.alt_balances) {
        if (bot.alt_balances.hasOwnProperty(alt) && !bot.alt_balances[alt].sell) {
            let altOrder = botBTCAvailable + (bot.alt_balances[alt].balance * bot.alt_balances[alt].top_asc);
            let fee = 0.0075 * altOrder;
            altOrder = altOrder - fee;
            botBTCAvailable += altOrder;
        }
    }

    let efficiency = 100 - (bot.start_btc_cash * 100 / botBTCAvailable);

    if (efficiency < 0) {
        tdEfficiency.className = "Failed"
    } else if (efficiency > 0) {
        tdEfficiency.className = "Success"
    } else if (efficiency === 0) {
        tdEfficiency.className = "Value"
    }


    tdEfficiency.innerHTML = efficiency.toFixed(3);
}

function UpdateIcon(id, bot) {
    let botIcon = document.getElementById("botIcon" + id);
    botIcon.style.background = 'url("../../img/' + bot.trade_strategy + '.jpg") no-repeat';
    botIcon.style.backgroundSize = "100%"
}

function UpdateCashTable(id, bot) {

    let startBTCCash = document.getElementById("StartBTCCash" + id);
    startBTCCash.innerHTML = bot.start_btc_cash.toFixed(8);

    let altToBTC = 0;
    for (let alt in bot.alt_balances) {
        if (bot.alt_balances.hasOwnProperty(alt) && !bot.alt_balances[alt].sell) {
            altToBTC = altToBTC + (bot.alt_balances[alt].balance * bot.alt_balances[alt].top_asc);
            let fee = (0.25 * altToBTC) / 100;
            altToBTC = altToBTC - fee;
        }
    }

    let availableBTCCash = document.getElementById("AvailableBTCCash" + id);
    availableBTCCash.innerHTML = bot.available__btc_cash.toFixed(8) + " <sub>+" + altToBTC.toFixed(8) + "</sub>";

    let inTradeStrategy = document.getElementById("InTradeStrategy" + id);
    inTradeStrategy.innerHTML = bot.in_trade_strategy.name;

    let outTradeStrategy = document.getElementById("OutTradeStrategy" + id);
    outTradeStrategy.innerHTML = bot.out_trade_strategy.name;

    let TradeStrategy = document.getElementById("TradeStrategy" + id);
    TradeStrategy.innerHTML = bot.trade_strategy;
}

function UpdateAltTable(id, bot) {

    let tradeSellStatus = document.getElementById("altCash" + id);

    for (let alt in bot.alt_balances) {
        if (bot.alt_balances.hasOwnProperty(alt)) {
            let trAlt = document.getElementById(alt);

            if (!trAlt) {
                trAlt = document.createElement("tr");
                trAlt.id = alt;

                function createTD(value, id) {
                    let td = document.createElement("td");
                    td.innerHTML = value;
                    td.id = id;
                    trAlt.appendChild(td);
                }

                createTD(bot.alt_balances[alt].alt_name, "name" + id + alt);
                createTD(bot.alt_balances[alt].balance.toFixed(8), "balance" + id + alt);
                createTD(bot.alt_balances[alt].buy_price.toFixed(8), "buyPrice" + id + alt);
                createTD(bot.alt_balances[alt].profit_price.toFixed(8), "profitPrice" + id + alt);
                createTD(bot.alt_balances[alt].grow_profit_price.toFixed(8), "growProfitPrice" + id + alt);
                createTD(bot.alt_balances[alt].top_asc.toFixed(8), "topAsc" + id + alt);
                createTD(bot.alt_balances[alt].sell_rate.toFixed(8), "sellRate" + id + alt);
                createTD("", "PercentRate" + id + alt);


                tradeSellStatus.appendChild(trAlt);
            } else {

                document.getElementById("name" + id + alt).innerHTML = bot.alt_balances[alt].alt_name;
                document.getElementById("balance" + id + alt).innerHTML = bot.alt_balances[alt].balance.toFixed(8);
                document.getElementById("buyPrice" + id + alt).innerHTML = bot.alt_balances[alt].buy_price.toFixed(8);
                document.getElementById("profitPrice" + id + alt).innerHTML = bot.alt_balances[alt].profit_price.toFixed(8);
                document.getElementById("growProfitPrice" + id + alt).innerHTML = bot.alt_balances[alt].grow_profit_price.toFixed(8);
                document.getElementById("topAsc" + id + alt).innerHTML = bot.alt_balances[alt].top_asc.toFixed(8);

                if (bot.alt_balances[alt].sell_rate === 0 || bot.alt_balances[alt].sell_rate === bot.alt_balances[alt].profit_price) {
                    document.getElementById("sellRate" + id + alt).innerHTML = "no sell";
                    document.getElementById("sellRate" + id + alt).style.color = "#a5d5ef";
                } else if (bot.alt_balances[alt].sell_rate > bot.alt_balances[alt].profit_price) {
                    document.getElementById("sellRate" + id + alt).innerHTML = bot.alt_balances[alt].sell_rate.toFixed(8);
                    document.getElementById("sellRate" + id + alt).style.color = "#1BFF04"
                } else if (bot.alt_balances[alt].sell_rate < bot.alt_balances[alt].profit_price) {
                    document.getElementById("sellRate" + id + alt).innerHTML = bot.alt_balances[alt].sell_rate.toFixed(8);
                    document.getElementById("sellRate" + id + alt).style.color = "#FF5F61";
                }

                let PercentRate;
                if (bot.alt_balances[alt].sell) {
                    PercentRate = 100 - (bot.alt_balances[alt].profit_price * 100 / bot.alt_balances[alt].sell_rate);
                } else {
                    PercentRate = 100 - (bot.alt_balances[alt].profit_price * 100 / bot.alt_balances[alt].top_asc);
                }

                document.getElementById("PercentRate" + id + alt).innerHTML = PercentRate.toFixed(3);

                if (PercentRate > 0) {
                    document.getElementById("PercentRate" + id + alt).style.color = "#1BFF04";
                } else if (PercentRate < 0) {
                    document.getElementById("PercentRate" + id + alt).style.color = "#FF5F61";
                } else {
                    document.getElementById("PercentRate" + id + alt).style.color = "#a5d5ef";
                }
            }
        }
    }
}

function UpdateBuyStatus(id, bot) {
    let buyStatusBlock = document.getElementById("BuyStatus" + id);

    // if (bot.available__btc_cash < 0.0005) {
    //     buyStatusBlock.innerHTML = "Нет денег :(";
    //     return
    // }
    // выводить полностью ордер по которому покует бот
    if (bot.active_markets && !bot.buy_order) {
        buyStatusBlock.innerHTML = "Анализирую рынок " + bot.active_markets;
    } else if (bot.active_markets && bot.buy_order) {
        buyStatusBlock.innerHTML = "Имею ордер на покупку в " + bot.buy_order.Exchange;
    } else if (!bot.active_markets && !bot.buy_order) {
        buyStatusBlock.innerHTML = "В активном поиске";
    }
}

function UpdateLog(id, bot) {
    let logBlock = document.getElementById("log" + id);

    for (let i = 0; i < bot.log.length; i++) {
        let logRow = document.getElementById("log" + id + ":" + i);

        if (!logRow) {
            let time = document.createElement("span");
            time.className = "timeLog";
            time.innerHTML = bot.log[i].time.split("T")[0] + " - " + bot.log[i].time.substring(11, 19) + " "; // лень думать

            logRow = document.createElement("span");
            logRow.className = "textLog";
            logRow.innerHTML = bot.log[i].log;
            logRow.id = "log" + id + ":" + i;

            logBlock.appendChild(time);
            logBlock.appendChild(logRow);
            logBlock.appendChild(document.createElement("br"));
        }
    }
}