# gondwana: cases

For each main program, input data that reaches the rules of its decision tables, and of the programs it calls, with what the program wrote and returned: the cases a rewrite has to match. `survey oracle` found them, and replays them to check them; `survey conform` runs a replacement for a program on them.

| Program | Rules reached | Cases |
| --- | --: | --: |
| [GLD010](GLD010/) | 39 of 96 | 22 |
| [GLD020](GLD020/) | 73 of 132 | 32 |
| [GLD030](GLD030/) | 87 of 197 | 53 |
| [GLD040](GLD040/) | 65 of 105 | 37 |
| [GLD050](GLD050/) | 60 of 94 | 39 |
| [GLD060](GLD060/) | 46 of 97 | 22 |
| [GLD070](GLD070/) | 62 of 102 | 32 |
| [GLD080](GLD080/) | 14 of 30 | 5 |
| [GLD090](GLD090/) | 23 of 42 | 11 |
| [GLM010](GLM010/) | 32 of 52 | 14 |
| [GLM020](GLM020/) | 45 of 69 | 23 |
| [GLM030](GLM030/) | 44 of 65 | 17 |
| [GLO010](GLO010/) | 24 of 45 | 14 |
| [GLY010](GLY010/) | 20 of 38 | 10 |
| [GLY020](GLY020/) | 40 of 64 | 24 |
| [GLY030](GLY030/) | 30 of 48 | 16 |

And for each job, a batch: the job run through its steps on inputs from its programs' cleanest cases, with the datasets it wrote and its output.

| Job | Steps |
| --- | --- |
| [OPEN](jobs/OPEN/) | GLO010 RC=0000 |
| [DAILY](jobs/DAILY/) | GLD010 RC=0000, GLD020 RC=0000, GLD030 RC=0000, GLD040 RC=0000, GLD060 RC=0000, GLD070 RC=0000, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [BILLING](jobs/BILLING/) | GLM010 RC=0000 |
| [RESULTS](jobs/RESULTS/) | GLD050 RC=0000 |
| [CLOSE](jobs/CLOSE/) | GLM020 RC=0000, GLM030 RC=0000, GLD080 RC=0000, UNLOAD RC=0000, UNLOAD.GLD090 RC=0000 |
| [ANNUAL](jobs/ANNUAL/) | GLY010 RC=0000, GLY020 RC=0001, GLD080 RC=0000 |
| [CERT](jobs/CERT/) | GLY030 RC=0000 |
