# Pangaea's own documentation

<img src="../brand/seal.png" alt="The seal of the First Bank of Pangaea" width="160" align="right">

What follows is written in the voice of the First Bank of Pangaea's operations team, as the bank would have written it for itself: how the day runs, which jobs run when, what the datasets are, and what to do when something fails. It is the kind of documentation a modernization usually inherits: operational, partial, written for people who already know the system, and silent on how the programs work inside. That part is what the [survey](../../survey/) is for; its reading of the bank is in [blueprints/pangaea/](../../blueprints/pangaea/README.md).

Like the bank, these documents are fiction written for this project. The procedures in them are real in the sense that matters: `make bank-days` and the console do what the run book says.

- [operations.md](operations.md): the business day, job logs, return codes, and restarting a failed day.
- [schedule.md](schedule.md): the jobs and their steps, in the order they run.
- [datasets.md](datasets.md): every dataset the bank keeps, who writes it and who reads it.
- [online.md](online.md): the card and ATM authorization service.
