function CreateHeaderBotTable(id) {
    let table = document.createElement("table");
    table.className = "headerTable";

    let tr = document.createElement("tr");
    table.appendChild(tr);

    let tdName = document.createElement("td");
    tdName.className = "header";
    tdName.id = "name" + id;
    tr.appendChild(tdName);

    let tdEfficiency = document.createElement("td");
    tdEfficiency.className = "Value";
    tdEfficiency.innerHTML = "Эффективность";
    tr.appendChild(tdEfficiency);

    let tdEfficiencyValue = document.createElement("td");
    tdEfficiencyValue.id = "Efficiency" + id;
    tr.appendChild(tdEfficiencyValue);

    return table;
}