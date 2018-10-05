function CreateCashTable(id) {
    let table = document.createElement("table");
    table.className = "cashTable";

    function CreateRow(optionText, idValue) {
        let tr = document.createElement("tr");

        let tdText = document.createElement("td");
        tdText.className = "option";
        tdText.innerHTML = optionText;
        tr.appendChild(tdText);

        let tdValue = document.createElement("td");
        tdValue.className = "Value";
        tdValue.id = idValue;
        tr.appendChild(tdValue);

        table.appendChild(tr);
    }

    CreateRow("Start BTC Cash", "StartBTCCash" + id);
    CreateRow("Available BTC Cash", "AvailableBTCCash" + id);
    CreateRow("In Trade Strategy", "InTradeStrategy" + id);
    CreateRow("Out Trade Strategy", "OutTradeStrategy" + id);

    return table;
}