// While the bank runs or resets, disable the form's button and say so.
document.addEventListener("submit", (event) => {
  const form = event.target;
  const button = form.querySelector("button");
  if (form.dataset.busy && button) {
    button.disabled = true;
    button.textContent = form.dataset.busy;
  }
});

// On the online page, follow the online channel's events as they come.
const feed = document.querySelector("#feed tbody");
if (feed) {
  const functions = { PURC: "purchase", REVS: "reversal", WDRL: "withdrawal", BALQ: "balance" };
  const money = (amount) =>
    Number(amount).toLocaleString("en-US", { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  const cell = (row, text, className) => {
    const td = row.insertCell();
    td.textContent = text;
    if (className) td.className = className;
    return td;
  };
  const counters = { requests: 0, approved: 0, declined: 0, held: 0 };
  const count = (name, by) => {
    counters[name] += by;
    document.getElementById(name).textContent =
      name === "held" ? money(counters[name] / 100) : counters[name].toLocaleString("en-US");
  };
  const accounts = document.querySelector("#accounts tbody");
  const status = document.getElementById("status");

  const showAccount = (e) => {
    const id = "account-" + e.account;
    document.getElementById(id)?.remove();
    const row = accounts.insertRow(0);
    row.id = id;
    const link = cell(row, "", "num").appendChild(document.createElement("a"));
    link.href = "/accounts/" + e.account;
    link.textContent = e.account;
    cell(row, functions[e.function] + " " + e.time);
    cell(row, money(e.balance), "num");
    cell(row, money(e.available), "num" + (Number(e.available) < 0 ? " debit" : ""));
    while (accounts.rows.length > 25) accounts.deleteRow(-1);
  };

  const showAuthorization = (e) => {
    count("requests", 1);
    count(e.approved ? "approved" : "declined", 1);
    if (e.approved && (e.function === "PURC" || e.function === "WDRL")) {
      count("held", Math.round((Number(e.amount) + Number(e.fee)) * 100));
    }
    if (e.approved && e.function === "REVS") count("held", -Math.round(Number(e.amount) * 100));
    const row = feed.insertRow(0);
    cell(row, e.time);
    cell(row, functions[e.function] || e.function);
    cell(row, "…" + e.card.slice(-4), "num");
    cell(row, e.account || "", "num");
    cell(row, e.terminal);
    cell(row, e.function === "BALQ" ? "" : money(e.approved ? e.amount : e.requested), "num");
    cell(row, Number(e.fee) ? money(e.fee) : "", "num");
    cell(row, e.result + " " + e.message, e.approved ? "ok" : "debit");
    cell(row, e.approved ? money(e.available) : "", "num");
    while (feed.rows.length > 200) feed.deleteRow(-1);
    if (e.approved && e.account) showAccount(e);
  };

  const events = new EventSource("/live/events");
  events.onmessage = (message) => {
    const e = JSON.parse(message.data);
    if (e.kind === "authorization") {
      showAuthorization(e);
    } else {
      status.textContent = e.kind === "opened" ? "Open" : "Closed: the batch runs";
      status.className = e.kind === "opened" ? "ok" : "note";
    }
  };
}
