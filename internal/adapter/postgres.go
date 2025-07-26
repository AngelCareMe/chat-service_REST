package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type PostgresAdapter struct {
	Pool   *pgxpool.Pool
	logger *logrus.Logger
}

func NewPostgresAdapter(pool *pgxpool.Pool, logger *logrus.Logger) *PostgresAdapter {
	return &PostgresAdapter{
		Pool:   pool,
		logger: logger,
	}
}

func (p *PostgresAdapter) Close() {
	p.logger.Info("closing database connection pool")
	if p.Pool != nil {
		p.Pool.Close()
	}
	p.logger.Info("database connection pool closed")
}

// BeginTx starts a new transaction
func (p *PostgresAdapter) BeginTx(ctx context.Context) (pgx.Tx, error) {
	p.logger.Debug("beginning database transaction")
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		p.logger.WithError(err).Error("failed to begin transaction")
		return nil, err
	}
	return tx, nil
}

// Exec executes a query
func (p *PostgresAdapter) Exec(ctx context.Context, query string, args ...interface{}) error {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("executing database query")

	start := time.Now()
	_, err := p.Pool.Exec(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		p.logger.WithError(err).WithFields(logrus.Fields{
			"query":      query,
			"duration":   duration,
			"args_count": len(args),
		}).Error("database query failed")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("database query executed successfully")

	return nil
}

// ExecTx executes a query within a transaction
func (p *PostgresAdapter) ExecTx(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) error {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("executing database query in transaction")

	start := time.Now()
	_, err := tx.Exec(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		p.logger.WithError(err).WithFields(logrus.Fields{
			"query":      query,
			"duration":   duration,
			"args_count": len(args),
		}).Error("database query in transaction failed")
		return err
	}

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("database query in transaction executed successfully")

	return nil
}

// QueryRow executes a query that returns a single row
func (p *PostgresAdapter) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("querying single row from database")

	start := time.Now()
	row := p.Pool.QueryRow(ctx, query, args...)
	duration := time.Since(start)

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("single row query executed")

	return row
}

// QueryRowTx executes a query that returns a single row within a transaction
func (p *PostgresAdapter) QueryRowTx(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) pgx.Row {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("querying single row from database in transaction")

	start := time.Now()
	row := tx.QueryRow(ctx, query, args...)
	duration := time.Since(start)

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("single row query in transaction executed")

	return row
}

// Query executes a query that returns multiple rows
func (p *PostgresAdapter) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("querying multiple rows from database")

	start := time.Now()
	rows, err := p.Pool.Query(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		p.logger.WithError(err).WithFields(logrus.Fields{
			"query":      query,
			"duration":   duration,
			"args_count": len(args),
		}).Error("multiple rows query failed")
		return nil, err
	}

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("multiple rows query executed successfully")

	return rows, nil
}

// QueryTx executes a query that returns multiple rows within a transaction
func (p *PostgresAdapter) QueryTx(ctx context.Context, tx pgx.Tx, query string, args ...interface{}) (pgx.Rows, error) {
	p.logger.WithFields(logrus.Fields{
		"query":      query,
		"args_count": len(args),
	}).Debug("querying multiple rows from database in transaction")

	start := time.Now()
	rows, err := tx.Query(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		p.logger.WithError(err).WithFields(logrus.Fields{
			"query":      query,
			"duration":   duration,
			"args_count": len(args),
		}).Error("multiple rows query in transaction failed")
		return nil, err
	}

	p.logger.WithFields(logrus.Fields{
		"duration":   duration,
		"args_count": len(args),
	}).Debug("multiple rows query in transaction executed successfully")

	return rows, nil
}
