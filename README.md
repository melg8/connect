<!--
SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>

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