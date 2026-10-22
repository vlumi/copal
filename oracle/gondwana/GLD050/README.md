# GLD050: cases

11 cases reach 15 of the 55 rules of GLD050's decision tables, and of the programs it calls, 5 of the 16 of GLC020. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD050

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾌｧｲﾙ-ｶﾝﾘｮｳ` Y |
| 2 | [case 01](01/) |

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 04](04/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 07](07/) |
| 3 | [case 08](08/) |
| 4 | [case 09](09/) |
| 5 | [case 10](10/) |
| 6 | [case 11](11/) |
| 7 | [case 01](01/) |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 06](06/) |
| 3 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` Y |
| 4 | [case 01](01/) |
| 5 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` Y |
| 6 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` Y |
| 7 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` N, `WS-ｴﾝﾄﾞ-ﾏﾁ AND ZG-ｴﾝﾄﾞ` Y |
| 8 | not reached: `WS-PRMRES-ｵﾜﾘ` N, `WS-PRMRES-OK` Y, `WS-ﾍｯﾀﾞ-ﾏﾁ AND ZG-ﾍｯﾀﾞ AND ZG-H-ｼｭﾍﾞﾂｺｰﾄﾞ = 91` N, `WS-ﾍｯﾀﾞ-ﾏﾁ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾃﾞｰﾀ` N, `WS-ﾃﾞｰﾀ-ﾏﾁ AND ZG-ﾄﾚｰﾗ` N, `WS-ｴﾝﾄﾞ-ﾏﾁ AND ZG-ｴﾝﾄﾞ` N |

### [ﾀﾒﾙ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾀﾒﾙ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `ZG-D-ﾌﾘｶｴｽﾞﾐ` Y |
| 2 | not reached: `ZG-D-ﾌﾘｶｴｽﾞﾐ` N |

### [ｹﾝｻﾝ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｹﾝｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `ZG-T-ｺﾞｳｹｲｹﾝｽｳ NOT = WS-ｹﾝｻ-ｹﾝｽｳ OR ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｷﾝｶﾞｸ` Y |
| 2 | not reached: `ZG-T-ｺﾞｳｹｲｹﾝｽｳ NOT = WS-ｹﾝｻ-ｹﾝｽｳ OR ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｷﾝｶﾞｸ` N, `ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ NOT = WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ OR ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ` Y |
| 3 | not reached: `ZG-T-ｺﾞｳｹｲｹﾝｽｳ NOT = WS-ｹﾝｻ-ｹﾝｽｳ OR ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｷﾝｶﾞｸ` N, `ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ NOT = WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ OR ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ` N, `ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ NOT = WS-ｹﾝｻ-ﾌﾉｳｹﾝｽｳ OR ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ﾌﾉｳｷﾝｶﾞｸ` Y |
| 4 | not reached: `ZG-T-ｺﾞｳｹｲｹﾝｽｳ NOT = WS-ｹﾝｻ-ｹﾝｽｳ OR ZG-T-ｺﾞｳｹｲｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｷﾝｶﾞｸ` N, `ZG-T-ﾌﾘｶｴｽﾞﾐｹﾝｽｳ NOT = WS-ｹﾝｻ-ｽﾞﾐｹﾝｽｳ OR ZG-T-ﾌﾘｶｴｽﾞﾐｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ｽﾞﾐｷﾝｶﾞｸ` N, `ZG-T-ﾌﾘｶｴﾌﾉｳｹﾝｽｳ NOT = WS-ｹﾝｻ-ﾌﾉｳｹﾝｽｳ OR ZG-T-ﾌﾘｶｴﾌﾉｳｷﾝｶﾞｸ NOT = WS-ｹﾝｻ-ﾌﾉｳｷﾝｶﾞｸ` N |

### [ﾊﾝｴｲ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾊﾝｴｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-POLMST-ﾅｼ` Y, `WS-ﾘﾕｳ-ﾅｼ` N |
| 2 | not reached: `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 3 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` N, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` N |
| 4 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` Y, `WS-POLMST-OK, again` N |
| 5 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` Y, `WS-POLMST-OK, again` Y |
| 6 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` N, `WS-POLMST-OK, again` N |
| 7 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y, `ZG-D-ﾌﾘｶｴｽﾞﾐ` N, `WS-POLMST-OK, again` Y |
| 8 | not reached: `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ` N, `WS-POLMST-OK` Y, `PM-ﾊﾗｲｺﾐﾁｭｳ` Y |

### [ｷﾝｶﾞｸ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｷﾝｶﾞｸ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ` Y |
| 2 | not reached: `ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ` N, `PM-ﾂｷﾊﾞﾗｲ AND ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * 2` Y |
| 3 | not reached: `ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ` N, `PM-ﾂｷﾊﾞﾗｲ AND ZG-D-ｷﾝｶﾞｸ = PM-ﾎｹﾝﾘｮｳ * 2` N |

### [ﾌﾘｶｴｽﾞﾐ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴｽﾞﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` Y, `PM-ﾕｳﾖﾁｭｳ` Y |
| 2 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` Y, `PM-ﾕｳﾖﾁｭｳ` N |
| 3 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾕｳﾖﾁｭｳ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` Y |
| 4 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾕｳﾖﾁｭｳ` N, `WS-ﾊﾗｲｺﾐ-MM > 12` Y |
| 5 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾕｳﾖﾁｭｳ` Y, `WS-ﾊﾗｲｺﾐ-MM > 12` N |
| 6 | not reached: `PM-ﾈﾝﾊﾞﾗｲ` N, `PM-ﾕｳﾖﾁｭｳ` N, `WS-ﾊﾗｲｺﾐ-MM > 12` N |

### [ﾌﾘｶｴﾌﾉｳ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ﾌﾘｶｴﾌﾉｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾕｳｺｳ` Y, `ZG-D-ｼｷﾝﾌｿｸ` Y |
| 2 | not reached: `PM-ﾕｳｺｳ` Y, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` Y |
| 3 | not reached: `PM-ﾕｳｺｳ` Y, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` Y |
| 4 | not reached: `PM-ﾕｳｺｳ` Y, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` N, `ZG-D-ｲﾗｲｼｮﾅｼ` Y |
| 5 | not reached: `PM-ﾕｳｺｳ` Y, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` N, `ZG-D-ｲﾗｲｼｮﾅｼ` N |
| 6 | not reached: `PM-ﾕｳｺｳ` N, `ZG-D-ｼｷﾝﾌｿｸ` Y |
| 7 | not reached: `PM-ﾕｳｺｳ` N, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` Y |
| 8 | not reached: `PM-ﾕｳｺｳ` N, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` Y |
| 9 | not reached: `PM-ﾕｳｺｳ` N, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` N, `ZG-D-ｲﾗｲｼｮﾅｼ` Y |
| 10 | not reached: `PM-ﾕｳｺｳ` N, `ZG-D-ｼｷﾝﾌｿｸ` N, `ZG-D-ﾄﾘﾋｷﾅｼ` N, `ZG-D-ﾖｷﾝｼｬﾂｺﾞｳ` N, `ZG-D-ｲﾗｲｼｮﾅｼ` N |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD050.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
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
