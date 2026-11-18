# GLD050: cases

39 cases reach 52 of the 64 rules of GLD050's decision tables, and of the programs it calls, 5 of the 16 of GLC020, 3 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 05](05/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 08](08/) |
| 3 | [case 09](09/) |
| 4 | [case 10](10/) |
| 5 | [case 11](11/) |
| 6 | [case 12](12/) |
| 7 | [case 01](01/) |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 07](07/) |
| 3 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` Y |
| 4 | [case 01](01/) |
| 5 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` Y |
| 6 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` Y |
| 7 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` N, `WS-ｴﾝﾄﾞ-ﾏﾁ AND ZG-ｴﾝﾄﾞ` Y |
| 8 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` N, `WS-ｴﾝﾄﾞ-ﾏﾁ AND ZG-ｴﾝﾄﾞ` N |

### [ﾀﾒﾙ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾀﾒﾙ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 17](17/) |

### [ｹﾝｻﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｹﾝｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 14](14/) |
| 3 | [case 15](15/) |
| 4 | [case 02](02/) |

### [ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | not reached: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 3 | [case 25](25/) |
| 4 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` Y, `WS-POLMST-OK, again` N |
| 5 | [case 16](16/) |
| 6 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` N, `WS-POLMST-OK, again` N |
| 7 | [case 17](17/) |
| 8 | [case 23](23/) |

### [ｷﾝｶﾞｸ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｷﾝｶﾞｸ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | [case 23](23/) |

### [ﾌﾘｶｴｽﾞﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴｽﾞﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 35](35/) |
| 2 | [case 26](26/) |
| 3 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` N, `PM-ﾕｳﾖﾁｭｳ` Y |
| 4 | [case 38](38/) |
| 5 | [case 36](36/) |
| 6 | [case 27](27/) |
| 7 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` N, `WS-ﾊﾗｲｺﾐ-MM > 12` N, `PM-ﾕｳﾖﾁｭｳ` Y, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y |
| 8 | [case 39](39/) |
| 9 | [case 24](24/) |
| 10 | [case 16](16/) |
| 11 | [case 37](37/) |
| 12 | [case 28](28/) |

### [ﾌﾘｶｴﾌﾉｳ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴﾌﾉｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 18](18/) |
| 3 | [case 19](19/) |
| 4 | [case 20](20/) |
| 5 | [case 21](21/) |
| 6 | [case 22](22/) |
| 7 | [case 29](29/) |
| 8 | [case 30](30/) |
| 9 | [case 31](31/) |
| 10 | [case 32](32/) |
| 11 | [case 33](33/) |
| 12 | [case 34](34/) |

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
| 1 | not reached: `LK-ｶｷｺﾐ` Y |
| 2 | not reached: `LK-ｶｷｺﾐ` N, `LK-ｸﾛｰｽﾞ AND WS-ｵｰﾌﾟﾝｽﾞﾐ` Y |
| 3 | [case 01](01/) |

### [ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y, `WS-PRMHST-OK, again` N |
| 2 | [case 16](16/) |
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
| 1 | not reached: `WS-PRMHST-OK` Y |
| 2 | [case 16](16/) |
