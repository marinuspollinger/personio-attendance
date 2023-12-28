module personio-server

go 1.21.4

require (
	github.com/codingconcepts/env v0.0.0-20200821220118-a8fbf8d84482
	github.com/dunv/uhttp v1.2.8
	github.com/dunv/ulog/v2 v2.0.10
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dunv/uhelpers v1.0.17 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0-traceLevel.1 // indirect
	golang.org/x/sys v0.9.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

// a fork of zap which has trace-level
replace go.uber.org/zap v1.24.0-traceLevel.1 => github.com/dunv/zap v1.24.0-traceLevel.1
