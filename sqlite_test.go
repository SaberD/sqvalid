package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSQLiteValid(t *testing.T) {

	const (
		testdataDir = "testdata"
		sqliteDir   = "sqlite"
		validDir    = "valid"
	)

	dir := filepath.Join(testdataDir, sqliteDir, validDir)
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
				err := validateSQLiteDirect(string(content))
				if err != nil {
					t.Errorf("Expected valid SQLite SQL in %s, got error: %v", path, err)
				}
			})
		}
	}
}

func TestSQLiteInvalid(t *testing.T) {

	const (
		testdataDir = "testdata"
		sqliteDir   = "sqlite"
		invalidDir  = "invalid"
	)

	dir := filepath.Join(testdataDir, sqliteDir, invalidDir)
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
				err := validateSQLiteDirect(string(content))
				if err == nil {
					t.Errorf("Expected SQLite SQL error in %s, got nil", path)
				}
			})
		}
	}
}
