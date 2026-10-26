# gondwana: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them; `survey conform` runs a replacement for a program on them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [GLD010](GLD010/) | 36 of 67 | 21 |
| [GLD020](GLD020/) | 29 of 66 | 12 |
| [GLD030](GLD030/) | 33 of 117 | 20 |
| [GLD040](GLD040/) | 50 of 73 | 30 |
| [GLD050](GLD050/) | 50 of 71 | 33 |
| [GLD060](GLD060/) | 15 of 29 | 5 |
| [GLD090](GLD090/) | 16 of 30 | 8 |
| [GLM010](GLM010/) | 28 of 46 | 13 |
| [GLM020](GLM020/) | 34 of 56 | 13 |
| [GLM030](GLM030/) | 19 of 51 | 9 |
| [GLO010](GLO010/) | 18 of 33 | 10 |
| [GLY010](GLY010/) | 20 of 38 | 10 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [OPEN](jobs/OPEN/) | GLO010 RC=0000 |
| [DAILY](jobs/DAILY/) | GLD010 RC=0000, GLD020 RC=0000, GLD030 RC=0000, GLD040 RC=0000, GLD060 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [BILLING](jobs/BILLING/) | GLM010 RC=0000 |
| [RESULTS](jobs/RESULTS/) | GLD050 RC=0000 |
| [CLOSE](jobs/CLOSE/) | GLM020 RC=0001, GLM030 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [ANNUAL](jobs/ANNUAL/) | GLY010 RC=0000 |
