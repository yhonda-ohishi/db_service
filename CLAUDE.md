# Claude Code Context for Database-repo

## プロジェクト概要
ryohi_sub_cal データベースのgRPCベースリポジトリサービス

## 技術スタック
- **言語**: Go 1.21
- **フレームワーク**: gRPC v1.58.0, GORM v1.25.5
- **データベース**: MySQL/MariaDB (db1)
- **設定管理**: godotenv（環境変数）

## プロジェクト構造
```
src/
├── models/      # GORMモデル（DTakoUriageKeihi, ETCMeisai, DTakoFerryRows）
├── repository/  # データアクセス層
├── service/     # gRPCサービス実装
├── proto/       # Protocol Buffer定義
└── config/      # 環境設定管理
```

## 主要テーブル
1. **dtako_uriage_keihi**: 経費精算（複合主キー: srch_id, datetime, keihi_c）
2. **etc_meisai**: ETC明細（主キー: id AUTO_INCREMENT）
3. **dtako_ferry_rows**: フェリー運行（主キー: id AUTO_INCREMENT）

## 重要な制約
- シークレットはハードコードしない（環境変数から取得）
- 仕様書等は日本語で表記
- 単一責任原則：データアクセスのみ（ビジネスロジック禁止）

## 最近の変更
- 2025-09-19: 初期設定、憲章作成、Phase 0-1完了
- Protocol Buffers定義作成
- データモデル仕様書作成

## 次のタスク
- GORMモデル実装
- リポジトリ層実装
- gRPCサービス実装