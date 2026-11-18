@echo off
rem CERT.BAT: 保険料控除証明書。年に一度、10月の最初の営業日の日次
rem 処理の後に、その年に収納した保険料の証明書を作り、印刷会社に渡す。
rem
rem 使い方: CERT.BAT YYYYMMDD  業務日付。その年の1月からこの日まで。
rem 環境: GLS_BIN プログラム、GLS_DATA データ
rem
rem 変更履歴
rem 2026-11-19 YN  GL-0827 初版。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥CERT.%DATE%.LOG

rem STEP GLY030: 保険料控除証明書
set DD_PRMHST=%DATA%¥MASTER¥PRMHST.IDX
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_PRODTB=%DATA%¥MASTER¥PRODTB.TXT
set DD_PRMCERT=%DATA%¥OUT¥PRMCERT.%DATE%.TXT
set DD_GLY030R=%DATA%¥RPT¥GLY030.%DATE%.RPT
"%BIN%¥GLY030.EXE" %DATE% > "%DATA%¥RPT¥GLY030.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLY030 GLY030 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
