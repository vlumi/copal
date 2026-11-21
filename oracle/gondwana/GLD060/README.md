# GLD060: cases

35 cases reach 29 of the 39 rules of GLD060's decision tables, and of the programs it calls, 13 of the 16 of GLC020, 4 of the 5 of GLC040, 21 of the 23 of GLC030, 7 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## GLD060

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 07](07/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 09](09/) |
| 3 | [case 10](10/) |
| 4 | [case 11](11/) |
| 5 | [case 01](01/) |

### [POLMST-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD060.md#polmst-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NOT WS-POLMST-OK AND NOT WS-POLMST-ｵﾜﾘ` Y |
| 2 | [case 01](01/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 29](29/) |
| 2 | [case 03](03/) |
| 3 | [case 01](01/) |
| 4 | [case 01](01/) |

### [ｷｹﾞﾝ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｷｹﾞﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｷｹﾞﾝ-MM = 12` Y |
| 2 | [case 01](01/) |

### [ﾌﾘｶｴ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ﾌﾘｶｴ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 21](21/) |
| 2 | [case 30](30/) |
| 3 | [case 12](12/) |
| 4 | [case 18](18/) |
| 5 | [case 29](29/) |
| 6 | [case 13](13/) |
| 7 | [case 03](03/) |

### [ﾌﾘｶｴ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ﾌﾘｶｴ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 2 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` Y |
| 3 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` N |
| 4 | not reached: `WS-ﾘｿｸ > 0` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |
| 5 | [case 30](30/) |
| 6 | [case 31](31/) |
| 7 | [case 35](35/) |
| 8 | [case 29](29/) |

### [POLMST-ｶｷﾓﾄﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD060.md#polmst-ｶｷﾓﾄﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-OK` N |
| 2 | [case 03](03/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD060.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | [case 29](29/) |

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
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | not reached: `LK-ﾜﾚｷ-ﾍ` N, `LK-ｾｲﾚｷ-ﾍ` N |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 05](05/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 03](03/) |
| 3 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` Y |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | [case 20](20/) |
| 6 | [case 06](06/) |
| 7 | [case 01](01/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | [case 29](29/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 19](19/) |
| 2 | [case 03](03/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 18](18/) |
| 3 | [case 29](29/) |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 18](18/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 27](27/) |
| 2 | [case 15](15/) |
| 3 | [case 17](17/) |
| 4 | [case 03](03/) |
| 5 | [case 14](14/) |
| 6 | [case 34](34/) |
| 7 | [case 26](26/) |
| 8 | [case 28](28/) |
| 9 | [case 16](16/) |
| 10 | [case 25](25/) |
| 11 | [case 33](33/) |
| 12 | [case 23](23/) |
| 13 | [case 24](24/) |
| 14 | [case 13](13/) |
| 15 | [case 22](22/) |

## GLC050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 29](29/) |
| 2 | [case 29](29/) |
| 3 | [case 01](01/) |

### [ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y, `WS-PRMHST-OK, again` N |
| 2 | [case 29](29/) |
| 3 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N, `WS-PRMHST-OK, again` N |
| 4 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` Y, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 5 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` Y |
| 6 | [case 29](29/) |
| 7 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 8 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |
| 9 | not reached: `WS-ｵｰﾌﾟﾝﾏｴ` N, `WS-PRMHST-OK` Y, `WS-ｵｰﾌﾟﾝｽﾞﾐ` N, `WS-ｵｰﾌﾟﾝｽﾞﾐ AND (WS-PRMHST-OK OR WS-PRMHST-ｵﾜﾘ)` N |

### [ﾚﾝﾊﾞﾝ-RTN](../../../blueprints/gondwana/rules/GLC050.md#ﾚﾝﾊﾞﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 29](29/) |
| 2 | [case 32](32/) |
