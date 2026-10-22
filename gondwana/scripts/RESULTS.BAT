@echo off
rem RESULTS.BAT: 口座振替結果処理。収納代行から結果ファイルが届いた日
rem (振替日の数営業日後) に、日次処理の前に実行する。
rem
rem 使い方: RESULTS.BAT YYYYMMDD  業務日付。結果ファイルの日付でもある。
rem 環境: GLS_BIN プログラム、GLS_DATA データ
rem
rem 変更履歴
rem 2003-10-06 MS  GL-0310 初版、口座振替の開始。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥RESULTS.%DATE%.LOG

rem STEP GLD050: 振替結果の反映
set DD_PRMRES=%DATA%¥IN¥PRMRES.%DATE%.DAT
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_GLJRNL=%DATA%¥OUT¥GLJRNL.%DATE%.TXT
set DD_RESREJ=%DATA%¥OUT¥RESREJ.%DATE%.TXT
set DD_GLD050R=%DATA%¥RPT¥GLD050.%DATE%.RPT
"%BIN%¥GLD050.EXE" %DATE% > "%DATA%¥RPT¥GLD050.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLD050 GLD050 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
