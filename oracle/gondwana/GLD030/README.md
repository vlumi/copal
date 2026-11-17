# GLD030: cases

20 cases reach 25 of the 128 rules of GLD030's decision tables, and of the programs it calls, 6 of the 16 of GLC020, 0 of the 23 of GLC030, 0 of the 5 of GLC040. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD030

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 16](16/) |
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
| 8 | not reached: `WS-CLMREQ-OK` Y, `WS-ｴﾗｰﾌｧｲﾙ = SPACES` N, `WS-PRODTB-OK` Y, `WS-POLMST-OK` Y, `WS-GLSCTL-OK` Y, `WS-PAYOUT-OK` Y, `WS-GLJRNL-OK` Y, `WS-REJECT-OK` Y, `WS-BENHST-OK` N |
| 9 | [case 09](09/) |
| 10 | [case 10](10/) |

### [CLMREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#clmreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CLMREQ-OK` Y |
| 2 | not reached: `WS-CLMREQ-OK` N, `WS-CLMREQ-ｵﾜﾘ` N |
| 3 | [case 15](15/) |

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
| 3 | [case 17](17/) |
| 4 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y |
| 5 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N |
| 6 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` Y, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N |
| 7 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N |
| 8 | [case 11](11/) |
| 9 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 10 | [case 18](18/) |
| 11 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y, `CR-ﾆｭｳｲﾝ` Y |
| 12 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ﾆｭｳｲﾝ` Y |
| 13 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` Y, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` Y |
| 14 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` Y |
| 15 | [case 12](12/) |
| 16 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 17 | [case 19](19/) |
| 18 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 19 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 20 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` Y, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 21 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 22 | [case 13](13/) |
| 23 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 24 | [case 20](20/) |
| 25 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` Y, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 26 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ, again` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` Y, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 27 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` Y, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 28 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-ﾘﾕｳ-ﾅｼ AND WS-ｼﾞﾕｳﾋﾞ > WS-ｷﾞｮｳﾑﾋﾂﾞｹ` N, `WS-ﾘﾕｳ-ﾅｼ AND NOT CR-ｶｲﾔｸ` N, `WS-POLMST-OK` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｹｲﾔｸﾋﾞ` N, `CR-ｶｲﾔｸ AND NOT PM-ﾎｼｮｳﾁｭｳ AND NOT PM-ｼｯｺｳ` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 29 | [case 14](14/) |

### [ﾎｼｮｳ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾎｼｮｳﾁｭｳ` Y |
| 2 | not reached: `PM-ﾎｼｮｳﾁｭｳ` N, `PM-ｼｯｺｳ` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｯｺｳﾋﾞ` N |
| 3 | not reached: `PM-ﾎｼｮｳﾁｭｳ` N, `PM-ｼｯｺｳ` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｯｺｳﾋﾞ` Y |
| 4 | not reached: `PM-ﾎｼｮｳﾁｭｳ` N, `PM-ｼｯｺｳ` N |

### [ｼｮｳﾋﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｳﾋﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y |
| 3 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` Y, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N |
| 4 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 5 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y |
| 6 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = PM-ｼｮｳﾋﾝ` N, `CR-ｼﾎﾞｳ AND WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N |

### [ﾆｭｳｲﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｭｳｲﾝ AND CR-ﾆｯｽｳ = 0` Y |
| 3 | not reached: `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｭｳｲﾝ AND CR-ﾆｯｽｳ = 0` N, `CR-ｼｭｼﾞｭﾂ AND CR-ﾆｯｽｳ NOT = 10 AND CR-ﾆｯｽｳ NOT = 20 AND CR-ﾆｯｽｳ NOT = 40` Y |
| 4 | not reached: `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｭｳｲﾝ AND CR-ﾆｯｽｳ = 0` N, `CR-ｼｭｼﾞｭﾂ AND CR-ﾆｯｽｳ NOT = 10 AND CR-ﾆｯｽｳ NOT = 20 AND CR-ﾆｯｽｳ NOT = 40` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` Y |
| 5 | not reached: `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `CR-ﾆｭｳｲﾝ AND CR-ﾆｯｽｳ = 0` N, `CR-ｼｭｼﾞｭﾂ AND CR-ﾆｯｽｳ NOT = 10 AND CR-ﾆｯｽｳ NOT = 20 AND CR-ﾆｯｽｳ NOT = 40` N, `WS-ｹｲｶﾆｯｽｳ < WS-ｼｮｳﾋﾝ-ﾒﾝｾｷﾆｯｽｳ(WS-ｼｮｳﾋﾝ-IX)` N |

### [ﾘﾚｷ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘﾚｷ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-BENHST-OK` Y, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` Y |
| 2 | not reached: `WS-BENHST-OK` Y, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` N, `WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000` Y |
| 3 | not reached: `WS-BENHST-OK` Y, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` N, `WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000` N |
| 4 | not reached: `WS-BENHST-OK` N, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` Y |
| 5 | not reached: `WS-BENHST-OK` N, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` N, `WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000` Y |
| 6 | not reached: `WS-BENHST-OK` N, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` N, `WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000` N |

### [ﾘﾚｷ-ｹﾝｻ-RTN, in its loop](../../../blueprints/gondwana/rules/GLD030.md#ﾘﾚｷ-ｹﾝｻ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `BH-ﾆｭｳｲﾝ` Y, `BH-ｼｭﾍﾞﾂ = CR-ｼｭﾍﾞﾂ AND BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ = CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ` Y |
| 2 | not reached: `BH-ﾆｭｳｲﾝ` Y, `BH-ｼｭﾍﾞﾂ = CR-ｼｭﾍﾞﾂ AND BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ = CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ` N |
| 3 | not reached: `BH-ﾆｭｳｲﾝ` N, `BH-ｼｭﾍﾞﾂ = CR-ｼｭﾍﾞﾂ AND BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ = CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ` Y |
| 4 | not reached: `BH-ﾆｭｳｲﾝ` N, `BH-ｼｭﾍﾞﾂ = CR-ｼｭﾍﾞﾂ AND BH-ｼﾞﾕｳﾊｯｾｲﾋﾞ = CR-ｼﾞﾕｳﾊｯｾｲﾋﾞ` N |

### [ｼﾊﾗｲ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾊﾗｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N |
| 2 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y |
| 3 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N |
| 4 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y |
| 5 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 6 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` Y |
| 7 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 8 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` Y |
| 9 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 10 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 11 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 12 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 13 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 14 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 15 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 16 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 17 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 18 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 19 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 20 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |

### [ｼﾎﾞｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾎﾞｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-ﾒﾝｾｷ-ｱﾘ` Y |
| 2 | not reached: `WS-ﾒﾝｾｷ-ｱﾘ` N, `WS-ﾐﾉｳｶﾞｸ > 0` Y |
| 3 | not reached: `WS-ﾒﾝｾｷ-ｱﾘ` N, `WS-ﾐﾉｳｶﾞｸ > 0` N |

### [ﾒﾝｾｷ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾒﾝｾｷ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ｼﾞｻﾂ` Y, `PM-ﾌｯｶﾂﾋﾞ = SPACES` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ﾒﾝｾｷｷｹﾞﾝ` Y |
| 2 | not reached: `CR-ｼﾞｻﾂ` Y, `PM-ﾌｯｶﾂﾋﾞ = SPACES` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ﾒﾝｾｷｷｹﾞﾝ` N |
| 3 | not reached: `CR-ｼﾞｻﾂ` Y, `PM-ﾌｯｶﾂﾋﾞ = SPACES` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ﾒﾝｾｷｷｹﾞﾝ` Y |
| 4 | not reached: `CR-ｼﾞｻﾂ` Y, `PM-ﾌｯｶﾂﾋﾞ = SPACES` Y, `WS-ｼﾞﾕｳﾋﾞ < WS-ﾒﾝｾｷｷｹﾞﾝ` N |
| 5 | not reached: `CR-ｼﾞｻﾂ` N |

### [ﾒﾝｾｷ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾒﾝｾｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-CSVTB-OK` N |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-ｷﾝｶﾞｸ > 0` Y |
| 3 | not reached: `WS-CSVTB-OK` Y, `WS-ｷﾝｶﾞｸ > 0` N |

### [ﾐﾉｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾐﾉｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾂｷﾊﾞﾗｲ` N, `PM-ﾕｳｺｳ` N |
| 2 | not reached: `PM-ﾂｷﾊﾞﾗｲ` N, `PM-ﾕｳｺｳ` Y |
| 3 | not reached: `PM-ﾂｷﾊﾞﾗｲ` Y, `WS-ﾐﾉｳﾂｷｽｳ < 0` Y |
| 4 | not reached: `PM-ﾂｷﾊﾞﾗｲ` Y, `WS-ﾐﾉｳﾂｷｽｳ < 0` N |

### [ﾘﾚｷ-ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘﾚｷ-ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-BENHST-OK` N |
| 2 | not reached: `WS-BENHST-OK` Y |

### [ﾆｭｳｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ﾆｯｽｳ > 120` Y, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` Y |
| 2 | not reached: `CR-ﾆｯｽｳ > 120` Y, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` N |
| 3 | not reached: `CR-ﾆｯｽｳ > 120` N, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` Y |
| 4 | not reached: `CR-ﾆｯｽｳ > 120` N, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` N |

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
| 2 | [case 17](17/) |

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
