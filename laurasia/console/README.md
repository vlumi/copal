# The console

The card processor's console: a web window onto the running Laurasia Card Services for its operators. It opens the processor, runs its days, and shows each day's night and billing cycle with the files they read and wrote. Like the rest of the processor, it knows nothing of the survey.

`make laurasia-console` serves it at <http://127.0.0.1:4649/> (`LAURASIA_ADDR` to change it):

- the processor today: the next day it runs, its cards in sum, by status, with their balances and limits, and every day it has run, newest first, with the presentments, rejects, postings, payments and statements of each and its worst return code; a weekend shows that no clearing file came;
- cards, searched by card number, or the start of one, or by account number, and for each card its account record field by field and every posting to it since the processor opened;
- each day: the night, on a business day, and the billing cycle, step by step with each step's return code; then every file the day read from `inbound/` and wrote to `outbound/` and `reports/`, a small one shown in place and a large one linked.

It reads the cards and postings through the layouts computed from the copybooks, the same way the processor's generator writes them, and reads back each day from the scripts' logs and the files the day left, as `make laurasia-days` reports it. Its buttons open a processor and run days by calling the same code as `make laurasia-open` and `make laurasia-days`; the console itself never writes the processor's files.

The processor has no online side of its own, so the console has no online page: the networks' switch authorizes around the clock, and the night reads only its log.
