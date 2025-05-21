# SQLParser

A MySQL SQL parsing library for Go, derived from and extending the Vitess SQL parser.

## Overview

SQLParser is a Go library that provides robust SQL parsing capabilities for MySQL syntax. This project is based on the SQL parsing components of the [Vitess](https://github.com/vitessio/vitess) project, with specific modifications and enhancements to better serve as a standalone parsing library.

Key modifications include:
- Renamed `HASH` function to `HASHFUNC` to avoid naming conflicts
- Added support for `HASH` and `PARALLEL` join types
- Modified syntax to improve parsing flexibility and functionality
- Streamlined for use as an importable Go library

## Installation

```bash
go get github.com/vedadiyan/sqlparser
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/vedadiyan/sqlparser"
)

func main() {
    sql := "SELECT * FROM users JOIN orders ON users.id = orders.user_id USING HASH"
    
    stmt, err := sqlparser.Parse(sql)
    if err != nil {
        fmt.Printf("Error parsing SQL: %v\n", err)
        return
    }
    
    fmt.Printf("Parsed statement: %v\n", stmt)
}
```

## Features

- Complete MySQL syntax parsing
- Support for `HASH` and `PARALLEL` join types
- AST (Abstract Syntax Tree) generation for SQL statements
- Thread-safe and efficient parsing

## License

This project is derived from the [Vitess](https://github.com/vitessio/vitess) project, which is licensed under the Apache License 2.0.

```
Copyright 2021 The Vitess Authors.
Copyright 2025 Pouya Vedadiyan.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Acknowledgments

This project builds upon the excellent work of the Vitess team and community. The original SQL parsing code from Vitess has been modified and adapted for this library's specific purposes.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request