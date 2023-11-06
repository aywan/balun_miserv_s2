module github.com/aywan/balun_miserv_s2/shared/lib/db

go 1.21.3

require (
	github.com/georgysavva/scany/v2 v2.0.0
	github.com/jackc/pgx-zap v0.0.0-20221202020421-94b1cb2f889f
	github.com/jackc/pgx/v5 v5.5.0
	go.uber.org/zap v1.26.0
	github.com/aywan/balun_miserv_s2/shared/lib/logger v0.0.0
)

replace (
	"github.com/aywan/balun_miserv_s2/shared/lib/logger" => "../logger"
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/text v0.11.0 // indirect
)
