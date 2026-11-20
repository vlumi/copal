# GLD050: cases

43 cases reach 59 of the 64 rules of GLD050's decision tables, and of the programs it calls, 5 of the 16 of GLC020, 6 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 06](06/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 07](07/) |
| 2 | [case 09](09/) |
| 3 | [case 10](10/) |
| 4 | [case 11](11/) |
| 5 | [case 12](12/) |
| 6 | [case 13](13/) |
| 7 | [case 01](01/) |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 08](08/) |
| 3 | [case 01](01/) |
| 4 | [case 03](03/) |
| 5 | [case 01](01/) |
| 6 | [case 01](01/) |
| 7 | [case 01](01/) |
| 8 | [case 02](02/) |

### [ﾀﾒﾙ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾀﾒﾙ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 18](18/) |

### [ｹﾝｻﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｹﾝｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 15](15/) |
| 3 | [case 16](16/) |
| 4 | [case 01](01/) |

### [ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 3 | [case 26](26/) |
| 4 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` Y, `WS-POLMST-OK, again` N |
| 5 | [case 17](17/) |
| 6 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` N, `WS-POLMST-OK, again` N |
| 7 | [case 18](18/) |
| 8 | [case 24](24/) |

### [ｷﾝｶﾞｸ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｷﾝｶﾞｸ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 24](24/) |

### [ﾌﾘｶｴｽﾞﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴｽﾞﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 42](42/) |
| 2 | [case 40](40/) |
| 3 | [case 37](37/) |
| 4 | [case 27](27/) |
| 5 | [case 38](38/) |
| 6 | [case 28](28/) |
| 7 | [case 43](43/) |
| 8 | [case 41](41/) |
| 9 | [case 39](39/) |
| 10 | [case 29](29/) |
| 11 | [case 25](25/) |
| 12 | [case 17](17/) |

### [ﾌﾘｶｴﾌﾉｳ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴﾌﾉｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 19](19/) |
| 3 | [case 20](20/) |
| 4 | [case 21](21/) |
| 5 | [case 22](22/) |
| 6 | [case 23](23/) |
| 7 | [case 31](31/) |
| 8 | [case 32](32/) |
| 9 | [case 33](33/) |
| 10 | [case 34](34/) |
| 11 | [case 35](35/) |
| 12 | [case 36](36/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRMHST-OK` N |
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
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ` Y, `WS-ｹﾞﾝｺﾞｳ-ｱﾘ` Y |
| 2 | not reached: `ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ` N, `WS-ｹﾞﾝｺﾞｳ-ｱﾘ` N |

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

## GLC050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 17](17/) |
| 3 | [case 01](01/) |

### [ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y, `WS-PRMHST-OK, again` N |
| 2 | [case 17](17/) |
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
| 1 | [case 17](17/) |
| 2 | [case 30](30/) |
