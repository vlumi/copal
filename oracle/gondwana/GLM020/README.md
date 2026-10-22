# GLM020: cases

13 cases reach 29 of the 40 rules of GLM020's decision tables, and of the programs it calls, 5 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLM020

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |
| 3 | [case 04](04/) |
| 4 | [case 02](02/) |
| 5 | [case 05](05/) |
| 6 | [case 06](06/) |
| 7 | [case 07](07/) |
| 8 | [case 11](11/) |
| 9 | [case 08](08/) |

### [ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RSVTB-ｵﾜﾘ` N |
| 2 | [case 08](08/) |

### [ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 09](09/) |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-ｵﾜﾘ` N |
| 2 | [case 09](09/) |

### [ｹｲｻﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ｹｲｻﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 09](09/) |
| 3 | [case 13](13/) |

### [ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 13](13/) |
| 3 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `PM-ｹｲｶﾈﾝｽｳ > 50` Y |
| 4 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `PM-ｹｲｶﾈﾝｽｳ > 50` N |

### [ｷｼﾞｭﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｷｼﾞｭﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ｿﾞｳｹﾞﾝ = 0` N, `WS-POLMST-OK` N |
| 2 | [case 09](09/) |
| 3 | not reached: `SR-ｿﾞｳｹﾞﾝ = 0` Y, `WS-POLMST-OK` N |
| 4 | [case 13](13/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` Y, `WS-GLJRNL-OK` N |
| 2 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` Y, `WS-GLJRNL-OK` Y |
| 3 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` N, `WS-GLJRNL-OK` N |
| 4 | [case 09](09/) |

### [ﾎｳｺｸ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ﾎｳｺｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 09](09/) |

### [ﾎｳｺｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ﾎｳｺｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 09](09/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 09](09/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 09](09/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
| 2 | [case 08](08/) |

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
