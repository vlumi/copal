# GLD020: cases

21 cases reach 32 of the 80 rules of GLD020's decision tables, and of the programs it calls, 6 of the 16 of GLC020, 1 of the 5 of GLC040, 10 of the 23 of GLC030, 0 of the 8 of GLC010. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD020

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 12](12/) |
| 3 | [case 10](10/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 04](04/) |
| 2 | [case 05](05/) |
| 3 | [case 02](02/) |
| 4 | [case 03](03/) |
| 5 | [case 06](06/) |
| 6 | [case 07](07/) |
| 7 | [case 08](08/) |
| 8 | [case 09](09/) |
| 9 | [case 10](10/) |

### [POLMST-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#polmst-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `NOT WS-POLMST-OK AND NOT WS-POLMST-ｵﾜﾘ` Y |
| 2 | [case 10](10/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾎｼｮｳﾁｭｳ` Y, `WS-ｵｳﾄｳ-ｱﾘ` Y, `WS-POLMST-OK` N |
| 2 | [case 13](13/) |
| 3 | [case 11](11/) |
| 4 | [case 10](10/) |

### [ｵｳﾄｳﾋﾞ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｳﾄｳﾋﾞ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｵｳﾄｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ｵｳﾄｳﾋﾞ > WS-ｹｲﾔｸﾋﾞ AND (WS-ｵｳﾄｳﾋﾞ > CT-ｷﾞｮｳﾑﾋﾂﾞｹ OR WS-ｵｳﾄｳﾋﾞ = WS-ｷﾞｮｳﾑﾋﾂﾞｹ)` Y |
| 2 | not reached: `WS-ｵｳﾄｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ｵｳﾄｳﾋﾞ > WS-ｹｲﾔｸﾋﾞ AND (WS-ｵｳﾄｳﾋﾞ > CT-ｷﾞｮｳﾑﾋﾂﾞｹ OR WS-ｵｳﾄｳﾋﾞ = WS-ｷﾞｮｳﾑﾋﾂﾞｹ)` N |
| 3 | [case 13](13/) |
| 4 | [case 11](11/) |

### [ｵｳﾄｳ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｵｳﾄｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` Y |
| 3 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾏﾝｷ` N |
| 4 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾏﾝｷ` Y |
| 5 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` N |
| 6 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y, `PM-ﾏﾝｷ` Y |
| 7 | [case 13](13/) |
| 8 | not reached: `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾎｼｮｳﾁｭｳ AND PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N, `PM-ﾏﾝｷ` Y |

### [ｼｭｹｲﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｹｲﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 15](15/) |
| 2 | [case 17](17/) |
| 3 | [case 14](14/) |
| 4 | not reached: `WS-COVMST-OK` Y, `PM-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y, `HS-ﾕｳｺｳ` Y, `PM-ﾖｳﾛｳ` Y |
| 5 | not reached: `WS-COVMST-OK` Y, `PM-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y, `HS-ﾕｳｺｳ` N, `PM-ﾖｳﾛｳ` Y |
| 6 | not reached: `WS-COVMST-OK` Y, `PM-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y, `HS-ﾕｳｺｳ` Y, `PM-ﾖｳﾛｳ` N |
| 7 | not reached: `WS-COVMST-OK` Y, `PM-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y, `HS-ﾕｳｺｳ` N, `PM-ﾖｳﾛｳ` N |

### [ｼｭｹｲﾔｸ-ｺｳｼﾝ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｹｲﾔｸ-ｺｳｼﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `HS-ﾕｳｺｳ` Y |
| 2 | not reached: `HS-ﾕｳｺｳ` N |

### [ﾄｸﾔｸ-ｿｳｻ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-ｿｳｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 2 | not reached: `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` N |
| 3 | not reached: `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 4 | [case 13](13/) |

### [ﾄｸﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` Y, `HS-ﾕｳｺｳ, again` Y |
| 2 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` Y, `HS-ﾕｳｺｳ, again` N |
| 3 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` N, `HS-ﾕｳｺｳ, again` Y, `HS-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y |
| 4 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` N, `HS-ﾕｳｺｳ, again` N, `HS-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` Y |
| 5 | not reached: `HS-ﾕｳｺｳ` Y, `PM-ﾏﾝｷ` N, `HS-ﾏﾝｷﾋﾞ = WS-ｵｳﾄｳﾋﾞ-ﾜﾚｷ` N |
| 6 | not reached: `HS-ﾕｳｺｳ` N |

### [ﾄｸﾔｸ-ﾒｲｻｲ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾄｸﾔｸ-ﾒｲｻｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `HS-ﾕｳｺｳ` Y |
| 2 | not reached: `HS-ﾕｳｺｳ` N |

### [ﾎｼｮｳ-ｶｷﾓﾄﾞｼ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾎｼｮｳ-ｶｷﾓﾄﾞｼ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-COVMST-OK` N |
| 2 | not reached: `WS-COVMST-OK` Y |

### [ﾘｿｸ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾘｿｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘｿｸ > 0` Y |
| 2 | [case 14](14/) |

### [ﾁｮｳｶ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾁｮｳｶ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 17](17/) |
| 2 | [case 16](16/) |
| 3 | not reached: `WS-CSVTB-OK` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > WS-ﾍﾝﾚｲｷﾝ` N |

### [ﾏﾝｷ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾏﾝｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾊｲﾄｳﾙｲｹｲ > 0` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y |
| 2 | not reached: `PM-ﾊｲﾄｳﾙｲｹｲ > 0` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N |
| 3 | not reached: `PM-ﾊｲﾄｳﾙｲｹｲ > 0` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y |
| 4 | not reached: `PM-ﾊｲﾄｳﾙｲｹｲ > 0` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N |

### [ｺｳｼﾝ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｺｳｼﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 3 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y |
| 4 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |
| 5 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 6 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` Y |
| 7 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` Y |
| 8 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ﾈﾝﾚｲ > WS-ｼｮｳﾋﾝ-ｻｲｺｳﾈﾝﾚｲ(WS-ｼｮｳﾋﾝ-IX)` N, `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾊﾝﾄｼﾊﾞﾗｲ` N |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | not reached: `WS-GLJRNL-OK` Y |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 10](10/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
| 2 | [case 10](10/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD020.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-RATETB-ｵﾜﾘ` N |
| 2 | [case 10](10/) |

### [ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-RTN, in its loop](../../../blueprints/gondwana/rules/GLD020.md#ﾘｮｳﾘﾂ-ﾖﾐｺﾐ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 11](11/) |
| 2 | [case 10](10/) |

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
| 2 | [case 11](11/) |

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

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | [case 14](14/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾆｯｽｳ > 0` Y |
| 2 | not reached: `WS-ﾆｯｽｳ > 0` N |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | [case 17](17/) |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 16](16/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 17](17/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 2 | [case 20](20/) |
| 3 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 4 | [case 16](16/) |
| 5 | [case 21](21/) |
| 6 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 7 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 8 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 9 | [case 18](18/) |
| 10 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 11 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 12 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 13 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 14 | [case 19](19/) |
| 15 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |

## GLC010

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC010.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 2 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` Y |
| 3 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` N |
| 4 | not reached: `WS-A-DD < WS-B-DD` Y, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` Y |
| 5 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` N |
| 6 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` Y, `WS-ｱﾏﾘ < 6` Y |
| 7 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` N |
| 8 | not reached: `WS-A-DD < WS-B-DD` N, `WS-ﾂｷｽｳ < 0` N, `WS-ｱﾏﾘ < 6` Y |
