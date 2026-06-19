package database

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	json "github.com/nef734/ggit/internal/jsonx"
	log "github.com/nef734/ggit/internal/log"
	"github.com/nef734/ggit/internal/osx"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// getTableType returns the type name of a table definition without package name,
// e.g. *database.LFSObject -> LFSObject.
func getTableType(t any) string {
	return strings.TrimPrefix(fmt.Sprintf("%T", t), "*database.")
}

// DumpDatabase dumps all data from database to file system in JSON Lines format.
func DumpDatabase(ctx context.Context, db *gorm.DB, dirPath string, verbose bool) error {
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return err
	}

	// All tables are now managed by GORM and included in the main 'Tables' list.

	for _, table := range Tables {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tableName := getTableType(table)
		if verbose {
			log.Trace("Dumping table %q...", tableName)
		}

		err := func() error {
			tableFile := filepath.Join(dirPath, tableName+".json")
			f, err := os.Create(tableFile)
			if err != nil {
				return fmt.Errorf("create table file"+": %v", err)
			}
			defer func() { _ = f.Close() }()

			return dumpTable(ctx, db, table, f)
		}()
		if err != nil {
			return fmt.Errorf("dump table %q: %v", tableName, err)
		}
	}

	return nil
}

func dumpTable(ctx context.Context, db *gorm.DB, table any, w io.Writer) error {
	query := db.WithContext(ctx).Model(table)
	switch table.(type) {
	case *LFSObject:
		query = query.Order("repo_id, oid ASC")
	default:
		query = query.Order("id ASC")
	}

	rows, err := query.Rows()
	if err != nil {
		return fmt.Errorf("select rows"+": %v", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		elem := reflect.New(reflect.TypeOf(table).Elem()).Interface()
		err = db.ScanRows(rows, elem)
		if err != nil {
			return fmt.Errorf("scan rows"+": %v", err)
		}

		switch e := elem.(type) {
		case *LFSObject:
			e.Created = e.Created.UTC()
		}

		err = json.NewEncoder(w).Encode(elem)
		if err != nil {
			return fmt.Errorf("encode JSON"+": %v", err)
		}
	}
	return rows.Err()
}

// ImportDatabase imports data from backup archive in JSON Lines format.
func ImportDatabase(ctx context.Context, db *gorm.DB, dirPath string, verbose bool) error {
	// All tables are now managed by GORM and included in the main 'Tables' list.

	for _, table := range Tables {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		tableName := strings.TrimPrefix(fmt.Sprintf("%T", table), "*database.")
		err := func() error {
			tableFile := filepath.Join(dirPath, tableName+".json")
			if !osx.IsFile(tableFile) {
				log.Info("Skipped table %q", tableName)
				return nil
			}

			if verbose {
				log.Trace("Importing table %q...", tableName)
			}

			f, err := os.Open(tableFile)
			if err != nil {
				return fmt.Errorf("open table file"+": %v", err)
			}
			defer func() { _ = f.Close() }()

			return importTable(ctx, db, table, f)
		}()
		if err != nil {
			return fmt.Errorf("import table %q: %v", tableName, err)
		}
	}

	return nil
}

func importTable(ctx context.Context, db *gorm.DB, table any, r io.Reader) error {
	err := db.WithContext(ctx).Migrator().DropTable(table)
	if err != nil {
		return fmt.Errorf("drop table"+": %v", err)
	}

	err = db.WithContext(ctx).Migrator().AutoMigrate(table)
	if err != nil {
		return fmt.Errorf("auto migrate"+": %v", err)
	}

	_, err = schema.Parse(table, &sync.Map{}, db.NamingStrategy)
	if err != nil {
		return fmt.Errorf("parse schema"+": %v", err)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		// PostgreSQL does not like the null characters (U+0000)
		cleaned := bytes.ReplaceAll(scanner.Bytes(), []byte("\\u0000"), []byte(""))

		elem := reflect.New(reflect.TypeOf(table).Elem()).Interface()
		err = json.Unmarshal(cleaned, elem)
		if err != nil {
			return fmt.Errorf("unmarshal JSON to struct"+": %v", err)
		}

		err = db.WithContext(ctx).Create(elem).Error
		if err != nil {
			return fmt.Errorf("create row"+": %v", err)
		}
	}

	return nil
}
