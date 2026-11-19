# GLM030: cases

18 cases reach 42 of the 49 rules of GLM030's decision tables, and of the programs it calls, 8 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLM030

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |
| 3 | [case 04](04/) |
| 4 | [case 05](05/) |
| 5 | [case 06](06/) |
| 6 | [case 07](07/) |
| 7 | [case 08](08/) |
| 8 | [case 13](13/) |
| 9 | [case 02](02/) |

### [ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-AGTMST-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ｾﾝﾀｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｾﾝﾀｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ｾﾝﾀｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM030.md#ｾﾝﾀｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 02](02/) |
| 3 | [case 14](14/) |
| 4 | [case 14](14/) |

### [ﾃｽｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾃｽｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 02](02/) |
| 3 | [case 02](02/) |
| 4 | [case 11](11/) |

### [ｼｶｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼｶｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 12](12/) |
| 2 | [case 02](02/) |

### [ﾓﾄﾞｼｲﾚ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾓﾄﾞｼｲﾚ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 14](14/) |

### [ﾓﾄﾞｼｲﾚ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾓﾄﾞｼｲﾚ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 18](18/) |
| 3 | [case 16](16/) |
| 4 | [case 17](17/) |
| 5 | [case 09](09/) |

### [ﾎｳｺｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾎｳｺｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 02](02/) |
| 3 | [case 12](12/) |

### [ﾎｳｺｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM030.md#ﾎｳｺｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 02](02/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ = WS-ｿｳﾀﾞｲﾘﾃﾝ` N, `WS-ｿｳﾀﾞｲﾘﾃﾝ = SPACES` N |
| 2 | [case 02](02/) |
| 3 | [case 02](02/) |

### [ﾐﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾐﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾀﾞｲﾘﾃﾝ-ｺｰﾄﾞ(WS-ﾀﾞｲﾘﾃﾝ-IX) = SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ` Y |
| 2 | [case 02](02/) |

### [ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾀﾞｲﾘﾃﾝ-ｹｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-COMMIS-OK` N |
| 2 | [case 09](09/) |
| 3 | [case 02](02/) |
| 4 | [case 10](10/) |
| 5 | [case 11](11/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | [case 02](02/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 02](02/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 09](09/) |
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
| 1 | [case 09](09/) |
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
| 7 | [case 09](09/) |
