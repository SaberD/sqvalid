package main

import (
	"os"
	"path/filepath"
	"testing"

	pg_query "github.com/pganalyze/pg_query_go/v6"
)

const (
	testdataDir = "testdata"
	postgresDir = "postgres"
	validDir    = "valid"
	invalidDir  = "invalid"
)

func TestPostgreSQLValid(t *testing.T) {
	dir := filepath.Join(testdataDir, postgresDir, validDir)
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read test directory %s: %v", dir, err)
	}

	if len(files) == 0 {
		t.Fatalf("No test files found in %s", dir)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			path := filepath.Join(dir, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read %s: %v", path, err)
				continue
			}

			t.Run(file.Name(), func(t *testing.T) {
				_, pgErr := pg_query.ParseToJSON(string(content))
				if pgErr != nil {
					t.Errorf("Expected valid PostgreSQL in %s, got error: %v", path, pgErr)
				}
			})
		}
	}
}

func TestPostgreSQLInvalid(t *testing.T) {
	dir := filepath.Join(testdataDir, postgresDir, invalidDir)
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read test directory %s: %v", dir, err)
	}

	if len(files) == 0 {
		t.Fatalf("No test files found in %s", dir)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			path := filepath.Join(dir, file.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read %s: %v", path, err)
				continue
			}

			t.Run(file.Name(), func(t *testing.T) {
				_, pgErr := pg_query.ParseToJSON(string(content))
				if pgErr == nil {
					t.Errorf("Expected PostgreSQL error in %s, got nil", path)
				}
			})
		}
	}
}
