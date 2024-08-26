# API

## Sample directory tree(DDD)

```txt
myproject/
├── cmd/
│   └── myproject/           # エントリーポイント、main関数を含む
│       └── main.go
├── internal/
│   ├── app/                 # アプリケーション層（ユースケース）
│   │   ├── service.go
│   │   └── interfaces.go
│   ├── domain/              # ドメイン層
│   │   ├── model/           # ドメインモデル（アグリゲート、エンティティ、値オブジェクト）
│   │   │   ├── order.go     # アグリゲートルート例: Order
│   │   │   ├── customer.go  # エンティティ例: Customer
│   │   │   └── product.go   # 値オブジェクト例: Product
│   │   ├── payment/         # 新しいドメイン層: Payment
│   │   │   ├── model/       # ドメインモデル（アグリゲート、エンティティ、値オブジェクト）
│   │   │   │   ├── payment.go # アグリゲートルート例: Payment
│   │   │   │   └── transaction.go # エンティティ例: Transaction
│   │   │   ├── repository.go  # リポジトリインターフェース
│   │   │   ├── service.go     # ドメインサービス
│   │   │   └── event.go       # ドメインイベント
│   ├── infrastructure/      # インフラ層
│   │   ├── persistence/     # データベース関連コード
│   │   │   ├── mysql/
│   │   │   │   ├── order_repository.go
│   │   │   │   └── payment_repository.go  # 新しいリポジトリ実装
│   │   └── api/             # 外部APIクライアント
│   │       └── payment_client.go
│   └── interfaces/          # インターフェースアダプタ層
│       ├── api/             # Web API（HTTPハンドラ）
│       │   ├── handlers.go
│       │   ├── payment_handlers.go # 新しいドメインに対応するハンドラ
│       │   └── routers.go
│       └── repository/      # リポジトリの実装
│           ├── order_repository.go
│           └── payment_repository.go # 新しいリポジトリ実装
└── pkg/                     # 外部公開するパッケージ
    └── myprojectlib/
        └── utils.go
```

### Packages

#### echo

#### gorm gen

<https://gorm.io/gen/query.html>

Generate Query

```bash
go run cmd/lib/gen.go
```
