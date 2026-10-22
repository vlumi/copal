# GLM010: cases

13 cases reach 23 of the 30 rules of GLM010's decision tables, and of the programs it calls, 5 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLM010

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 06](06/) |
| 3 | [case 07](07/) |
| 4 | [case 08](08/) |
| 5 | [case 09](09/) |
| 6 | [case 05](05/) |
| 7 | [case 01](01/) |
| 8 | [case 02](02/) |

### [ｾﾝﾀｸ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ｾﾝﾀｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ｾﾝﾀｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM010.md#ｾﾝﾀｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 02](02/) |

### [ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾊﾗｲｺﾐﾁｭｳ AND PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ < WS-ｾｲｷｭｳ-YYYYMM` Y, `PM-ﾈﾝﾊﾞﾗｲ` Y, `PM-ｹｲﾔｸﾋﾞ(4:2) = WS-ｾｲｷｭｳ-YYYYMM(5:2)` Y |
| 2 | [case 11](11/) |
| 3 | [case 04](04/) |
| 4 | not reached: `PM-ﾊﾗｲｺﾐﾁｭｳ AND PM-ﾊﾗｲｺﾐｽﾞﾐﾈﾝｹﾞﾂ < WS-ｾｲｷｭｳ-YYYYMM` Y, `PM-ﾈﾝﾊﾞﾗｲ` N, `WS-ｶｲｽｳ > 2` N |
| 5 | [case 02](02/) |

### [ｻｸｾｲ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ｻｸｾｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `RETURN SORT-F AT END` Y, `WS-ｷﾞﾝｺｳｹﾝｽｳ > 0` Y |
| 2 | [case 02](02/) |
| 3 | [case 04](04/) |
| 4 | not reached: `RETURN SORT-F AT END` N, `WS-ｷﾞﾝｺｳｹﾝｽｳ > 0` N |

### [ｻｸｾｲ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM010.md#ｻｸｾｲ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 10](10/) |

### [ﾃﾞｰﾀ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ﾃﾞｰﾀ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 13](13/) |
| 3 | [case 12](12/) |

### [PRMREQ-ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM010.md#prmreq-ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRMREQ-OK` N |
| 2 | [case 02](02/) |

### [ｾｲｷﾞｮ-RTN](../../../blueprints/gondwana/rules/GLM010.md#ｾｲｷﾞｮ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 02](02/) |

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
