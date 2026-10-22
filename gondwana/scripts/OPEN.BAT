@echo off
rem OPEN.BAT: システム開設。契約抽出ファイルから契約マスタの索引
rem ファイルと制御レコードを作る。一度だけ実行する。
rem
rem 使い方: OPEN.BAT YYYYMMDD  業務日付
rem 環境: GLS_BIN プログラム、GLS_SCRIPTS バッチファイル、GLS_DATA データ
rem
rem 変更履歴
rem 1999-11-08 MS  GL-0233 初版、PC システム移行時。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥OPEN.LOG
if not exist "%DATA%¥MASTER" mkdir "%DATA%¥MASTER"
if not exist "%DATA%¥OUT" mkdir "%DATA%¥OUT"
if not exist "%DATA%¥RPT" mkdir "%DATA%¥RPT"
if not exist "%DATA%¥LOG" mkdir "%DATA%¥LOG"

rem STEP GLO010: 契約マスタの初期登録と制御レコード
set DD_EXTRACT=%DATA%¥IN¥EXTRACT.TXT
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_GLO010R=%DATA%¥RPT¥GLO010.%DATE%.RPT
"%BIN%¥GLO010.EXE" %DATE% > "%DATA%¥RPT¥GLO010.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLO010 GLO010 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
