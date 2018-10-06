function CreateTradeStatusBlock(id) {
    let tradeStatus = document.createElement("div");
    tradeStatus.className = "TradeStatus";

    let tradeBuyStatus = CreateTradeBuyStatus(id);
    tradeStatus.appendChild(tradeBuyStatus);

    let tradeSellStatus = CreateTradeSellStatus(id);
    tradeStatus.appendChild(tradeSellStatus);

    return tradeStatus;
}

function CreateTradeBuyStatus(id) {
    let tradeBuyStatus = document.createElement("div");
    tradeBuyStatus.className = "TradeBuyStatus";

    let spanHead = document.createElement("span");
    spanHead.innerHTML = "Buy status";
    spanHead.className = "header";
    spanHead.style.marginTop = "5px";
    spanHead.style.display = "block";
    spanHead.style.fontSize = "12px";
    tradeBuyStatus.appendChild(spanHead);

    let statusBlock = document.createElement("div");
    statusBlock.className = "statusBuyBlock";
    statusBlock.id = "BuyStatus" + id;
    tradeBuyStatus.appendChild(statusBlock);

    return tradeBuyStatus;
}

function CreateTradeSellStatus(id) {
    let tradeSellStatus = document.createElement("div");
    tradeSellStatus.className = "TradeSellStatus";
    tradeSellStatus.id = "TradeSellStatus";

    let spanHead = document.createElement("span");
    spanHead.innerHTML = "Sell status";
    spanHead.className = "header";
    spanHead.style.marginTop = "5px";
    spanHead.style.display = "block";
    spanHead.style.fontSize = "12px";
    tradeSellStatus.appendChild(spanHead);

    let altCashTable = document.createElement("table");
    altCashTable.className = "altCash";
    altCashTable.id = "altCash" + id;


    let tr = document.createElement("tr");

    function createTD (text) {
        let tdName = document.createElement("td");
        tdName.className = "Value";
        tdName.innerHTML = text;
        tr.appendChild(tdName);
    }

    createTD("Name");
    createTD("Balance");
    createTD("BuyPrice");
    createTD("ProfitPrice");
    createTD("GrowProfitPrice");
    createTD("TopAscOrder");

    altCashTable.appendChild(tr);


    tradeSellStatus.appendChild(altCashTable);

    return tradeSellStatus;
}