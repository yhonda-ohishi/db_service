# Database-repo for ryohi_sub_cal

gRPCベースのデータベースリポジトリサービス

## 概要

ryohi_sub_calデータベースの3つのテーブル（dtako_uriage_keihi、etc_meisai、dtako_ferry_rows）に対するCRUD操作を提供するgRPCサービス。

## 技術スタック

- **言語**: Go 1.21+
- **フレームワーク**:
  - gRPC v1.75+
  - GORM v1.25.5
  - Protocol Buffers
- **データベース**:
  - MySQL/MariaDB (本番DB、読み取り専用)
  - SQL Server (CAPE#01データベース)
- **設定管理**: godotenv

## プロジェクト構造

```
db_service/
├── src/
│   ├── proto/       # Protocol Buffers定義とコンパイル済みファイル
│   ├── models/      # GORMモデル定義
│   ├── repository/  # データアクセス層
│   ├── service/     # gRPCサービス実装
│   ├── config/      # 設定管理
│   └── registry/    # サービス登録
├── sql_server_tables/ # SQL Serverテーブル定義（UTF-8）
├── tests/
│   ├── contract/    # 契約テスト
│   └── integration/ # 統合テスト
├── cmd/
│   └── server/      # サーバーエントリポイント
└── Makefile         # ビルドコマンド
```

## セットアップ

### 1. 前提条件

- Go 1.21以上
- Protocol Buffers Compiler (protoc)
- MySQL/MariaDB

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. 環境変数の設定

`.env`ファイルを作成して以下を設定：

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password  # ハードコード禁止
DB_NAME=db1
GRPC_PORT=50051
```

### 4. Protocol Buffersのコンパイル

```bash
make proto
# または
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       src/proto/ryohi.proto
```

## 実行方法

### サーバー起動

```bash
# ビルドして実行
make build && ./bin/server

# または直接実行
go run cmd/server/main.go

# またはMakefile経由
make run
```

### テスト実行

```bash
# 全テスト
make test

# カバレッジ付き
make test-coverage

# 統合テストのみ
make test-integration
```

## 他のプロジェクトからdb_serviceを利用する

### インストール

```bash
go get github.com/yhonda-ohishi/db_service
```

### 同一プロセスで統合（推奨）

db_serviceを他のgRPCサーバーと同じプロセス・ポートで統合できます:

```go
import (
    "google.golang.org/grpc"
    "github.com/yhonda-ohishi/db_service/src/registry"
)

func main() {
    // 既存のgRPCサーバーを作成
    grpcServer := grpc.NewServer()

    // 自分のサービスを登録
    pb.RegisterYourServiceServer(grpcServer, yourService)

    // db_serviceのサービスを自動登録（1行で完了！）
    registry.Register(grpcServer)

    // サーバーを起動
    listener, _ := net.Listen("tcp", ":50051")
    grpcServer.Serve(listener)
}
```

### 新しいサービスの追加時の自動対応

db_serviceに新しいサービスが追加された場合:

1. `src/registry/registry.go`の`ServiceRegistry`構造体に新しいサービスを追加
2. `NewServiceRegistry()`で新しいリポジトリとサービスを初期化
3. `RegisterAll()`で新しいサービスを登録

これにより、他のプロジェクト（desktop-server等）では**コード変更不要**で新しいサービスが利用可能になります。

## API仕様

### DTakoUriageKeihiService

経費精算データ管理（複合主キー: srch_id, datetime, keihi_c）

- `Create`: 新規作成
- `Get`: 取得
- `Update`: 更新
- `Delete`: 削除
- `List`: 一覧取得

### ETCMeisaiService

ETC明細データ管理（主キー: id AUTO_INCREMENT）

- `Create`: 新規作成
- `Get`: 取得
- `Update`: 更新
- `Delete`: 削除
- `List`: 一覧取得

### DTakoFerryRowsService

フェリー運行データ管理（主キー: id AUTO_INCREMENT）

- `Create`: 新規作成
- `Get`: 取得
- `Update`: 更新
- `Delete`: 削除
- `List`: 一覧取得

## 動作確認

### grpcurlを使用

```bash
# サービス一覧
grpcurl -plaintext localhost:50051 list

# データ取得例
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:50051 ryohi.DTakoUriageKeihiService/List
```

## セキュリティ

- **環境変数必須**: データベース認証情報は環境変数から取得
- **ハードコード禁止**: シークレット情報のハードコードは絶対禁止
- **.gitignore**: `.env`ファイルはバージョン管理対象外

## トラブルシューティング

### データベース接続エラー

1. `.env`ファイルの設定を確認
2. MySQLサービスが起動しているか確認
3. ユーザー権限を確認

### protocエラー

```bash
# Chocolateyでインストール
choco install protoc

# Go用プラグインインストール
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## ライセンス

(ライセンス情報を記載)

## 作成者

yhonda-ohishi