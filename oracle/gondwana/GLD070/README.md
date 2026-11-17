# GLD070: cases

31 cases reach 35 of the 58 rules of GLD070's decision tables, and of the programs it calls, 8 of the 16 of GLC020, 4 of the 5 of GLC040, 14 of the 23 of GLC030. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD070

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 07](07/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 09](09/) |
| 3 | [case 10](10/) |
| 4 | [case 11](11/) |
| 5 | [case 12](12/) |
| 6 | [case 13](13/) |
| 7 | [case 14](14/) |
| 8 | [case 01](01/) |

### [LONREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD070.md#lonreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | not reached: `WS-LONREQ-OK` N, `WS-LONREQ-ｵﾜﾘ` N |
| 3 | [case 06](06/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | [case 03](03/) |
| 4 | not reached: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 5 | not reached: `WS-POLMST-ﾅｼ` N, `LR-ｶｼﾂｹ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y |
| 6 | [case 21](21/) |
| 7 | [case 19](19/) |
| 8 | [case 18](18/) |
| 9 | not reached: `WS-POLMST-ﾅｼ` N, `LR-ｶｼﾂｹ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `LR-ﾍﾝｻｲ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y |
| 10 | [case 16](16/) |

### [ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | [case 18](18/) |
| 3 | [case 21](21/) |
| 4 | [case 20](20/) |
| 5 | [case 05](05/) |

### [ｹﾞﾝﾄﾞ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｹﾞﾝﾄﾞ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 22](22/) |
| 3 | not reached: `LR-ｶｼﾂｹ` Y, `WS-CSVTB-OK` Y, `WS-ｻﾞﾝﾀﾞｶ + LR-ｷﾝｶﾞｸ > WS-ｹﾞﾝﾄﾞｶﾞｸ` N |
| 4 | [case 15](15/) |
| 5 | [case 19](19/) |

### [ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 2 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` Y |
| 3 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 4 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` Y |
| 5 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 6 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` Y |
| 7 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 8 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` Y |
| 9 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 10 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` Y |
| 11 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 12 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` Y |
| 13 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 14 | [case 15](15/) |
| 15 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 16 | [case 19](19/) |

### [ｶｼﾂｹ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｶｼﾂｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PAYOUT-OK` N |
| 2 | not reached: `WS-PAYOUT-OK` Y |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | [case 15](15/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 18](18/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｭｳﾘｮｳ-rtn)

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
| 1 | [case 04](04/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 04](04/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 19](19/) |
| 2 | not reached: `ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ` N, `WS-ｹﾞﾝｺﾞｳ-ｱﾘ` N |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ﾈﾝ = 0` Y, `LK-ﾍﾝｶﾝｽﾞﾐ` N |
| 2 | [case 27](27/) |
| 3 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` Y |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` N |
| 6 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` Y, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` Y |
| 7 | [case 19](19/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 19](19/) |
| 2 | [case 27](27/) |
| 3 | [case 05](05/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾆｯｽｳ > 0` Y |
| 2 | [case 19](19/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 22](22/) |
| 2 | [case 17](17/) |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 22](22/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 17](17/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 2 | [case 25](25/) |
| 3 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 4 | [case 22](22/) |
| 5 | [case 26](26/) |
| 6 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 7 | [case 28](28/) |
| 8 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 9 | [case 23](23/) |
| 10 | [case 29](29/) |
| 11 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 12 | [case 30](30/) |
| 13 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 14 | [case 24](24/) |
| 15 | [case 31](31/) |
