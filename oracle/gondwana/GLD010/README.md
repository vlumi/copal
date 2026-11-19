# GLD010: cases

27 cases reach 43 of the 61 rules of GLD010's decision tables, and of the programs it calls, 5 of the 16 of GLC020, 4 of the 8 of GLC010, 1 of the 14 of GLC050. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD010

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 07](07/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 08](08/) |
| 2 | [case 11](11/) |
| 3 | [case 12](12/) |
| 4 | [case 09](09/) |
| 5 | [case 10](10/) |
| 6 | [case 13](13/) |
| 7 | [case 14](14/) |
| 8 | [case 15](15/) |
| 9 | [case 16](16/) |
| 10 | [case 01](01/) |

### [NEWPOL-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD010.md#newpol-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 26](26/) |
| 2 | not reached: `WS-NEWPOL-OK` N, `WS-NEWPOL-ｵﾜﾘ` N |
| 3 | [case 06](06/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘﾕｳ-ﾅｼ` Y |
| 2 | [case 01](01/) |

### [ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 17](17/) |
| 3 | not reached: `NP-ｼｮｳﾋﾝ = SPACES` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` N, `WS-POLMST-ﾅｼ` N |
| 4 | [case 04](04/) |
| 5 | [case 20](20/) |
| 6 | [case 21](21/) |
| 7 | [case 18](18/) |
| 8 | [case 19](19/) |
| 9 | [case 01](01/) |

### [ｼｮｳﾋﾝ-ｹﾝｻｸ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｮｳﾋﾝ-ｹﾝｻｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 04](04/) |
| 3 | [case 03](03/) |

### [ﾈﾝﾚｲ-ｷﾝｶﾞｸ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾈﾝﾚｲ-ｷﾝｶﾞｸ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 22](22/) |
| 2 | not reached: `WS-ﾈﾝﾚｲ < WS-ｼｮｳﾋﾝ-ｻｲﾃｲﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX) OR WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-N = 1` N |
| 3 | [case 25](25/) |
| 4 | not reached: `WS-ﾈﾝﾚｲ < WS-ｼｮｳﾋﾝ-ｻｲﾃｲﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX) OR WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `WS-N = 1` N, `WS-ﾎｼｮｳ-ｷﾝｶﾞｸ(WS-N) < WS-ｼｮｳﾋﾝ-ｻｲﾃｲｷﾝｶﾞｸ(WS-ｼｮｳﾋﾝ-IX) OR WS-ﾎｼｮｳ-ｷﾝｶﾞｸ(WS-N) > WS-ｼｮｳﾋﾝ-ｻｲｺｳｷﾝｶﾞｸ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 5 | [case 19](19/) |
| 6 | [case 24](24/) |

### [ﾄｸﾔｸ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾄｸﾔｸ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(WS-ﾄｸﾔｸ-N) = SPACES` N, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(WS-ﾄｸﾔｸ-N)` Y, `WS-ﾘﾕｳ-ﾅｼ` Y |
| 2 | [case 24](24/) |
| 3 | not reached: `NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(WS-ﾄｸﾔｸ-N) = SPACES` N, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(WS-ﾄｸﾔｸ-N)` N, `WS-ﾘﾕｳ-ﾅｼ` N |
| 4 | [case 26](26/) |

### [ﾄｸﾔｸ-ｷｿｸ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾄｸﾔｸ-ｷｿｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 24](24/) |
| 2 | not reached: `NOT WS-ｼｮｳﾋﾝ-ﾄｸﾔｸ(WS-ｼｮｳﾋﾝ-IX) OR WS-ｼｮｳﾋﾝ-ｼｭﾙｲ(WS-ｼｮｳﾋﾝ-IX) = WS-ｼｭｹｲﾔｸ-ｼｭﾙｲ OR (WS-ﾄｸﾔｸ-N = 2 AND NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(2) = NP-ﾄｸﾔｸ-ｼｮｳﾋﾝ(1))` N |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-OK` N |
| 2 | not reached: `WS-POLMST-OK` Y |

### [ﾎｼｮｳ-ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾎｼｮｳ-ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 26](26/) |
| 2 | not reached: `WS-COVMST-OK` Y |

### [ﾎｹﾝﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾎｹﾝﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NP-ﾊﾗｲｺﾐﾎｳﾎｳ = '2'` Y |
| 2 | not reached: `NP-ﾊﾗｲｺﾐﾎｳﾎｳ = '2'` N, `NP-ﾊﾗｲｺﾐﾎｳﾎｳ = '3'` Y |
| 3 | [case 26](26/) |

### [ﾏﾝｷﾋﾞ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾏﾝｷﾋﾞ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ﾎｹﾝｷｶﾝ(WS-ｼｮｳﾋﾝ-IX) = 0` Y |
| 2 | [case 26](26/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
| 2 | [case 01](01/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RATETB-ｵﾜﾘ` N |
| 2 | [case 01](01/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN, in its loop](../../../blueprints/gondwana/rules/GLD010.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 03](03/) |

### [ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD010.md#ｼｭｳﾉｳ-ﾖﾋﾞﾀﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRMHST-OK` N |
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
| 1 | [case 05](05/) |
| 2 | [case 01](01/) |

### [ﾜﾚｷ-ﾍﾝｶﾝ-RTN, in its loop](../../../blueprints/gondwana/rules/GLC020.md#ﾜﾚｷ-ﾍﾝｶﾝ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 05](05/) |

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

## GLC010

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC010.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 2 | [case 22](22/) |
| 3 | [case 27](27/) |
| 4 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` Y |
| 5 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 6 | [case 23](23/) |
| 7 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` N |
| 8 | [case 19](19/) |

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
