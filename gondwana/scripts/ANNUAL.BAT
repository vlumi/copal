@echo off
rem ANNUAL.BAT: 年次処理。年度末 (3月31日) の月次締めの後に、契約者
rem 配当の割り当て。
rem
rem 使い方: ANNUAL.BAT YYYYMMDD  年度末の業務日付
rem 環境: GLS_BIN プログラム、GLS_DATA データ
rem
rem 変更履歴
rem 1999-11-08 MS  GL-0233 初版、PC システム移行時 (旧 JCL から)。
setlocal
set DATE=%1
set BIN=%GLS_BIN%
set DATA=%GLS_DATA%
set LOG=%DATA%¥LOG¥ANNUAL.%DATE%.LOG

rem STEP GLY010: 契約者配当
set DD_POLMST=%DATA%¥MASTER¥POLMST.IDX
set DD_GLSCTL=%DATA%¥MASTER¥GLSCTL.BIN
set DD_PRODTB=%DATA%¥MASTER¥PRODTB.TXT
set DD_GLJRNL=%DATA%¥OUT¥GLJRNL.%DATE%.TXT
set DD_GLY010R=%DATA%¥RPT¥GLY010.%DATE%.RPT
"%BIN%¥GLY010.EXE" %DATE% > "%DATA%¥RPT¥GLY010.%DATE%.OUT" 2>&1
>>"%LOG%" echo GLY010 GLY010 %ERRORLEVEL%
if errorlevel 8 exit /b %ERRORLEVEL%
exit /b 0
