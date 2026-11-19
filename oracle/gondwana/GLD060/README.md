# GLD060: cases

22 cases reach 21 of the 39 rules of GLD060's decision tables, and of the programs it calls, 7 of the 16 of GLC020, 3 of the 5 of GLC040, 14 of the 23 of GLC030, 1 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD060

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 04](04/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 06](06/) |
| 3 | [case 07](07/) |
| 4 | [case 08](08/) |
| 5 | [case 01](01/) |

### [POLMST-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD060.md#polmst-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NOT WS-POLMST-OK AND NOT WS-POLMST-ｵﾜﾘ` Y |
| 2 | [case 01](01/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾕｳﾖﾁｭｳ` Y, `WS-ｷﾞｮｳﾑ-YYYYMM > WS-ｷｹﾞﾝ-YYYYMM` Y, `WS-ﾌﾘｶｴ-ｶ` Y |
| 2 | [case 10](10/) |
| 3 | [case 01](01/) |
| 4 | [case 02](02/) |

### [ｷｹﾞﾝ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｷｹﾞﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｷｹﾞﾝ-MM = 12` Y |
| 2 | [case 01](01/) |

### [ﾌﾘｶｴ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ﾌﾘｶｴ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | not reached: `PM-ﾂｷﾊﾞﾗｲ` N, `WS-CSVTB-OK` Y, `WS-ﾌﾘｶｴ-ｶｲｽｳ > 0` Y, `WS-ｻﾞﾝﾀﾞｶ > WS-ﾍﾝﾚｲｷﾝ` N |
| 3 | [case 12](12/) |
| 4 | [case 11](11/) |
| 5 | not reached: `PM-ﾂｷﾊﾞﾗｲ` Y, `WS-CSVTB-OK` Y, `WS-ﾌﾘｶｴ-ｶｲｽｳ > 0` Y, `WS-ｻﾞﾝﾀﾞｶ > WS-ﾍﾝﾚｲｷﾝ` N |
| 6 | [case 10](10/) |
| 7 | [case 13](13/) |

### [ﾌﾘｶｴ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ﾌﾘｶｴ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 2 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` Y |
| 3 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` N |
| 4 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |
| 5 | not reached: `WS-ﾘｿｸ > 0` N, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 6 | not reached: `WS-ﾘｿｸ > 0` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` Y |
| 7 | not reached: `WS-ﾘｿｸ > 0` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` N |
| 8 | not reached: `WS-ﾘｿｸ > 0` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |

### [POLMST-ｶｷﾓﾄﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD060.md#polmst-ｶｷﾓﾄﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-OK` N |
| 2 | [case 10](10/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | not reached: `WS-GLJRNL-OK` Y |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRMHST-OK` N |
| 2 | [case 01](01/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
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
| 1 | [case 01](01/) |
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
| 7 | [case 01](01/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | [case 09](09/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾆｯｽｳ > 0` Y |
| 2 | [case 10](10/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 11](11/) |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 11](11/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 2 | [case 16](16/) |
| 3 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 4 | [case 10](10/) |
| 5 | [case 17](17/) |
| 6 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 7 | [case 19](19/) |
| 8 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 9 | [case 14](14/) |
| 10 | [case 20](20/) |
| 11 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 12 | [case 21](21/) |
| 13 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 14 | [case 15](15/) |
| 15 | [case 22](22/) |

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
| 2 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y, `WS-PRMHST-OK, again` Y |
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
| 2 | not reached: `WS-PRMHST-OK` N |
