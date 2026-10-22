      *----------------------------------------------------------------
      * GLC010: 保険年齢
      *
      * 生年月日と基準日 (ともに YYYYMMDD)、年齢を受けて呼ばれる。
      * 保険年齢は満年齢に、次の1年のうち6か月以上が過ぎていれば
      * 1 を足したもの: 39歳7か月の人は40歳として加入する。
      * 月数は満月数で、誕生日の日より前の日は数えない。
      *
      * 変更履歴
      * 日付       担当 変更
      * 1987-04-01 KT   GL-0001 初版。
      *----------------------------------------------------------------
       IDENTIFICATION DIVISION.
       PROGRAM-ID. GLC010.
       AUTHOR. ゴンドワナ生命 システム部.
       DATE-WRITTEN. 1987-04-01.

       DATA DIVISION.
       WORKING-STORAGE SECTION.
       01  WS-ｾｲﾈﾝｶﾞｯﾋﾟ.
           05  WS-B-YYYY               PIC 9(4).
           05  WS-B-MM                 PIC 99.
           05  WS-B-DD                 PIC 99.
       01  WS-ｷｼﾞｭﾝﾋﾞ.
           05  WS-A-YYYY               PIC 9(4).
           05  WS-A-MM                 PIC 99.
           05  WS-A-DD                 PIC 99.
       01  WS-ﾂｷｽｳ                     PIC S9(5).
       01  WS-ｱﾏﾘ                      PIC 99.

       LINKAGE SECTION.
       01  LK-ｾｲﾈﾝｶﾞｯﾋﾟ                PIC 9(8).
       01  LK-ｷｼﾞｭﾝﾋﾞ                  PIC 9(8).
       01  LK-ﾈﾝﾚｲ                     PIC 99.

       PROCEDURE DIVISION USING LK-ｾｲﾈﾝｶﾞｯﾋﾟ LK-ｷｼﾞｭﾝﾋﾞ LK-ﾈﾝﾚｲ.
      *----------------------------------------------------------------
      * ﾒｲﾝ-RTN: 満月数を求め、6か月で切り上げて年にする。
      *----------------------------------------------------------------
       ﾒｲﾝ-RTN.
           MOVE LK-ｾｲﾈﾝｶﾞｯﾋﾟ TO WS-ｾｲﾈﾝｶﾞｯﾋﾟ
           MOVE LK-ｷｼﾞｭﾝﾋﾞ TO WS-ｷｼﾞｭﾝﾋﾞ
           COMPUTE WS-ﾂｷｽｳ = (WS-A-YYYY - WS-B-YYYY) * 12
                           + WS-A-MM - WS-B-MM
           IF WS-A-DD < WS-B-DD
               SUBTRACT 1 FROM WS-ﾂｷｽｳ
           END-IF
           IF WS-ﾂｷｽｳ < 0
               MOVE 0 TO WS-ﾂｷｽｳ
           END-IF
           DIVIDE WS-ﾂｷｽｳ BY 12 GIVING LK-ﾈﾝﾚｲ
               REMAINDER WS-ｱﾏﾘ
           IF WS-ｱﾏﾘ NOT < 6
               ADD 1 TO LK-ﾈﾝﾚｲ
           END-IF
           GOBACK.
