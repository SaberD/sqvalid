# sqlvalid

SQL file validator for Go projects. Validates all `.sql` files recursively.

**Supported databases:**
- PostgreSQL (full SQL parsing with pganalyze)
- SQLite (actual SQLite parser)

## Features

- ✓ PostgreSQL SQL parsing (full validation with pganalyze)
- ✓ SQLite SQL validation (actual SQLite parser)
- ✓ Recursive directory scanning
- ✓ Clear error reporting with file paths
- ✓ Exit code for CI/CD integration (0 = success, 1 = errors)
- ✓ Works with modular project structures
- ✓ Fast validation

## Installation

```bash
go install github.com/saberd/sqlvalid/cmd/sqlvalid@latest
```

Or a specific version:

```bash
go install github.com/saberd/sqlvalid/cmd/sqlvalid@v1.0.0
```

## Usage

### PostgreSQL (default)

```bash
sqlvalid ./sql
```

### SQLite

```bash
sqlvalid -sqlite ./sql
```

### Install Globally

If you prefer installing once:

```bash
go install github.com/saberd/sqlvalid@v1.0.0
```

Then use in makefile or scripts:

```makefile
.PHONY: validate
validate:
	sqlvalid ./sql              # PostgreSQL
	# OR
	sqlvalid -sqlite ./sql      # SQLite
```

### Makefile Integration

```makefile
.PHONY: validate
validate:
	@go run github.com/saberd/sqlvalid@v1.0.0 ./sql

.PHONY: build
build: validate
	CGO_ENABLED=1 go build -o myapp
```

### Database Flags

- Default: PostgreSQL validation
- `-sqlite`: Validate SQLite SQL

### Output

**When valid:**
```
✓ posts/sql/getPosts.sql
✓ posts/sql/createPost.sql
✓ posts/sql/updatePost.sql
✓ users/sql/getUser.sql
✓ users/sql/createUser.sql
✓ comments/sql/getComments.sql
✓ comments/sql/createComment.sql
...
✓ All 15 SQL files validated
```

**When invalid:**
```
✓ posts/sql/getPosts.sql
✓ posts/sql/createPost.sql
✓ posts/sql/updatePost.sql
✓ users/sql/getUser.sql
✓ users/sql/updateUser.sql
✓ comments/sql/getComments.sql

✗ users/sql/createUser.sql: invalid SQL
  syntax error at or near "INSRT"
✗ comments/sql/createComment.sql: invalid SQL
  syntax error at or near "VALES"

✗ 2 SQL errors found
exit status 1
```

## CI/CD Integration

### GitHub Actions

Add to `.github/workflows/build.yml`:

```yaml
- name: Validate SQL
  run: |
    go install github.com/saberd/sqlvalid@v1.0.0
    sqlvalid ./sql
```

Or with versioning in makefile:

```yaml
- name: Build
  run: make go
  # SQL validation happens automatically
```

## Why This Tool?

- **Modular projects**: Works with any directory structure
- **Build-time validation**: Catch SQL errors before deployment
- **No configuration**: Just point to directory
- **PostgreSQL & SQLite**: Works with both databases
- **Simple**: Single responsibility - validate SQL files
- **AI-Friendly**: Raw SQL files are easier for AI agents to understand and modify

## Perfect for AI Agents

pgvalid enables a workflow that's **ideal for AI code generation and modification:**

### Raw SQL Files

AI agents work better with:
- ✓ Explicit SQL files (not hidden in ORMs)
- ✓ Clear file organization (one query per file)
- ✓ Readable, modifiable SQL
- ✓ No generated boilerplate code

### Example AI Workflow

```
1. AI reads: users/sql/getUser.sql
2. AI understands the exact query
3. AI can modify or extend it safely
4. pgvalid validates the changes
5. Build passes ✓

Without pgvalid:
1. AI needs to understand ORM
2. AI needs to understand schema generation
3. AI needs to understand type mappings
4. Risk of type mismatches
5. Very complex
```

### Why ORMs Don't Work Well with AI

- ✗ Type generation is implicit
- ✗ Schema coupling is complex
- ✗ Hard to predict generated code
- ✗ Changes cascade unexpectedly
- ✗ AI needs domain knowledge of tool

### Why Raw SQL Works Well with AI

- ✓ SQL is standard and predictable
- ✓ Files are editable and readable
- ✓ Changes are explicit and traceable
- ✓ Validation is independent
- ✓ AI can reason about it directly

### Use Case: AI-Assisted Database Optimization

```
1. AI analyzes queries in users/sql/
2. AI suggests indexes, rewrites slow queries
3. AI modifies files directly
4. pgvalid validates all changes
5. Team reviews changes (clear diffs)
6. Merge with confidence
```

This would be **impossible with sqlc or ORMs** - you'd need AI to understand code generation, type mappings, and schema coupling.

## Modular SQL Architecture

This tool enables a clean, modular approach to SQL in Go projects:

```
posts/
├── sql/
│   ├── getPosts.sql
│   ├── createPost.sql
│   └── updatePost.sql
├── handlers.go
└── queries.go

users/
├── sql/
│   ├── getUser.sql
│   ├── createUser.sql
│   └── updateUser.sql
├── handlers.go
└── queries.go

comments/
├── sql/
│   ├── getComments.sql
│   └── createComment.sql
├── handlers.go
└── queries.go
```

Each module owns its SQL queries, keeping them organized and maintainable.

### Embedding SQL Files

Use Go's `embed` package to bundle SQL files into your binary:

```go
// posts/queries.go
package posts

import (
	"embed"
	"os"
	"path/filepath"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

//go:embed sql/*.sql
var sqlFS embed.FS

type Queries struct {
	// Queries loaded at startup
}

func LoadQueries() (*Queries, error) {
	// Validate all SQL files
	entries, err := sqlFS.ReadDir("sql")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := sqlFS.ReadFile(filepath.Join("sql", entry.Name()))
		if err != nil {
			return nil, err
		}

		// Validate with pganalyze
		_, err = pg_query.ParseToJSON(string(content))
		if err != nil {
			return nil, fmt.Errorf("invalid SQL in %s: %w", entry.Name(), err)
		}
	}

	return &Queries{}, nil
}
```

### Build-Time Validation

Add to your makefile:

```makefile
.PHONY: validate
validate:
	@go run github.com/saberd/sqlvalid@v1.0.0 ./sql

.PHONY: build
build: validate
	go build -o myapp
```

Now when you `make build`:
1. sqlvalid validates all `.sql` files
2. If any SQL is invalid → Build stops
3. If all valid → Binary is created

## How It Works

### PostgreSQL Mode (default)
1. Walks through all files in the given directory recursively
2. Finds all `.sql` files
3. Parses each with pganalyze's PostgreSQL parser
4. Reports errors with file paths and error messages
5. Exits with code 1 if any errors found, 0 if all valid

### SQLite Mode (`-db sqlite`)
1. Walks through all files in the given directory recursively
2. Finds all `.sql` files
3. Validates SQL syntax using SQLite's actual parser (mattn/go-sqlite3)
4. For DML statements (SELECT, INSERT, UPDATE, DELETE), creates a dummy table to avoid "table not found" errors
5. Reports syntax errors with file paths
6. Exits with code 1 if any errors found, 0 if all valid

### Why This Approach?

**Pragmatic and Frugal:**
- ✓ PostgreSQL: Use proven pganalyze parser (PostgreSQL's own parser)
- ✓ SQLite: Use actual SQLite parser (you're already using CGO if you use SQLite in Go)
- ✓ No complex regex or manual validation
- ✓ Validates what actually matters: SQL syntax correctness
- ✓ Simple, maintainable code

## Development

```bash
git clone https://github.com/saberd/sqlvalid.git
cd sqlvalid

# Build locally
go build -o sqlvalid

# Test
./sqlvalid ./testdata
./sqlvalid -sqlite ./testdata/sqlite/valid
```

## License

MIT

## Contributing

Pull requests welcome! This is a simple tool, so scope changes to SQL validation only.

## Examples

See [examples/](examples/) for sample SQL projects.

## Support

Found a bug or have a feature request? [Open an issue](https://github.com/saberd/sqlvalid/issues).
