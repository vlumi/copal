# gondwana: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them; `survey conform` runs a replacement for a program on them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [GLD010](GLD010/) | 36 of 68 | 21 |
| [GLD020](GLD020/) | 44 of 107 | 19 |
| [GLD030](GLD030/) | 31 of 172 | 20 |
| [GLD040](GLD040/) | 52 of 76 | 32 |
| [GLD050](GLD050/) | 56 of 79 | 39 |
| [GLD060](GLD060/) | 39 of 81 | 20 |
| [GLD070](GLD070/) | 61 of 102 | 31 |
| [GLD090](GLD090/) | 19 of 36 | 9 |
| [GLM010](GLM010/) | 31 of 50 | 14 |
| [GLM020](GLM020/) | 34 of 56 | 13 |
| [GLM030](GLM030/) | 19 of 51 | 9 |
| [GLO010](GLO010/) | 18 of 33 | 10 |
| [GLY010](GLY010/) | 20 of 38 | 10 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [OPEN](jobs/OPEN/) | GLO010 RC=0000 |
| [DAILY](jobs/DAILY/) | GLD010 RC=0000, GLD020 RC=0000, GLD030 RC=0000, GLD040 RC=0000, GLD060 RC=0000, GLD070 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [BILLING](jobs/BILLING/) | GLM010 RC=0000 |
| [RESULTS](jobs/RESULTS/) | GLD050 RC=0000 |
| [CLOSE](jobs/CLOSE/) | GLM020 RC=0001, GLM030 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [ANNUAL](jobs/ANNUAL/) | GLY010 RC=0000 |
