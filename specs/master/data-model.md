# データモデル仕様書

## エンティティ定義

### 1. DTakoUriageKeihi（経費精算データ）

**テーブル名**: `dtako_uriage_keihi`

**フィールド定義**:
| フィールド名 | 型 | 制約 | 説明 |
|------------|-----|------|------|
| srch_id | VARCHAR(44) | PRIMARY KEY | 検索ID |
| datetime | DATETIME | PRIMARY KEY | 日時 |
| keihi_c | INT(11) | PRIMARY KEY | 経費コード |
| price | DOUBLE | NOT NULL | 金額 |
| km | DOUBLE | NULL | 距離 |
| dtako_row_id | VARCHAR(24) | NOT NULL | 運行NO |
| dtako_row_id_r | VARCHAR(23) | NOT NULL | 運行NO（参照） |
| start_srch_id | VARCHAR(44) | NULL | 開始検索ID |
| start_srch_time | DATETIME | NULL | 開始時刻 |
| start_srch_place | VARCHAR(50) | NULL | 開始場所 |
| start_srch_tokui | VARCHAR(9) | NULL | 開始得意先 |
| end_srch_id | VARCHAR(44) | NULL | 終了検索ID |
| end_srch_time | DATETIME | NULL | 終了時刻 |
| end_srch_place | VARCHAR(50) | NULL | 終了場所 |
| manual | TINYINT(1) | NULL | 手動フラグ |

**インデックス**:
- PRIMARY KEY: (srch_id, datetime, keihi_c)
- INDEX: dtako_row_id, srch_id
- INDEX: start_srch_id
- INDEX: end_srch_id

**バリデーション**:
- srch_id: 必須、最大44文字
- price: 必須、0以上
- dtako_row_id: 必須、最大24文字

### 2. ETCMeisai（ETC明細データ）

**テーブル名**: `etc_meisai`

**フィールド定義**:
| フィールド名 | 型 | 制約 | 説明 |
|------------|-----|------|------|
| id | BIGINT(20) | PRIMARY KEY, AUTO_INCREMENT | ID |
| date_fr | DATETIME | NULL | 開始日時 |
| date_to | DATETIME | NOT NULL | 終了日時 |
| date_to_date | DATE | NOT NULL | 終了日 |
| IC_fr | VARCHAR(30) | NOT NULL | 入口IC |
| IC_to | VARCHAR(30) | NOT NULL | 出口IC |
| price_bf | INT(11) | NULL | 割引前料金 |
| descount | INT(11) | NULL | 割引額 |
| price | INT(11) | NOT NULL | 料金 |
| shashu | INT(11) | NOT NULL | 車種 |
| car_id_num | INT(11) | NULL | 車両ID番号 |
| etc_num | VARCHAR(20) | NOT NULL | ETCカード番号 |
| detail | VARCHAR(40) | NULL | 詳細 |
| dtako_row_id | VARCHAR(24) | NULL | 運行NO |

**インデックス**:
- PRIMARY KEY: (id)
- INDEX: date_to, id
- INDEX: dtako_row_id, id

**バリデーション**:
- date_to: 必須
- IC_fr, IC_to: 必須、最大30文字
- price: 必須、0以上
- etc_num: 必須、最大20文字

### 3. DTakoFerryRows（フェリー運行データ）

**テーブル名**: `dtako_ferry_rows`

**フィールド定義**:
| フィールド名 | 型 | 制約 | 説明 |
|------------|-----|------|------|
| id | INT(11) | PRIMARY KEY, AUTO_INCREMENT | ID |
| 運行NO | VARCHAR(23) | NOT NULL | 運行番号 |
| 運行日 | DATE | NOT NULL | 運行日 |
| 読取日 | DATE | NOT NULL | 読取日 |
| 事業所CD | INT(11) | NOT NULL | 事業所コード |
| 事業所名 | VARCHAR(20) | NOT NULL | 事業所名 |
| 車輌CD | INT(11) | NOT NULL | 車両コード |
| 車輌名 | VARCHAR(20) | NOT NULL | 車両名 |
| 乗務員CD1 | INT(11) | NOT NULL | 乗務員コード1 |
| 乗務員名１ | VARCHAR(20) | NOT NULL | 乗務員名1 |
| 対象乗務員区分 | INT(11) | NOT NULL | 対象乗務員区分 |
| 開始日時 | DATETIME | NOT NULL | 開始日時 |
| 終了日時 | DATETIME | NOT NULL | 終了日時 |
| フェリー会社CD | INT(11) | NOT NULL | フェリー会社コード |
| フェリー会社名 | VARCHAR(20) | NOT NULL | フェリー会社名 |
| 乗場CD | INT(11) | NOT NULL | 乗場コード |
| 乗場名 | VARCHAR(20) | NOT NULL | 乗場名 |
| 便 | VARCHAR(10) | NOT NULL | 便 |
| 降場CD | INT(11) | NOT NULL | 降場コード |
| 降場名 | VARCHAR(20) | NOT NULL | 降場名 |
| 精算区分 | INT(11) | NOT NULL | 精算区分 |
| 精算区分名 | VARCHAR(20) | NOT NULL | 精算区分名 |
| 標準料金 | INT(11) | NOT NULL | 標準料金 |
| 契約料金 | INT(11) | NOT NULL | 契約料金 |
| 航送車種区分 | INT(11) | NOT NULL | 航送車種区分 |
| 航送車種区分名 | VARCHAR(20) | NOT NULL | 航送車種区分名 |
| 見なし距離 | INT(11) | NOT NULL | 見なし距離 |
| ferry_srch | VARCHAR(60) | NULL | フェリー検索 |

**インデックス**:
- PRIMARY KEY: (id)

**バリデーション**:
- 運行NO: 必須、最大23文字
- 各日付フィールド: 必須、有効な日付
- 各名称フィールド: 必須、最大20文字
- 料金フィールド: 0以上

## リレーションシップ

### DTakoUriageKeihi ↔ DTakoRows
- dtako_row_idフィールドによる関連
- 1対多の関係（1つの運行に複数の経費）

### ETCMeisai ↔ DTakoRows
- dtako_row_idフィールドによる関連
- 1対多の関係（1つの運行に複数のETC明細）

## 状態遷移

### DTakoUriageKeihi
- 作成 → 更新可能 → 削除可能
- manualフラグによる手動/自動の区別

### ETCMeisai
- 作成（自動インクリメント） → 更新可能 → 削除可能

### DTakoFerryRows
- 作成（自動インクリメント） → 更新可能 → 削除可能

## トランザクション要件

1. **複数レコード一括登録**
   - 同一運行NOの複数経費を一括登録
   - 失敗時は全ロールバック

2. **更新時の整合性**
   - 関連データの同時更新
   - 部分的な更新の防止

3. **削除時のカスケード**
   - 関連データの確認
   - 依存関係の検証