# GLD030: cases

46 cases reach 97 of the 153 rules of GLD030's decision tables, and of the programs it calls, 9 of the 16 of GLC020, 14 of the 23 of GLC030, 4 of the 5 of GLC040. Each case is a directory: the input files, field by field, the files the program wrote and what it displayed, and what it returned. The cases were found by running the program, compiled with GnuCOBOL's trace code, on inputs made from its own conditions and with files that cannot be opened or read, and are as few as reach the same rules. A rule no case reaches is listed with the outcomes of its conditions: most are rules no run can reach, as a table does not weigh what one condition's outcome says about another's. A rule that finds two totals the program adds up itself to differ, items of its own that it only adds to, says so: no input field sets either, so the search has nothing to change to reach it. A rule that finds a write or rewrite by key to fail says so too: the search fails a file only when it is opened or read, and a program that writes keys it made new, or rewrites the record it read, meets that failure only when something else changes the file meanwhile.

## GLD030

### [ｼｮｷ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 03](03/) |
| 2 | [case 04](04/) |
| 3 | [case 01](01/) |

### [ｵｰﾌﾟﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｵｰﾌﾟﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 05](05/) |
| 2 | [case 08](08/) |
| 3 | [case 06](06/) |
| 4 | [case 07](07/) |
| 5 | [case 09](09/) |
| 6 | [case 10](10/) |
| 7 | [case 11](11/) |
| 8 | [case 12](12/) |
| 9 | not reached: `WS-CLMREQ-OK` Y, `WS-ｴﾗｰﾌｧｲﾙ = SPACES` N, `WS-PRODTB-OK` Y, `WS-POLMST-OK` Y, `WS-COVMST-OK` Y, `WS-GLSCTL-OK` Y, `WS-PAYOUT-OK` Y, `WS-GLJRNL-OK` Y, `WS-REJECT-OK` Y, `WS-BENHST-OK` N |
| 10 | [case 13](13/) |
| 11 | [case 01](01/) |

### [CLMREQ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#clmreq-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `WS-CLMREQ-OK` N, `WS-CLMREQ-ｵﾜﾘ` N |
| 3 | [case 01](01/) |

### [ｼｮﾘ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 21](21/) |
| 2 | [case 01](01/) |

### [ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N |
| 3 | [case 15](15/) |
| 4 | [case 16](16/) |
| 5 | [case 17](17/) |
| 6 | [case 18](18/) |
| 7 | [case 14](14/) |
| 8 | [case 01](01/) |
| 9 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 10 | [case 15](15/) |
| 11 | [case 16](16/) |
| 12 | [case 17](17/) |
| 13 | [case 18](18/) |
| 14 | [case 14](14/) |
| 15 | [case 01](01/) |
| 16 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` Y, `CR-ﾆｭｳｲﾝ` N |
| 17 | [case 15](15/) |
| 18 | [case 16](16/) |
| 19 | [case 17](17/) |
| 20 | [case 24](24/) |
| 21 | [case 14](14/) |
| 22 | [case 01](01/) |
| 23 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ﾘﾕｳ-ﾅｼ` Y, `WS-POLMST-ﾅｼ` N, `WS-POLMST-OK` N, `CR-ｶｲﾔｸ` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 24 | [case 15](15/) |
| 25 | [case 16](16/) |
| 26 | [case 17](17/) |
| 27 | [case 33](33/) |
| 28 | [case 21](21/) |
| 29 | [case 02](02/) |

### [ﾎｼｮｳ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 14](14/) |
| 3 | [case 20](20/) |
| 4 | [case 21](21/) |

### [ﾎｼｮｳ-ｿｳｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ｿｳｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 2 | [case 18](18/) |
| 3 | [case 18](18/) |
| 4 | [case 23](23/) |
| 5 | not reached: `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 6 | [case 21](21/) |
| 7 | [case 21](21/) |
| 8 | not reached: `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` N, `CR-ｼﾎﾞｳ AND WS-ｼﾎﾞｳｷﾝｶﾞｸ = 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N |

### [ﾎｼｮｳ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 32](32/) |
| 2 | [case 31](31/) |
| 3 | not reached: `HS-ﾕｳｺｳ` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｮｳﾒﾂﾋﾞ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` Y |
| 4 | not reached: `HS-ﾕｳｺｳ` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｮｳﾒﾂﾋﾞ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` N |
| 5 | not reached: `HS-ﾕｳｺｳ` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｮｳﾒﾂﾋﾞ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 6 | not reached: `HS-ﾕｳｺｳ` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｮｳﾒﾂﾋﾞ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` Y |
| 7 | not reached: `HS-ﾕｳｺｳ` N, `WS-ｼﾞﾕｳﾋﾞ < WS-ｼｮｳﾒﾂﾋﾞ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` N |
| 8 | [case 23](23/) |
| 9 | [case 33](33/) |
| 10 | not reached: `HS-ﾕｳｺｳ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` Y, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` N |
| 11 | not reached: `HS-ﾕｳｺｳ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` N |
| 12 | not reached: `HS-ﾕｳｺｳ` Y, `WS-ﾎｼｮｳ-ｱﾘ` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = HS-ｼｮｳﾋﾝ` N, `WS-ｼｮｳﾋﾝ-ｲﾘｮｳ(WS-ｼｮｳﾋﾝ-IX)` Y, `WS-ﾆｭｳｲﾝﾎｼｮｳ-ﾅｼ` Y |
| 13 | [case 22](22/) |

### [ﾆｭｳｲﾝ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 34](34/) |
| 3 | [case 33](33/) |
| 4 | [case 36](36/) |
| 5 | [case 35](35/) |

### [ﾘﾚｷ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘﾚｷ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-BENHST-OK` Y, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` Y |
| 2 | not reached: `WS-BENHST-OK` Y, `NOT WS-BENHST-OK AND NOT WS-BENHST-ｵﾜﾘ` N, `WS-ﾘﾕｳ-ﾅｼ AND CR-ﾆｭｳｲﾝ AND WS-ﾂｳｻﾝﾆｯｽｳ NOT < 1000` Y |
| 3 | [case 35](35/) |
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
| 2 | [case 24](24/) |
| 3 | not reached: `CR-ｼﾎﾞｳ` Y, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N |
| 4 | [case 25](25/) |
| 5 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 6 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` Y |
| 7 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` Y |
| 8 | [case 35](35/) |
| 9 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 10 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 11 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 12 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` Y, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` Y |
| 13 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 14 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 15 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` Y |
| 16 | [case 21](21/) |
| 17 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 18 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` Y, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 19 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` N, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |
| 20 | not reached: `CR-ｼﾎﾞｳ` N, `WS-ｷﾝｶﾞｸ > 0` N, `CR-ﾆｭｳｲﾝ OR CR-ｼｭｼﾞｭﾂ` N, `WS-POLMST-OK` Y, `CR-ﾆｭｳｲﾝ` N, `CR-ｼｭｼﾞｭﾂ` N, `CR-ｶｲﾔｸ` N |

### [ｼﾎﾞｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾎﾞｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 25](25/) |
| 2 | [case 38](38/) |
| 3 | [case 24](24/) |

### [ﾒﾝｾｷ-ﾊﾝﾃｲ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾒﾝｾｷ-ﾊﾝﾃｲ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 25](25/) |
| 2 | [case 24](24/) |
| 3 | [case 23](23/) |
| 4 | [case 41](41/) |
| 5 | [case 37](37/) |

### [ﾒﾝｾｷ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾒﾝｾｷ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 26](26/) |
| 2 | [case 42](42/) |
| 3 | [case 25](25/) |

### [ﾐﾉｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾐﾉｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `PM-ﾂｷﾊﾞﾗｲ` N, `PM-ﾕｳｺｳ` N |
| 2 | [case 40](40/) |
| 3 | [case 39](39/) |
| 4 | [case 24](24/) |

### [ﾘﾚｷ-ｶｷｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘﾚｷ-ｶｷｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-BENHST-OK` N |
| 2 | [case 35](35/) |

### [ﾆｭｳｲﾝ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾆｭｳｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `CR-ﾆｯｽｳ > 120` Y, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` Y |
| 2 | not reached: `CR-ﾆｯｽｳ > 120` Y, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` N |
| 3 | not reached: `CR-ﾆｯｽｳ > 120` N, `WS-ﾂｳｻﾝﾆｯｽｳ + WS-ﾆｯｽｳ > 1000` Y |
| 4 | [case 35](35/) |

### [ｶｲﾔｸ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｶｲﾔｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 19](19/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-ｷﾝｶﾞｸ > 0` Y |
| 3 | [case 21](21/) |

### [ﾎｼｮｳ-ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-COVMST-OK` Y, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 2 | [case 24](24/) |
| 3 | not reached: `WS-COVMST-OK` N, `NOT WS-COVMST-OK AND NOT WS-COVMST-ｵﾜﾘ` Y |
| 4 | [case 21](21/) |

### [ﾎｼｮｳ-ｼｭｳﾘｮｳ-RTN, in its loop](../../../blueprints/gondwana/rules/GLD030.md#ﾎｼｮｳ-ｼｭｳﾘｮｳ-rtn-in-its-loop)

| Rule | Reached by |
| --: | --- |
| 1 | [case 24](24/) |
| 2 | not reached: `HS-ﾕｳｺｳ` Y, `CR-ｼﾎﾞｳ` Y, `WS-COVMST-OK` N |
| 3 | [case 36](36/) |
| 4 | not reached: `HS-ﾕｳｺｳ` Y, `CR-ｼﾎﾞｳ` N, `WS-COVMST-OK` N |
| 5 | [case 32](32/) |

### [ﾘｿｸ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾘｿｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 14](14/) |

### [ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ﾊｲﾄｳ-ｶｼﾂｹ-ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 24](24/) |
| 2 | [case 42](42/) |
| 3 | [case 37](37/) |
| 4 | not reached: `CR-ｼﾎﾞｳ AND PM-ﾊｲﾄｳﾙｲｹｲ > 0` N, `PM-ｶｼﾂｹｻﾞﾝﾀﾞｶ > 0` N |

### [ｼﾜｹ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼﾜｹ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLJRNL-OK` N |
| 2 | [case 18](18/) |

### [ｼｭｳﾘｮｳ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｭｳﾘｮｳ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-GLSCTL-OK` N |
| 2 | [case 01](01/) |

### [ｼｮｳﾋﾝ-ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLD030.md#ｼｮｳﾋﾝ-ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `WS-PRODTB-ｵﾜﾘ` N |
| 2 | [case 01](01/) |

## GLC020

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 01](01/) |
| 2 | [case 14](14/) |
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
| 1 | [case 14](14/) |
| 2 | [case 15](15/) |

### [ﾊﾝｲ-ｹﾝｻ-RTN](../../../blueprints/gondwana/rules/GLC020.md#ﾊﾝｲ-ｹﾝｻ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ﾈﾝ = 0` Y, `LK-ﾍﾝｶﾝｽﾞﾐ` N |
| 2 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` Y |
| 3 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` Y |
| 4 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` Y |
| 5 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` N, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` N |
| 6 | not reached: `LK-ﾈﾝ = 0` N, `LK-ﾍﾝｶﾝｽﾞﾐ` Y, `WS-MM < 1 OR WS-MM > 12` N, `WS-DD < 1 OR WS-DD > 31` N, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX)` N, `ｹﾞﾝｺﾞｳ-IX < ｹﾞﾝｺﾞｳ-ｹﾝｽｳ` Y, `WS-YYYYMMDD < ｹﾞﾝｺﾞｳ-ｶｲｼﾋﾞ(ｹﾞﾝｺﾞｳ-IX), again` Y |
| 7 | [case 14](14/) |

## GLC030

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 21](21/) |
| 2 | [case 19](19/) |
| 3 | not reached: `WS-ﾖﾐｺﾐｽﾞﾐ` Y |

### [ﾖﾐｺﾐ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾖﾐｺﾐ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 21](21/) |
| 2 | not reached: `WS-CSVTB-OK` Y, `WS-CSVTB-ｵﾜﾘ` N |
| 3 | [case 19](19/) |

### [ﾄｳﾛｸ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ﾄｳﾛｸ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | not reached: `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = CV-ｼｮｳﾋﾝ OR WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = ZEROS` N |

### [ｹｲｻﾝ-RTN](../../../blueprints/gondwana/rules/GLC030.md#ｹｲｻﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 2 | [case 29](29/) |
| 3 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` Y, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y |
| 4 | [case 21](21/) |
| 5 | [case 30](30/) |
| 6 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 7 | [case 43](43/) |
| 8 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` Y |
| 9 | [case 27](27/) |
| 10 | [case 44](44/) |
| 11 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` Y, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 12 | [case 45](45/) |
| 13 | not reached: `LK-ｶﾆｭｳﾈﾝﾚｲ < 20` N, `LK-ｹｲｶﾈﾝｽｳ > 40` N, `WS-ﾈﾝｽｳ > 0` Y, `WS-ｼｮｳﾋﾝ-ｺｰﾄﾞ(WS-ｼｮｳﾋﾝ-IX) = LK-ｼｮｳﾋﾝ` Y, `LK-ｶﾆｭｳﾈﾝﾚｲ < 60` N |
| 14 | [case 28](28/) |
| 15 | [case 46](46/) |

## GLC040

### [ﾒｲﾝ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾒｲﾝ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 14](14/) |
| 2 | not reached: `LK-ｻﾞﾝﾀﾞｶ > 0 AND LK-ｷｻﾝﾋﾞ NOT = SPACES` Y, `WS-ﾍﾝｶﾝｹｯｶ = 0` N |
| 3 | [case 42](42/) |

### [ﾋﾜﾘ-RTN](../../../blueprints/gondwana/rules/GLC040.md#ﾋﾜﾘ-rtn)

| Rule | Reached by |
| --: | --- |
| 1 | [case 18](18/) |
| 2 | [case 14](14/) |
