function UpdateStatus(jsonMessage) {
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
    let efficiency = 100 - (bot.start_btc_cash * 100 / bot.available__btc_cash);
    tdEfficiency.innerHTML = efficiency;

    if (efficiency < 0) {
        tdEfficiency.className = "Failed"
    } else if (efficiency > 0) {
        tdEfficiency.className = "Success"
    } else if (efficiency === 0 || isNaN(efficiency)) {
        tdEfficiency.className = "Value"
    }
}

function UpdateIcon(id, bot) {
    let botIcon = document.getElementById("botIcon" + id);
    botIcon.style.background = 'url("../../img/' + bot.in_trade_strategy.name + '.jpg") no-repeat';
    botIcon.style.backgroundSize = "100%"
}

function UpdateCashTable(id, bot) {

    let startBTCCash = document.getElementById("StartBTCCash" + id);
    startBTCCash.innerHTML = bot.start_btc_cash;

    let availableBTCCash = document.getElementById("AvailableBTCCash" + id);
    availableBTCCash.innerHTML = bot.available__btc_cash;

    let inTradeStrategy = document.getElementById("InTradeStrategy" + id);
    inTradeStrategy.innerHTML = bot.in_trade_strategy.name;

    let outTradeStrategy = document.getElementById("OutTradeStrategy" + id);
    outTradeStrategy.innerHTML = bot.out_trade_strategy.name;
}

function UpdateAltTable(id, bot) {

    let tradeSellStatus = document.getElementById("TradeSellStatus" + id);

    for (let alt in bot.alt_balances) {
        if (bot.alt_balances.hasOwnProperty(alt)) {
            //TODO заполнение таблицы алтов
        }
    }
}

function UpdateBuyStatus(id, bot) {
    let buyStatusBlock = document.getElementById("BuyStatus" + id);

    if (bot.available__btc_cash < 0.0005) {
        buyStatusBlock.innerHTML = "Нет денег :(";
        return
    }

    if (bot.active_markets) {
        //TODO
    } else {
        buyStatusBlock.innerHTML = "В активном поиске"
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
        }
    }
}