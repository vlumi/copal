@echo off
rem BILLING.BAT: 保険料請求。月に一度、振替日 (27日) の前に、収納代行
rem に送る口座振替依頼ファイルを作る。
rem
rem 使い方: BILLING.BAT YYYYMMDD  実行日。その月が請求月。
rem 環境: GLS_BIN プログラム、GLS_DATA データ
rem
rem 変更履歴
rem 2003-10-06 MS  GL-0310 初版、口座振替の開始。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥BILLING.%DATE%.LOG

rem STEP GLM010: 口座振替依頼ファイルの作成
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_PRMREQ=%DATA%¥OUT¥PRMREQ.%DATE%.DAT
set DD_GLM010R=%DATA%¥RPT¥GLM010.%DATE%.RPT
"%BIN%¥GLM010.EXE" %DATE% > "%DATA%¥RPT¥GLM010.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLM010 GLM010 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
