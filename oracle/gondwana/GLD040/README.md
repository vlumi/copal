# GLD040: cases

32 cases reach 61 of the 75 rules of GLD040's decision tables, and of the programs it calls, 13 of the 16 of GLC020, 6 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## GLD040

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 04](04/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 06](06/) |
| 3 | [case 07](07/) |
| 4 | [case 08](08/) |
| 5 | [case 09](09/) |
| 6 | [case 10](10/) |
| 7 | [case 11](11/) |
| 8 | [case 01](01/) |

### [MNTREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD040.md#mntreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached, a file found to fail once open: `WS-MNTREQ-OK` N, `WS-MNTREQ-ｵﾜﾘ` N |
| 3 | [case 01](01/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 01](01/) |
| 4 | [case 01](01/) |
| 5 | [case 01](01/) |
| 6 | [case 01](01/) |
| 7 | [case 02](02/) |
| 8 | not reached, a file found to fail once open: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 9 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 10 | [case 15](15/) |
| 11 | [case 14](14/) |
| 12 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 13 | [case 14](14/) |
| 14 | [case 15](15/) |
| 15 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 16 | [case 14](14/) |
| 17 | [case 16](16/) |
| 18 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 19 | [case 14](14/) |
| 20 | [case 17](17/) |
| 21 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 22 | [case 18](18/) |
| 23 | [case 14](14/) |
| 24 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` N, `MR-ﾄｸﾔｸｶｲﾔｸ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 25 | [case 19](19/) |
| 26 | [case 13](13/) |
| 27 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` N, `MR-ﾄｸﾔｸｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 28 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` N, `MR-ﾄｸﾔｸｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` Y |
| 29 | [case 12](12/) |

### [ﾍﾝｺｳ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾍﾝｺｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 15](15/) |
| 3 | [case 15](15/) |
| 4 | [case 14](14/) |
| 5 | [case 16](16/) |
| 6 | [case 14](14/) |
| 7 | [case 17](17/) |
| 8 | [case 14](14/) |
| 9 | [case 14](14/) |
| 10 | [case 13](13/) |
| 11 | [case 12](12/) |

### [ﾌｯｶﾂ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾌｯｶﾂ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 29](29/) |
| 3 | [case 28](28/) |
| 4 | [case 20](20/) |
| 5 | [case 31](31/) |
| 6 | [case 30](30/) |
| 7 | [case 21](21/) |
| 8 | [case 23](23/) |
| 9 | [case 22](22/) |
| 10 | [case 18](18/) |

### [ﾄｸﾔｸｶｲﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾄｸﾔｸｶｲﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 22](22/) |
| 2 | [case 13](13/) |
| 3 | not reached, a file found to fail once open: `PM-ﾕｳｺｳ` Y, `WS-COVMST-ﾅｼ` N, `WS-COVMST-OK` N |
| 4 | [case 13](13/) |
| 5 | [case 19](19/) |

### [ﾄｸﾔｸｶｲﾔｸ-ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾄｸﾔｸｶｲﾔｸ-ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-COVMST-OK` N |
| 2 | [case 19](19/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRMHST-OK` N |
| 2 | [case 01](01/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 18](18/) |
| 3 | not reached: `LK-ﾜﾚｷ-ﾍ` N, `LK-ｾｲﾚｷ-ﾍ` N |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 25](25/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 24](24/) |
| 2 | [case 23](23/) |
| 3 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` Y |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | [case 32](32/) |
| 6 | [case 26](26/) |
| 7 | [case 18](18/) |

## GLC050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 18](18/) |
| 3 | [case 01](01/) |

### [ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y, `WS-PRMHST-OK, again` N |
| 2 | [case 18](18/) |
| 3 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N, `WS-PRMHST-OK, again` N |
| 4 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 5 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y |
| 6 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y |
| 7 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 8 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 9 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |

### [ﾚﾝﾊﾞﾝ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ﾚﾝﾊﾞﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 27](27/) |
