# GLD020: cases

49 cases reach 58 of the 80 rules of GLD020's decision tables, and of the programs it calls, 15 of the 16 of GLC020, 4 of the 5 of GLC040, 20 of the 23 of GLC030, 5 of the 8 of GLC010. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

## GLD020

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 08](08/) |
| 3 | [case 02](02/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 12](12/) |
| 3 | [case 09](09/) |
| 4 | [case 10](10/) |
| 5 | [case 13](13/) |
| 6 | [case 14](14/) |
| 7 | [case 15](15/) |
| 8 | [case 16](16/) |
| 9 | [case 02](02/) |

### [POLMST-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#polmst-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `NOT WS-POLMST-OK AND NOT WS-POLMST-ｵﾜﾘ` Y |
| 2 | [case 02](02/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `PM-ﾎｼｮｳﾁｭｳ` Y, `WS-ｵｳﾄｳ-ｱﾘ` Y, `WS-POLMST-OK` N |
| 2 | [case 06](06/) |
| 3 | [case 02](02/) |
| 4 | [case 02](02/) |

### [ｵｳﾄｳﾋﾞ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｳﾄｳﾋﾞ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 20](20/) |
| 2 | [case 03](03/) |
| 3 | [case 06](06/) |
| 4 | [case 02](02/) |

### [ｵｳﾄｳ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｳﾄｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` Y |
| 3 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾏﾝｷ` N |
| 4 | [case 29](29/) |
| 5 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` N |
| 6 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` Y |
| 7 | [case 17](17/) |
| 8 | [case 43](43/) |

### [ｼｭｹｲﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｹｲﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 17](17/) |
| 3 | [case 35](35/) |
| 4 | not reached: `WS-COVMST-OK` Y, `PM-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y, `HS-ﾕｳｺｳ` Y, `PM-ﾖｳﾛｳ` Y |
| 5 | [case 43](43/) |
| 6 | [case 06](06/) |
| 7 | [case 29](29/) |

### [ｼｭｹｲﾔｸ-ｺｳｼﾝ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｹｲﾔｸ-ｺｳｼﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 29](29/) |

### [ﾄｸﾔｸ-ｿｳｻ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-ｿｳｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 2 | [case 25](25/) |
| 3 | not reached, a file found to fail once open: `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 4 | [case 06](06/) |

### [ﾄｸﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` Y, `HS-ﾕｳｺｳ, again` Y |
| 2 | [case 45](45/) |
| 3 | [case 33](33/) |
| 4 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` N, `HS-ﾕｳｺｳ, again` N, `HS-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y |
| 5 | [case 25](25/) |
| 6 | [case 35](35/) |

### [ﾄｸﾔｸ-ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 33](33/) |
| 2 | not reached: `HS-ﾕｳｺｳ` N |

### [ﾎｼｮｳ-ｶｷﾓﾄﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾎｼｮｳ-ｶｷﾓﾄﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-COVMST-OK` N |
| 2 | [case 06](06/) |

### [ﾘｿｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾘｿｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 06](06/) |

### [ﾁｮｳｶ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾁｮｳｶ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 26](26/) |
| 2 | [case 06](06/) |
| 3 | not reached: `WS-CSVTB-OK` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > WS-ﾍﾝﾚｲｷﾝ` N |

### [ﾏﾝｷ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾏﾝｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 45](45/) |
| 2 | [case 43](43/) |
| 3 | [case 44](44/) |
| 4 | [case 49](49/) |

### [ｺｳｼﾝ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｺｳｼﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 29](29/) |
| 2 | [case 27](27/) |
| 3 | [case 28](28/) |
| 4 | [case 06](06/) |
| 5 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 6 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 7 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y |
| 8 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLJRNL-OK` N |
| 2 | [case 18](18/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLSCTL-OK` N |
| 2 | [case 02](02/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-PRODTB-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-RATETB-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN, in its loop](../../../blueprints/gondwana/rules/GLD020.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 02](02/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |
| 3 | not reached: `LK-ﾜﾚｷ-ﾍ` N, `LK-ｾｲﾚｷ-ﾍ` N |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 01](01/) |

### [ｾｲﾚｷ-ﾍﾝｶﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ｾｲﾚｷ-ﾍﾝｶﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 05](05/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 29](29/) |
| 3 | [case 07](07/) |
| 4 | [case 31](31/) |
| 5 | [case 34](34/) |
| 6 | [case 06](06/) |
| 7 | [case 02](02/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | [case 21](21/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 06](06/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | [case 26](26/) |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | not reached, a file found to fail once open: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 26](26/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 06](06/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 42](42/) |
| 2 | [case 23](23/) |
| 3 | [case 24](24/) |
| 4 | [case 06](06/) |
| 5 | [case 19](19/) |
| 6 | [case 48](48/) |
| 7 | [case 40](40/) |
| 8 | [case 41](41/) |
| 9 | [case 22](22/) |
| 10 | [case 39](39/) |
| 11 | [case 47](47/) |
| 12 | [case 37](37/) |
| 13 | [case 38](38/) |
| 14 | [case 18](18/) |
| 15 | [case 36](36/) |

## GLC010

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC010.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 2 | [case 32](32/) |
| 3 | [case 46](46/) |
| 4 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` Y |
| 5 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 6 | [case 30](30/) |
| 7 | [case 29](29/) |
| 8 | [case 06](06/) |
