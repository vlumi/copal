# GLM030: cases

22 cases reach 43 of the 49 rules of GLM030's decision tables, and of the programs it calls, 14 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile. A rule that finds a file to have failed once open says so: it takes the file's status after a read, a write, a rewrite, a delete or a close to be none that the file's records make, a record or the end for a read in order, a record or none of that key for a read by key, which is the device's doing. And a rule of a program called that takes a parameter to be what no call of this program passes, a direction it never asks for, says so, as do the rules of the paragraphs performed only from such rules: no run of this program reaches them, whatever its input.

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
| 8 | [case 02](02/) |
| 9 | [case 09](09/) |

### [ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾀﾞｲﾘﾃﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-AGTMST-ｵﾜﾘ` N |
| 2 | [case 02](02/) |

### [ｾﾝﾀｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｾﾝﾀｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-POLMST-ｵﾜﾘ` N |
| 2 | [case 09](09/) |

### [ｾﾝﾀｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM030.md#ｾﾝﾀｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 14](14/) |
| 3 | [case 09](09/) |
| 4 | [case 09](09/) |

### [ﾃｽｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾃｽｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 15](15/) |
| 3 | [case 15](15/) |
| 4 | [case 14](14/) |

### [ｼｶｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼｶｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 15](15/) |

### [ﾓﾄﾞｼｲﾚ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾓﾄﾞｼｲﾚ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 09](09/) |

### [ﾓﾄﾞｼｲﾚ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾓﾄﾞｼｲﾚ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 19](19/) |
| 2 | [case 22](22/) |
| 3 | [case 20](20/) |
| 4 | [case 21](21/) |
| 5 | [case 18](18/) |

### [ﾎｳｺｸ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾎｳｺｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 09](09/) |
| 2 | [case 14](14/) |
| 3 | [case 16](16/) |

### [ﾎｳｺｸ-RTN, in its loop](../../../blueprints/gondwana/rules/GLM030.md#ﾎｳｺｸ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 14](14/) |

### [ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `SR-ﾀﾞｲﾘﾃﾝｺｰﾄﾞ = WS-ｿｳﾀﾞｲﾘﾃﾝ` N, `WS-ｿｳﾀﾞｲﾘﾃﾝ = SPACES` N |
| 2 | [case 14](14/) |
| 3 | [case 14](14/) |

### [ﾐﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾐﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | [case 15](15/) |

### [ﾀﾞｲﾘﾃﾝ-ｹｲ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ﾀﾞｲﾘﾃﾝ-ｹｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-COMMIS-OK` N |
| 2 | [case 18](18/) |
| 3 | [case 15](15/) |
| 4 | [case 19](19/) |
| 5 | [case 14](14/) |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLJRNL-OK` N |
| 2 | [case 15](15/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLM030.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached, a file found to fail once open: `WS-GLSCTL-OK` N |
| 2 | [case 09](09/) |

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
| 2 | [case 11](11/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | [case 09](09/) |
| 3 | [case 13](13/) |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | [case 17](17/) |
| 6 | [case 12](12/) |
| 7 | [case 18](18/) |
