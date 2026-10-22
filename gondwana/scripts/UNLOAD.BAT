@echo off
rem UNLOAD.BAT: 契約マスタのアンロード。日次と月次の最後に呼ばれ、
rem 下流に渡すファイルを作り、制御レコードの業務日付を進める。
rem
rem 使い方: UNLOAD.BAT YYYYMMDD  業務日付
rem 環境: GLS_BIN プログラム、GLS_DATA データ
rem
rem 変更履歴
rem 1999-11-08 MS  GL-0233 初版、PC システム移行時。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥UNLOAD.%DATE%.LOG

rem STEP GLD090: アンロードと制御レコード
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_UNLOAD=%DATA%¥OUT¥POLMST.%DATE%.TXT
set DD_GLD090R=%DATA%¥RPT¥GLD090.%DATE%.RPT
"%BIN%¥GLD090.EXE" %DATE% > "%DATA%¥RPT¥GLD090.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLD090 GLD090 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
