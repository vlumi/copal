# BANKDATE against wrong.sh

wrong.sh conforms on 0 of 3 cases.

| Case | Differences |
| ---: | --- |
| 01 | CONTROL-OUT, record 1: CONTROL-OUT-RECORD "20261016", recorded "20261019" |
| 02 | return code 0, recorded 16; displayed nothing, recorded "BANKDATE: NO BANK CONTROL RECORD" |
| 03 | return code 0, recorded 16; displayed nothing, recorded "BANKDATE: CANNOT OPEN FILES 35 00" |
