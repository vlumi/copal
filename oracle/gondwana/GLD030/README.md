# GLD030: cases

18 cases reach 23 of the 77 rules of GLD030's decision tables, and of the programs it calls, 6 of the 16 of GLC020, 0 of the 23 of GLC030, 0 of the 5 of GLC040. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD030

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 15](15/) |
| 3 | [case 10](10/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 02](02/) |
| 2 | [case 04](04/) |
| 3 | [case 03](03/) |
| 4 | [case 05](05/) |
| 5 | [case 06](06/) |
| 6 | [case 07](07/) |
| 7 | [case 08](08/) |
| 8 | [case 09](09/) |
| 9 | [case 10](10/) |

### [CLMREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#clmreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CLMREQ-OK` Y |
| 2 | not reached: `WS-CLMREQ-OK` N, `WS-CLMREQ-ｵﾜﾘ` N |
| 3 | [case 14](14/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘﾕｳ-ﾅｼ` Y |
| 2 | [case 10](10/) |

### [ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 10](10/) |
| 2 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 3 | [case 16](16/) |
| 4 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y |
| 5 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ` Y |
| 6 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ` N, `CR-ﾆｭｳｲﾝ` N |
| 7 | [case 11](11/) |
| 8 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 9 | [case 17](17/) |
| 10 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y, `CR-ﾆｭｳｲﾝ` Y |
| 11 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ` Y, `CR-ﾆｭｳｲﾝ` Y |
| 12 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ` N, `CR-ﾆｭｳｲﾝ` Y |
| 13 | [case 12](12/) |
| 14 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｶｲﾔｸ` Y |
| 15 | [case 18](18/) |
| 16 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｶｲﾔｸ` Y |
| 17 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｶｲﾔｸ` Y |
| 18 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `NOT CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｶｲﾔｸ` Y |
| 19 | [case 13](13/) |

### [ﾆｭｳｲﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` Y |
| 3 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 4 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 5 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 6 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` Y |
| 7 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 8 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｯｽｳ = 0` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` N |

### [ｼﾊﾗｲ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾊﾗｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N |
| 2 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y |
| 3 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N |
| 4 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y |
| 5 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` Y, `WS-POLMST-OK` N |
| 6 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` Y, `WS-POLMST-OK` Y |
| 7 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` Y, `WS-POLMST-OK` N |
| 8 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` Y, `WS-POLMST-OK` Y |
| 9 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` Y |
| 10 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y, `CR-ｶｲﾔｸ` Y |
| 11 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` Y |
| 12 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y, `CR-ｶｲﾔｸ` Y |
| 13 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` N |
| 14 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y, `CR-ｶｲﾔｸ` N |
| 15 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` N |
| 16 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ` N, `WS-POLMST-OK` Y, `CR-ｶｲﾔｸ` N |

### [ﾆｭｳｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ﾆｯｽｳ > 120` Y |
| 2 | not reached: `CR-ﾆｯｽｳ > 120` N |

### [ｶｲﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｶｲﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CSVTB-OK` N |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-ｷﾝｶﾞｸ > 0` Y |
| 3 | not reached: `WS-CSVTB-OK` Y, `WS-ｷﾝｶﾞｸ > 0` N |

### [ﾘｿｸ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘｿｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾘｿｸ > 0` Y |
| 2 | not reached: `WS-ﾘｿｸ > 0` N |

### [ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y |
| 2 | not reached: `CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0` Y, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N |
| 3 | not reached: `CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` Y |
| 4 | not reached: `CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | not reached: `WS-GLJRNL-OK` Y |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 10](10/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
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
| 2 | [case 16](16/) |

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

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` N, `WS-ﾖﾐｺﾐｽﾞﾐ, again` Y |
| 2 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` N, `WS-ﾖﾐｺﾐｽﾞﾐ, again` N |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` Y |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | not reached: `WS-CSVTB-OK` N |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` Y |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 2 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N |
| 3 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 4 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N |
| 5 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` N |
| 6 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 7 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 8 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 9 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 10 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 11 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 12 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 13 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 14 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 15 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` N, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` N |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾆｯｽｳ > 0` Y |
| 2 | not reached: `WS-ﾆｯｽｳ > 0` N |
