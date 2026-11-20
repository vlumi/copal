# gondwana: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them; `survey conform` runs a replacement for a program on them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [GLD010](GLD010/) | 64 of 99 | 29 |
| [GLD020](GLD020/) | 92 of 132 | 44 |
| [GLD030](GLD030/) | 145 of 197 | 61 |
| [GLD040](GLD040/) | 78 of 105 | 30 |
| [GLD050](GLD050/) | 70 of 94 | 43 |
| [GLD060](GLD060/) | 72 of 97 | 33 |
| [GLD070](GLD070/) | 70 of 102 | 35 |
| [GLD080](GLD080/) | 14 of 30 | 4 |
| [GLD090](GLD090/) | 23 of 42 | 10 |
| [GLM010](GLM010/) | 32 of 52 | 12 |
| [GLM020](GLM020/) | 45 of 69 | 21 |
| [GLM030](GLM030/) | 53 of 65 | 20 |
| [GLO010](GLO010/) | 29 of 45 | 9 |
| [GLY010](GLY010/) | 22 of 38 | 9 |
| [GLY020](GLY020/) | 47 of 64 | 27 |
| [GLY030](GLY030/) | 31 of 48 | 13 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [OPEN](jobs/OPEN/) | GLO010 RC=0004 |
| [DAILY](jobs/DAILY/) | GLD010 RC=0000, GLD020 RC=0000, GLD030 RC=0004, GLD040 RC=0004, GLD060 RC=0000, GLD070 RC=0004, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [BILLING](jobs/BILLING/) | GLM010 RC=0000 |
| [RESULTS](jobs/RESULTS/) | GLD050 RC=0004 |
| [CLOSE](jobs/CLOSE/) | GLM020 RC=0000, GLM030 RC=0000, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [ANNUAL](jobs/ANNUAL/) | GLY010 RC=0000, GLY020 RC=0000, GLD080 RC=0000 |
| [CERT](jobs/CERT/) | GLY030 RC=0004 |
