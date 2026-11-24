# GLD070: cases

37 cases reach 42 of the 58 rules of GLD070's decision tables, and of the programs it calls, 11 of the 16 of GLC020, 5 of the 5 of GLC040, 21 of the 23 of GLC030. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## GLD070

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 05](05/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 07](07/) |
| 3 | [case 08](08/) |
| 4 | [case 09](09/) |
| 5 | [case 10](10/) |
| 6 | [case 11](11/) |
| 7 | [case 12](12/) |
| 8 | [case 01](01/) |

### [LONREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD070.md#lonreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached, a file found to fail once open: `WS-LONREQ-OK` N, `WS-LONREQ-ｵﾜﾘ` N |
| 3 | [case 01](01/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 02](02/) |
| 3 | [case 01](01/) |
| 4 | not reached, a file found to fail once open: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 5 | [case 34](34/) |
| 6 | [case 04](04/) |
| 7 | [case 16](16/) |
| 8 | [case 04](04/) |
| 9 | not reached: `WS-POLMST-ﾅｼ` N, `LR-ｶｼﾂｹ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `LR-ﾍﾝｻｲ` N, `WS-POLMST-OK` Y, `WS-ﾘﾕｳ-ﾅｼ, again` Y |
| 10 | [case 13](13/) |

### [ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 13](13/) |
| 2 | [case 14](14/) |
| 3 | [case 04](04/) |
| 4 | [case 04](04/) |
| 5 | [case 15](15/) |

### [ｹﾞﾝﾄﾞ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｹﾞﾝﾄﾞ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 23](23/) |
| 2 | [case 15](15/) |
| 3 | [case 34](34/) |
| 4 | [case 16](16/) |
| 5 | [case 16](16/) |

### [ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 2 | not reached: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` Y |
| 3 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 4 | [case 37](37/) |
| 5 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 6 | [case 31](31/) |
| 7 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` Y, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 8 | [case 22](22/) |
| 9 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 10 | not reached: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` Y |
| 11 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 12 | [case 34](34/) |
| 13 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` Y, `WS-POLMST-OK` N |
| 14 | [case 16](16/) |
| 15 | not reached, a file found to fail once open: `WS-ﾘｿｸ > 0` N, `LR-ｶｼﾂｹ` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ = 0` N, `WS-POLMST-OK` N |
| 16 | [case 34](34/) |

### [ｶｼﾂｹ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｶｼﾂｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-PAYOUT-OK` N |
| 2 | [case 34](34/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLJRNL-OK` N |
| 2 | [case 16](16/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD070.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 15](15/) |
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
| 1 | [case 15](15/) |
| 2 | not reached: `ｹﾞﾝｺﾞｳ-ｷｺﾞｳ(ｹﾞﾝｺﾞｳ-IX) = LK-ｹﾞﾝｺﾞｳ` N, `WS-ｹﾞﾝｺﾞｳ-ｱﾘ` N |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ﾈﾝ = 0` Y, `LK-ﾍﾝｶﾝｽﾞﾐ` N |
| 2 | [case 34](34/) |
| 3 | [case 36](36/) |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` N |
| 6 | [case 35](35/) |
| 7 | [case 15](15/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 34](34/) |
| 3 | [case 16](16/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 20](20/) |
| 2 | [case 15](15/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 23](23/) |
| 3 | [case 15](15/) |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | not reached, a file found to fail once open: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 23](23/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 29](29/) |
| 2 | [case 19](19/) |
| 3 | [case 34](34/) |
| 4 | [case 15](15/) |
| 5 | [case 18](18/) |
| 6 | [case 33](33/) |
| 7 | [case 28](28/) |
| 8 | [case 30](30/) |
| 9 | [case 21](21/) |
| 10 | [case 27](27/) |
| 11 | [case 32](32/) |
| 12 | [case 25](25/) |
| 13 | [case 26](26/) |
| 14 | [case 17](17/) |
| 15 | [case 24](24/) |
