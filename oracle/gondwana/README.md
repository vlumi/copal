# gondwana: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them; `survey conform` runs a replacement for a program on them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [GLD010](GLD010/) | 53 of 99 | 26 |
| [GLD020](GLD020/) | 80 of 132 | 35 |
| [GLD030](GLD030/) | 124 of 197 | 46 |
| [GLD040](GLD040/) | 70 of 105 | 25 |
| [GLD050](GLD050/) | 67 of 94 | 40 |
| [GLD060](GLD060/) | 48 of 97 | 22 |
| [GLD070](GLD070/) | 66 of 102 | 29 |
| [GLD080](GLD080/) | 14 of 30 | 4 |
| [GLD090](GLD090/) | 23 of 42 | 10 |
| [GLM010](GLM010/) | 32 of 52 | 12 |
| [GLM020](GLM020/) | 45 of 69 | 21 |
| [GLM030](GLM030/) | 50 of 65 | 18 |
| [GLO010](GLO010/) | 29 of 45 | 9 |
| [GLY010](GLY010/) | 22 of 38 | 9 |
| [GLY020](GLY020/) | 41 of 64 | 22 |
| [GLY030](GLY030/) | 31 of 48 | 13 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [OPEN](jobs/OPEN/) | GLO010 RC=0004 |
| [DAILY](jobs/DAILY/) | GLD010 RC=0004, GLD020 RC=0000, GLD030 RC=0001, GLD040 RC=0004, GLD060 RC=0000, GLD070 RC=0000, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [BILLING](jobs/BILLING/) | GLM010 RC=0000 |
| [RESULTS](jobs/RESULTS/) | GLD050 RC=0000 |
| [CLOSE](jobs/CLOSE/) | GLM020 RC=0000, GLM030 RC=0000, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [ANNUAL](jobs/ANNUAL/) | GLY010 RC=0000, GLY020 RC=0001, GLD080 RC=0000 |
| [CERT](jobs/CERT/) | GLY030 RC=0000 |
