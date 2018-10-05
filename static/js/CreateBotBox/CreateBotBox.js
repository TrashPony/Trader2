function CreateBotBox(id) {
    let main = document.getElementById("main");

    let botBox = document.createElement("div");
    botBox.className = "bots";
    botBox.id = id;
    main.appendChild(botBox);

    let botIcon = document.createElement("div");
    botIcon.className = "botIcon";
    botIcon.id = "botIcon" + id;
    botBox.appendChild(botIcon);

    let headerTable = CreateHeaderBotTable(id);
    botBox.appendChild(headerTable);

    let cashTable = CreateCashTable(id);
    botBox.appendChild(cashTable);

    let clearDiv = document.createElement("div");
    clearDiv.style.clear = "both";
    botBox.appendChild(clearDiv);

    let tradeStatus = CreateTradeStatusBlock(id);
    botBox.appendChild(tradeStatus);

    let logBlock = CreateLogBlock(id);
    botBox.appendChild(logBlock);
}

function CreateLogBlock(id) {
    let log = document.createElement("div");
    log.className = "log";
    log.id = "log" + id;

    return log
}