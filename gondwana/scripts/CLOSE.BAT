@echo off
rem CLOSE.BAT: 月次締め。月末の日次処理の後に、責任準備金の計算、
rem 代理店手数料の計算、そしてアンロード。
rem
rem 使い方: CLOSE.BAT YYYYMMDD  月末の業務日付
rem 環境: GLS_BIN プログラム、GLS_SCRIPTS バッチファイル、GLS_DATA データ
rem
rem 変更履歴
rem 1999-11-08 MS  GL-0233 初版、PC システム移行時 (旧 JCL から)。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥CLOSE.%DATE%.LOG
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_PRODTB=%DATA%¥MASTER¥PRODTB.TXT
set DD_GLJRNL=%DATA%¥OUT¥GLJRNL.%DATE%.TXT

rem STEP GLM020: 責任準備金計算
set DD_RSVTB=%DATA%¥MASTER¥RSVTB.TXT
set DD_GLM020R=%DATA%¥RPT¥GLM020.%DATE%.RPT
"%BIN%¥GLM020.EXE" %DATE% > "%DATA%¥RPT¥GLM020.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLM020 GLM020 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%

rem STEP GLM030: 代理店手数料計算
set DD_AGTMST=%DATA%¥MASTER¥AGTMST.TXT
set DD_COMMIS=%DATA%¥OUT¥COMMIS.%DATE%.TXT
set DD_GLM030R=%DATA%¥RPT¥GLM030.%DATE%.RPT
"%BIN%¥GLM030.EXE" %DATE% > "%DATA%¥RPT¥GLM030.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLM030 GLM030 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%

rem STEP UNLOAD: 月末のアンロード
call "%GLS_SCRIPTS%¥UNLOAD.BAT" %DATE%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
