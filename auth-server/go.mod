module github.com/aywan/balun_miserv_s2/auth-server

go 1.21.3

require (
	github.com/Masterminds/squirrel v1.5.4
	github.com/aywan/balun_miserv_s2/shared/lib/db v0.0.0
	github.com/aywan/balun_miserv_s2/shared/lib/logger v0.0.0
	github.com/aywan/balun_miserv_s2/shared/lib/runutil v0.0.0
	github.com/brianvoe/gofakeit/v6 v6.24.0
	github.com/fatih/color v1.15.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/fx v1.20.1
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.13.0
	google.golang.org/grpc v1.58.2
	google.golang.org/protobuf v1.31.0
)

replace (
	github.com/aywan/balun_miserv_s2/shared/lib/db => ./../shared/lib/db
	github.com/aywan/balun_miserv_s2/shared/lib/logger => ./../shared/lib/logger
	github.com/aywan/balun_miserv_s2/shared/lib/runutil => ./../shared/lib/runutil
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/georgysavva/scany/v2 v2.0.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx-zap v0.0.0-20221202020421-94b1cb2f889f // indirect
	github.com/jackc/pgx/v5 v5.5.0 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/net v0.15.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
