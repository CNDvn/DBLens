### Tổng quan kiến trúc

```mermaid
flowchart TD
    DB[(MySQL / InnoDB)]
    C[Collector Worker]
    RAW[(Layer 1: Raw Store)]
    N[Normalizer Worker]
    DOMAIN[(Layer 2: Normalized Store)]
    A[Lock Analyzer]
    EVENTS[(Lock Events)]
    API[REST API]
    UI[Web Dashboard]
    METRICS[Prometheus Metrics]
    G[Grafana]

    DB --> C
    C --> RAW
    RAW --> N
    N --> DOMAIN
    DOMAIN --> A
    A --> EVENTS
    DOMAIN --> API
    EVENTS --> API
    API --> UI
    C --> METRICS
    N --> METRICS
    A --> METRICS
    METRICS --> G
```

#### Migration Tutorial
##### Cài đặt CLI qua Go
- go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
##### Tạo các file migration
- migrate create -ext sql -dir db/migrations -seq create_users_table
##### Run migration
- Cách 1: chạy bằng CLI 
    - Up: `migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/mydb?sslmode=disable" up`
    - Rollback: `migrate -path db/migrations -database "postgres://postgres:secretpassword@localhost:5432/mydb?sslmode=disable" down 1`
- cách 2: chạy qua Go code: 
    ```Go
    package main

    import (
        "errors"
        "log"

        "github.com/golang-migrate/migrate/v4"
        _ "github.com/golang-migrate/migrate/v4/database/postgres"
        _ "github.com/golang-migrate/migrate/v4/source/file"
    )

    func main() {
        // Đường dẫn tới thư mục chứa file migration và chuỗi kết nối DB
        m, err := migrate.New(
            "file://db/migrations",
            "postgres://postgres:secretpassword@localhost:5432/mydb?sslmode=disable",
        )
        if err != nil {
            log.Fatalf("Không thể khởi tạo migration: %v", err)
        }

        // Chạy migration lên phiên bản mới nhất
        if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
            log.Fatalf("Chạy migration thất bại: %v", err)
        }

        log.Println("Migration thành công!")
    }
    ```