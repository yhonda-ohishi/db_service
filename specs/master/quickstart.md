# クイックスタートガイド

## 前提条件

- Go 1.21以上がインストール済み
- MySQL/MariaDBが起動済み
- データベース`db1`が作成済み

## セットアップ

### 1. 環境変数の設定

`.env.example`を`.env`にコピーして編集：

```bash
cp .env.example .env
```

`.env`ファイルを編集：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=db1
GRPC_PORT=50051
```

### 2. 依存関係のインストール

```bash
go mod download
```

### 3. Protocol Buffersのコンパイル

```bash
# protoc-gen-goとprotoc-gen-go-grpcをインストール
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# protoファイルをコンパイル
protoc --go_out=. --go-grpc_out=. src/proto/ryohi.proto
```

### 4. サーバーの起動

```bash
go run cmd/server/main.go
```

## 動作確認

### gRPCクライアントでの接続テスト

```bash
# grpcurlをインストール（未インストールの場合）
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# サービス一覧の確認
grpcurl -plaintext localhost:50051 list

# DTakoUriageKeihiの一覧取得
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:50051 ryohi.DTakoUriageKeihiService/List
```

### サンプルリクエスト

#### 経費データの作成

```json
{
  "dtako_uriage_keihi": {
    "srch_id": "TEST001",
    "datetime": "2025-09-19T10:00:00Z",
    "keihi_c": 1,
    "price": 1000.0,
    "km": 50.5,
    "dtako_row_id": "DTAKO001",
    "dtako_row_id_r": "DTAKO001R"
  }
}
```

#### ETC明細の作成

```json
{
  "etc_meisai": {
    "date_to": "2025-09-19T10:00:00Z",
    "date_to_date": "2025-09-19",
    "ic_fr": "東京IC",
    "ic_to": "横浜IC",
    "price": 1500,
    "shashu": 1,
    "etc_num": "1234567890123456"
  }
}
```

## トラブルシューティング

### データベース接続エラー

1. `.env`ファイルの設定を確認
2. MySQLが起動しているか確認
3. ユーザー権限を確認

```sql
GRANT ALL PRIVILEGES ON db1.* TO 'your_username'@'localhost';
FLUSH PRIVILEGES;
```

### gRPC接続エラー

1. ポート50051が使用されていないか確認
```bash
netstat -an | grep 50051
```

2. ファイアウォールの設定を確認

### コンパイルエラー

1. Go バージョンの確認
```bash
go version
```

2. 依存関係の再インストール
```bash
go mod tidy
go mod download
```

## テストの実行

```bash
# 単体テスト
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 統合テスト（要データベース）
go test -tags=integration ./...
```

## 次のステップ

1. [データモデル仕様書](./data-model.md)を確認
2. [API仕様書](./contracts/)を確認
3. ビジネスロジックの実装
4. 監視・ロギングの設定