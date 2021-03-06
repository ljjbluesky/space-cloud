package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spaceuptech/space-cloud/gateway/utils"
)

// RawBatch performs a batch operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawBatch(ctx context.Context, queries []string) error {
	// Skip if length of queries == 0
	if len(queries) == 0 {
		return nil
	}

	tx, err := s.client.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	for _, query := range queries {
		_, err := tx.ExecContext(ctx, query)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

// RawExec performs an operation for schema creation
// NOTE: not to be exposed externally
func (s *SQL) RawExec(ctx context.Context, query string) error {
	_, err := s.client.ExecContext(ctx, query, []interface{}{}...)
	return err
}

// GetConnectionState : Function to get connection state
func (s *SQL) GetConnectionState(ctx context.Context) bool {
	if !s.enabled || s.client == nil {
		return false
	}

	// Ping to check if connection is established
	err := s.client.PingContext(ctx)
	return err == nil
}

// CreateProjectIfNotExist creates a schema / database
func (s *SQL) CreateProjectIfNotExist(ctx context.Context, project string) error {
	var sql string
	switch utils.DBType(s.dbType) {
	case utils.MySQL:
		sql = "create database if not exists " + project
	case utils.Postgres:
		sql = "create schema if not exists " + project
	case utils.SqlServer:
		sql = `IF (NOT EXISTS (SELECT * FROM sys.schemas WHERE name = '` + project + `')) 
					BEGIN
    					EXEC ('CREATE SCHEMA [` + project + `] ')
					END`
	default:
		return fmt.Errorf("invalid db type (%s) provided", s.dbType)
	}
	return s.RawExec(ctx, sql)
}
