# GLD040: cases

30 cases reach 44 of the 57 rules of GLD040's decision tables, and of the programs it calls, 6 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD040

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 09](09/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 11](11/) |
| 3 | [case 12](12/) |
| 4 | [case 13](13/) |
| 5 | [case 14](14/) |
| 6 | [case 15](15/) |
| 7 | [case 01](01/) |

### [MNTREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD040.md#mntreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-MNTREQ-OK` Y |
| 2 | not reached: `WS-MNTREQ-OK` N, `WS-MNTREQ-ｵﾜﾘ` N |
| 3 | [case 08](08/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | [case 03](03/) |
| 4 | [case 04](04/) |
| 5 | [case 05](05/) |
| 6 | [case 06](06/) |
| 7 | not reached: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 8 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 9 | [case 16](16/) |
| 10 | [case 23](23/) |
| 11 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 12 | [case 17](17/) |
| 13 | [case 24](24/) |
| 14 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 15 | [case 18](18/) |
| 16 | [case 25](25/) |
| 17 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 18 | [case 19](19/) |
| 19 | [case 26](26/) |
| 20 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 21 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` Y |
| 22 | [case 20](20/) |
| 23 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` N |
| 24 | not reached: `WS-POLMST-ﾅｼ` N, `MR-ｼﾞｭｳｼｮﾍﾝｺｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `MR-ｺｳｻﾞﾍﾝｺｳ` N, `MR-ｼﾒｲﾍﾝｺｳ` N, `MR-ｳｹﾄﾘﾆﾝﾍﾝｺｳ` N, `MR-ﾌｯｶﾂ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y, `WS-POLMST-OK, again` Y |
| 25 | [case 21](21/) |

### [ﾍﾝｺｳ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾍﾝｺｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 23](23/) |
| 2 | [case 16](16/) |
| 3 | [case 24](24/) |
| 4 | [case 17](17/) |
| 5 | [case 25](25/) |
| 6 | [case 18](18/) |
| 7 | [case 26](26/) |
| 8 | [case 19](19/) |
| 9 | [case 20](20/) |
| 10 | [case 21](21/) |

### [ﾌｯｶﾂ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ﾌｯｶﾂ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 20](20/) |
| 2 | [case 27](27/) |
| 3 | [case 29](29/) |
| 4 | not reached: `PM-ｼｯｺｳ` Y, `PM-ﾈﾝﾊﾞﾗｲ` Y, `WS-ｹｲｶﾆｯｽｳ > 365` N, `MR-ﾆｭｳｷﾝｶﾞｸ < WS-ﾐﾉｳｶﾞｸ` N |
| 5 | [case 22](22/) |
| 6 | [case 28](28/) |
| 7 | [case 30](30/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD040.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `LK-ﾜﾚｷ-ﾍ` N, `LK-ｾｲﾚｷ-ﾍ` Y |
| 3 | not reached: `LK-ﾜﾚｷ-ﾍ` N, `LK-ｾｲﾚｷ-ﾍ` N |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 07](07/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ` Y, `WS-ｹﾞﾝｺﾞｳ-ｱﾘ` Y |
| 2 | [case 22](22/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ﾈﾝ = 0` Y, `LK-ﾍﾝｶﾝｽﾞﾐ` N |
| 2 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` Y |
| 3 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` Y |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` N |
| 6 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` Y, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` Y |
| 7 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` Y, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` N |
