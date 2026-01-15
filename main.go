package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	pg_query "github.com/pganalyze/pg_query_go/v6"
	_ "github.com/mattn/go-sqlite3"
)

type sqlError struct {
	path string
	err  string
}

func main() {
	useSQLite := flag.Bool("sqlite", false, "Validate SQLite SQL (default: PostgreSQL)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Printf("Usage: sqlvalid [-sqlite] <directory>\n")
		os.Exit(1)
	}

	rootDir := args[0]

	validCount := 0
	var errors []sqlError

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			errors = append(errors, sqlError{path: path, err: "failed to read"})
			return nil
		}

		sqlStr := string(content)
		var validationErr error

		// Validate based on database type
		if *useSQLite {
			validationErr = validateSQLiteDirect(sqlStr)
		} else {
			_, validationErr = pg_query.ParseToJSON(sqlStr)
		}

		if validationErr != nil {
			errors = append(errors, sqlError{path: path, err: validationErr.Error()})
			return nil
		}

		fmt.Printf("✓ %s\n", path)
		validCount++
		return nil
	})

	if len(errors) > 0 {
		fmt.Printf("\n")
		for _, sqlErr := range errors {
			fmt.Printf("✗ %s: invalid SQL\n  %s\n", sqlErr.path, sqlErr.err)
		}
		fmt.Printf("\n✗ %d SQL errors found\n", len(errors))
		os.Exit(1)
	}

	fmt.Printf("✓ All %d SQL files validated\n", validCount)
}

// validateSQLiteDirect validates SQL syntax using SQLite's parser
// We use a pragmatic approach: try to prepare statements, creating a dummy table for DML
func validateSQLiteDirect(sqlStr string) error {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return fmt.Errorf("cannot open sqlite: %w", err)
	}
	defer db.Close()

	// For SELECT/INSERT/UPDATE/DELETE, create a dummy table so table references don't fail
	upperSQL := strings.ToUpper(strings.TrimSpace(sqlStr))
	if !strings.HasPrefix(upperSQL, "CREATE") && !strings.HasPrefix(upperSQL, "ALTER") && !strings.HasPrefix(upperSQL, "DROP") {
		// Create dummy tables for common operations
		_, err = db.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, email TEXT, active INTEGER)")
		if err != nil {
			return fmt.Errorf("setup error: %w", err)
		}
		_, err = db.Exec("CREATE TABLE posts (id INTEGER PRIMARY KEY, user_id INTEGER, title TEXT)")
		if err != nil {
			return fmt.Errorf("setup error: %w", err)
		}
	}

	// Try to prepare the statement - syntax errors will be caught
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("SQL error: %w", err)
	}
	stmt.Close()

	return nil
}
