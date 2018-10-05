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

function UpdateLog(id, bot) {

    for (let i = 0; i < bot.log.length; i++) {

    }
}