# Ryohi Database Service - Protocol Buffers定義

このディレクトリには、ryohi_sub_calデータベース向けのgRPCサービス定義が含まれています。

## 概要

- **パッケージ名**: `db_service`
- **Go パッケージ**: `github.com/yhonda-ohishi/db_service/src/proto`
- **Protocol Buffers バージョン**: proto3

## 提供サービス

### SQL Serverテーブル（CAPE#01）

1. **DTakoUriageKeihiService** - 経費精算データ管理
2. **ETCMeisaiService** - ETC明細データ管理
3. **DTakoFerryRowsService** - フェリー運行データ管理
4. **ETCMeisaiMappingService** - ETC明細マッピング管理
5. **UntenNippoMeisaiService** - 運転日報明細管理
6. **ShainMasterService** - 社員マスタ管理
7. **ChiikiMasterService** - 地域マスタ管理
8. **ChikuMasterService** - 地区マスタ管理

### MySQLテーブル（本番DB、読み取り専用）

1. **DTakoCarsService** - 車輌マスタ管理
2. **DTakoEventsService** - イベント情報管理
3. **DTakoRowsService** - 運行データ管理
4. **ETCNumService** - ETCカード番号マスタ管理
5. **DTakoFerryRowsProdService** - フェリー運行データ管理
6. **CarsService** - 車両マスタ管理
7. **DriversService** - ドライバーマスタ管理

## 使用方法

### Go言語での使用

```go
import (
    pb "github.com/yhonda-ohishi/db_service/src/proto"
)
```

### protoファイルのインポート

他のプロジェクトからこのproto定義を使用する場合：

```protobuf
import "db_service.proto";
```

### Bufを使用したコード生成

このディレクトリには、Buf（buf.build）用の設定ファイルが含まれています：

- `buf.yaml` - Buf設定とlintルール
- `buf.gen.yaml` - コード生成プラグイン設定

#### コード生成コマンド

```bash
# src/protoディレクトリで実行
buf generate
```

生成されるファイル：
- `*.pb.go` - Protocol Buffersメッセージ定義
- `*_grpc.pb.go` - gRPCサービス定義
- `*.pb.gw.go` - gRPC-Gateway用リバースプロキシ
- `swagger/*.swagger.json` - OpenAPI/Swagger定義

## 依存関係

- `buf.build/googleapis/googleapis` - Google APIの共通定義（HTTP annotations等）

## ライセンス

プライベートリポジトリ - 内部使用のみ
