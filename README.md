<!--
SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>

SPDX-License-Identifier: MIT
-->

# connect
Emulate multiple client connections to your game server.


# commands

Run main application
```bash
go run .\cmd\connect\main.go
```

Run unit tests of project:
```bash
go test ./... --cover --count=1
```

Build with compiler explanation of heap vs stack memory for variables:
```bash
go build -gcflags "-m=2"
```

Run linters for project:
```bash
golangci-lint run ./...
```

Run specific linter for project with fixes:
```bash
golangci-lint run --fix --disable-all --enable=wsl ./...
```

Run concrete benchmark with profiling:
```bash
E:\Go\bin\go.exe test -benchmem -cpuprofile=cpu_out -memprofile=mem_out  -run=^$ -bench ^BenchmarkEncryptor_Write$ github.com/melg8/connect/internal/connect/crypt
```

Run pprof tool for profiling:
```bash
go tool pprof -http=localhost:8080 mem_out
```

Run pprof tool for profiling:
```bash
go tool pprof -http=localhost:8080 cpu_out
```
