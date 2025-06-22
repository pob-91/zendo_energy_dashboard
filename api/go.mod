module zendo/api

go 1.24.4

replace zendo/lib_zendo => ../lib_zendo

require (
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.27.0
	zendo/lib_zendo v0.0.0-00010101000000-000000000000
)

require go.uber.org/multierr v1.10.0 // indirect
