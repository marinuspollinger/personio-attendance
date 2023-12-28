module personio-cli

go 1.21.4

// a fork of zap which has trace-level
replace go.uber.org/zap v1.24.0-traceLevel.1 => github.com/dunv/zap v1.24.0-traceLevel.1

require (
	github.com/codingconcepts/env v0.0.0-20200821220118-a8fbf8d84482
	github.com/dunv/ulog/v2 v2.0.10
)

require (
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0-traceLevel.1 // indirect
)
