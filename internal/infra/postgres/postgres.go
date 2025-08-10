package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"github.com/sunr3d/subscription-aggregator/internal/config"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/models"
)

var _ infra.Database = (*PostgresDB)(nil)

type PostgresDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(cfg config.PostgresConfig, log *zap.Logger) (infra.Database, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}
	poolCfg.MinConns = int32(cfg.MinConns)
	poolCfg.MaxConns = int32(cfg.MaxConns)
	poolCfg.MaxConnLifetime = cfg.MaxConnTTL
	poolCfg.HealthCheckPeriod = cfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres New() -> Ping: %w", err)
	}

	log.Info("Postgres Pool инициализирован",
		zap.String("component", "infra.Database(PostgresDB)"),
		zap.String("host", cfg.Host),
		zap.String("db", cfg.DBName),
		zap.Int32("maxConns", poolCfg.MaxConns),
	)

	return &PostgresDB{pool: pool, logger: log}, nil
}

func (db *PostgresDB) Close() {
	if db.pool != nil {
		db.logger.Info("Postgres Pool закрыт",
			zap.String("component", "infra.Database(PostgresDB)"),
		)
		db.pool.Close()
	}
}

func (db *PostgresDB) Create(ctx context.Context, data models.Subscription) (int, error) {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id int

	if err := db.pool.QueryRow(ctx, query,
		data.ServiceName, data.Price, data.UserID, data.StartDate, data.EndDate,
	).Scan(&id); err != nil {
		return -1, fmt.Errorf("postgres Create(): %w", err)
	}

	return id, nil
}

func (db *PostgresDB) GetByID(ctx context.Context, id int) (models.Subscription, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1;
	`
	var data models.Subscription

	if err := db.pool.QueryRow(ctx, query, id).Scan(
		&data.ID, &data.ServiceName, &data.Price, &data.UserID, &data.StartDate, &data.EndDate,
	); err != nil {
		return models.Subscription{}, fmt.Errorf("postgres GetByID(): %w", err)
	}

	return data, nil
}

func (db *PostgresDB) Update(ctx context.Context, data models.Subscription) error {
	const query = `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6;
	`

	ct, err := db.pool.Exec(ctx, query,
		data.ServiceName, data.Price, data.UserID, data.StartDate, data.EndDate, data.ID,
	)
	if err != nil {
		return fmt.Errorf("postgres Update(): %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("postgres Update(): запись не найдена в БД")
	}

	return nil
}

func (db *PostgresDB) Delete(ctx context.Context, id int) error {
	const query = `
		DELETE FROM subscriptions WHERE id = $1;
	`

	ct, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("postgres Delete(): %w", err)
	}
	if ct.RowsAffected() == 0 {
		return fmt.Errorf("postgres Delete(): запись не найдена в БД")
	}
	return nil
}

func (db *PostgresDB) List(ctx context.Context, filter infra.ListFilter) ([]models.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
	`
	var (
		conds []string
		args  []any
		i     = 1
	)

	if filter.UserID != nil {
		conds = append(conds, fmt.Sprintf("user_id = $%d", i))
		args = append(args, *filter.UserID)
		i++
	}

	if filter.ServiceName != nil {
		conds = append(conds, fmt.Sprintf("service_name = $%d", i))
		args = append(args, *filter.ServiceName)
		i++
	}

	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	query += " ORDER BY id DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres List(): %w", err)
	}

	var data []models.Subscription
	for rows.Next() {
		var dataItem models.Subscription
		if err := rows.Scan(
			&dataItem.ID, &dataItem.ServiceName, &dataItem.Price, &dataItem.UserID,
			&dataItem.StartDate, &dataItem.EndDate,
		); err != nil {
			return nil, fmt.Errorf("postgres List(), rows.Scan(): %w", err)
		}
		data = append(data, dataItem)
	}

	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres List(), rows.Err(): %w", err)
	}

	return data, nil
}
