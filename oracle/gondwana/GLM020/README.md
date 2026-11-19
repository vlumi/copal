# GLM020: cases

21 cases reach 40 of the 53 rules of GLM020's decision tables, and of the programs it calls, 5 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLM020

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 14](14/) |
| 3 | [case 15](15/) |
| 4 | [case 13](13/) |
| 5 | [case 16](16/) |
| 6 | [case 17](17/) |
| 7 | [case 18](18/) |
| 8 | [case 19](19/) |
| 9 | [case 12](12/) |
| 10 | [case 01](01/) |

### [ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RSVTB-ｵﾜﾘ` N |
| 2 | [case 01](01/) |

### [ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ｼﾞｭﾝﾋﾞｷﾝ-ﾖﾐｺﾐ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 01](01/) |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-ｵﾜﾘ` N |
| 2 | [case 01](01/) |

### [ｹｲｻﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ｹｲｻﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |
| 3 | [case 06](06/) |

### [ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ｹｲｶﾈﾝｽｳ > 50` Y, `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 2 | [case 05](05/) |
| 3 | not reached: `PM-ｹｲｶﾈﾝｽｳ > 50` Y, `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 4 | [case 05](05/) |
| 5 | not reached: `PM-ｹｲｶﾈﾝｽｳ > 50` N, `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 6 | [case 01](01/) |
| 7 | not reached: `PM-ｹｲｶﾈﾝｽｳ > 50` N, `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 8 | [case 01](01/) |

### [ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ｾｷﾆﾝｼﾞｭﾝﾋﾞｷﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 07](07/) |

### [ﾎｼｮｳ-ｼﾞｭﾝﾋﾞｷﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ﾎｼｮｳ-ｼﾞｭﾝﾋﾞｷﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 01](01/) |
| 3 | [case 21](21/) |
| 4 | [case 08](08/) |
| 5 | [case 20](20/) |
| 6 | [case 04](04/) |

### [ｷｼﾞｭﾝ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｷｼﾞｭﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ｿﾞｳｹﾞﾝ = 0` N, `WS-POLMST-OK` N |
| 2 | [case 01](01/) |
| 3 | not reached: `SR-ｿﾞｳｹﾞﾝ = 0` Y, `WS-POLMST-OK` N |
| 4 | [case 06](06/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` Y, `WS-GLJRNL-OK` N |
| 2 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` Y, `WS-GLJRNL-OK` Y |
| 3 | not reached: `SR-ｿﾞｳｹﾞﾝ > 0` N, `WS-GLJRNL-OK` N |
| 4 | [case 01](01/) |

### [ﾎｳｺｸ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ﾎｳｺｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 01](01/) |

### [ﾎｳｺｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM020.md#ﾎｳｺｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 01](01/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 02](02/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM020.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
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
| 1 | [case 10](10/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 10](10/) |

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
