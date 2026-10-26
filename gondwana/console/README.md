# The console

The life insurer's console: a web window onto the running Gondwana Life Insurance for its operators. It opens the insurer, runs its business days, and shows each day's batch files with the files they read and wrote. Like the rest of the insurer, it knows nothing of the survey.

`make gondwana-console` serves it at <http://127.0.0.1:2525/> (`GONDWANA_ADDR` to change it):

- the insurer today: the next business day it runs, its policies by state with the premiums of those in force, and every business day it has run, newest first, with the batch files it ran, the new policies, benefits, changes and rejects of each, the debits billed, collected and unpaid on their days, the portfolio after it, and its worst errorlevel;
- policies, searched by policy number, or the start of one, or by part of the policyholder's name in kanji or kana, and for each policy its master record field by field under the copybook's katakana names;
- each business day: every batch file it ran, the unload the daily run and the close call among them, step by step with each step's errorlevel; then every file the day read from `IN` and wrote to `OUT` and `RPT`, a small one shown in place and a large one linked, the banks' fixed-length files a record a line.

Its labels are English, as the Go around the insurer is, and the data is shown as the insurer keeps it, in Japanese: everything is read from Shift_JIS, the names, the reports and the files alike, and nothing is translated. The master record's dates stay in their era, as the master keeps them; the lists give them in the Western calendar.

It reads the policies through the layouts computed from the copybooks, the same way the insurer's generator writes them, and reads back each day from the batch files' logs and the files the day left, as `make gondwana-days` reports it. Its buttons open an insurer and run business days by calling the same code as `make gondwana-open` and `make gondwana-days`; the console itself never writes the insurer's files.
