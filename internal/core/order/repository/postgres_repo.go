package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Jeff1507/go-to-save-api/internal/core/order"
	"github.com/Jeff1507/go-to-save-api/internal/infra/transaction"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepo struct {
	pg *pgxpool.Pool
}

func NewPostgresRepo(pg *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{
		pg: pg,
	}
}

func (repo *PostgresRepo) Save(ctx context.Context, o *order.Order) error {
	tx, ok := transaction.Querier(ctx)
	if ok {
		return repo.save(ctx, tx, o)
	}

	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := repo.save(ctx, tx, o); err != nil {
		return fmt.Errorf("save with transaction: %w", err)
	}

	return tx.Commit(ctx)
}

func (repo *PostgresRepo) save(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	r := toRow(o)

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO orders(
			id,
			customer_id,
			store_id,
			status,
			total_price,
			pickup_start,
			pickup_end,
			confirmation_code,
			completed_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`,
		r.ID,
		r.CustomerID,
		r.StoreID,
		r.Status,
		r.TotalPrice,
		r.PickupStart,
		r.PickupEnd,
		r.ConfirmationCode,
		r.CompletedAt,
		r.CreatedAt,
		r.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert bag: %w", err)
	}

	batch := &pgx.Batch{}
	for _, item := range r.Items {
		batch.Queue(
			`
			INSERT INTO order_items(order_id, bag_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			`,
			r.ID, item.BagID, item.Quantity, item.UnitPrice,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range r.Items {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) Update(ctx context.Context, o *order.Order) error {
	tx, ok := transaction.Querier(ctx)
	if ok {
		return repo.update(ctx, tx, o)
	}

	tx, err := repo.pg.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := repo.update(ctx, tx, o); err != nil {
		return fmt.Errorf("update with transaction: %w", err)
	}

	return tx.Commit(ctx)
}

func (repo *PostgresRepo) update(ctx context.Context, tx pgx.Tx, o *order.Order) error {
	r := toRow(o)

	_, err := tx.Exec(
		ctx,
		`
		UPDATE orders
		SET status = $1, completed_at = $2, updated_at = $3
		WHERE id = $4
		`,
		r.Status, r.CompletedAt, r.UpdatedAt, r.ID,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM order_items WHERE order_id = $1
		`,
		r.ID,
	)
	if err != nil {
		return fmt.Errorf("delete order items: %w", err)
	}

	batch := &pgx.Batch{}
	for _, item := range r.Items {
		batch.Queue(
			`
			INSERT INTO order_items(order_id, bag_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			`,
			r.ID, item.BagID, item.Quantity, item.UnitPrice,
		)
	}

	results := tx.SendBatch(ctx, batch)
	defer results.Close()

	for range r.Items {
		if _, err := results.Exec(); err != nil {
			return fmt.Errorf("insert order item: %w", err)
		}
	}

	if err := results.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	return nil
}

func (repo *PostgresRepo) HasPendingByBag(ctx context.Context, bagID string) (bool, error) {
	var exists bool
	if err := repo.pg.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM order_items i
			JOIN orders o ON o.id = i.order_id
			WHERE i.bag_id = $1 AND o.status = 'PENDING'
		)
		`,
		bagID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check pending order for bag: %w", err)
	}

	return exists, nil
}

func (repo *PostgresRepo) HasPendingByCustomer(
	ctx context.Context,
	customerID string,
) (bool, error) {
	var exists bool
	if err := repo.pg.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM orders WHERE customer_id = $1 AND status = 'PENDING')`,
		customerID,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("check pending order for customer: %w", err)
	}

	return exists, nil
}

func (repo *PostgresRepo) Find(ctx context.Context, id string) (*order.Order, error) {
	row := repo.pg.QueryRow(
		ctx,
		`
		SELECT
			id,
			customer_id,
			store_id,
			status,
			total_price,
			pickup_start,
			pickup_end,
			confirmation_code,
			completed_at,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		`,
		id,
	)

	var r orderRow
	if err := row.Scan(
		&r.ID,
		&r.CustomerID,
		&r.StoreID,
		&r.Status,
		&r.TotalPrice,
		&r.PickupStart,
		&r.PickupEnd,
		&r.ConfirmationCode,
		&r.CompletedAt,
		&r.CreatedAt,
		&r.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, order.ErrNotFound
		}

		return nil, fmt.Errorf("scan order: %w", err)
	}

	rows, err := repo.pg.Query(
		ctx,
		`
		SELECT bag_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1
		`,
		r.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("select order items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i itemRow
		if err := rows.Scan(
			&i.BagID,
			&i.Quantity,
			&i.UnitPrice,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}

		r.Items = append(r.Items, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return toEntity(r), nil
}
