package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	j "github.com/OlesyaNovikova/metricsallert.git/internal/models"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(ctx context.Context, db *sql.DB) (*PostgresDB, error) {
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Ошибка соединения с базой: %v \n", err)
		return nil, err
	}

	err = retry(ctx,
		func(ctx context.Context) error {
			_, err = db.ExecContext(ctx,
				`CREATE TABLE IF NOT EXISTS gauge("name" varchar(50) UNIQUE,"value" double precision)`)
			return err
		})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы gauge: %v \n", err)
		return nil, err
	}

	err = retry(ctx,
		func(ctx context.Context) error {
			_, err = db.ExecContext(ctx,
				`CREATE TABLE IF NOT EXISTS counter("name" varchar(50) UNIQUE,"delta" bigint)`)
			return err
		})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы counter: %v \n", err)
		return nil, err
	}

	pdb := PostgresDB{
		db: db,
	}

	return &pdb, nil
}

func (p *PostgresDB) UpdateGauge(ctx context.Context, name string, value float64) error {
	err := retry(ctx,
		func(ctx context.Context) error {
			return p.updateGauge(ctx, name, value)
		})
	return err
}

func (p *PostgresDB) updateGauge(ctx context.Context, name string, value float64) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO gauge (name, value) VALUES($1,$2)
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`, name, value)
	return err
}

func (p *PostgresDB) UpdateCounter(ctx context.Context, name string, delta int64) (int64, error) {
	var err error
	var val int64
	err = retry(ctx,
		func(ctx context.Context) error {
			val, err = p.updateCounter(ctx, name, delta)
			return err
		})
	return val, err
}

func (p *PostgresDB) updateCounter(ctx context.Context, name string, delta int64) (int64, error) {
	row := p.db.QueryRowContext(ctx,
		`INSERT INTO counter AS c (name, delta) VALUES($1,$2)
		ON CONFLICT (name) DO UPDATE SET delta = c.delta + EXCLUDED.delta 
		RETURNING delta`, name, delta)
	var val int64
	err := row.Scan(&val)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return val, nil
}

func (p *PostgresDB) GetGauge(ctx context.Context, name string) (float64, error) {
	var err error
	var val float64
	err = retry(ctx,
		func(ctx context.Context) error {
			val, err = p.getGauge(ctx, name)
			return err
		})
	return val, err
}

func (p *PostgresDB) getGauge(ctx context.Context, name string) (float64, error) {
	row := p.db.QueryRowContext(ctx,
		"SELECT value FROM gauge WHERE name = $1", name)
	var val float64
	err := row.Scan(&val)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return val, nil
}

func (p *PostgresDB) GetCounter(ctx context.Context, name string) (int64, error) {
	var err error
	var val int64
	err = retry(ctx,
		func(ctx context.Context) error {
			val, err = p.getCounter(ctx, name)
			return err
		})
	return val, err
}

func (p *PostgresDB) getCounter(ctx context.Context, name string) (int64, error) {
	row := p.db.QueryRowContext(ctx,
		"SELECT delta FROM counter WHERE name = $1", name)
	var val int64
	err := row.Scan(&val)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return val, nil
}

func (p *PostgresDB) GetAll(ctx context.Context) (map[string]string, error) {
	var err error
	allMems := make(map[string]string)
	err = retry(ctx,
		func(ctx context.Context) error {
			allMems, err = p.getAll(ctx)
			return err
		})
	return allMems, err
}

func (p *PostgresDB) getAll(ctx context.Context) (map[string]string, error) {

	allMems := make(map[string]string)
	var name string
	var value float64
	var delta int64

	//запрашиваем таблицу gauge
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM gauge")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			return nil, err
		}
		val := strconv.FormatFloat(float64(value), 'f', -1, 64)
		allMems[name] = val
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	//запрашиваем таблицу counter
	rows, err = p.db.QueryContext(ctx, "SELECT * FROM counter")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// пробегаем по всем записям
	for rows.Next() {
		err = rows.Scan(&name, &delta)
		if err != nil {
			return nil, err
		}
		del := strconv.FormatInt(int64(delta), 10)
		allMems[name] = del
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return allMems, nil
}

func (p *PostgresDB) Ping(ctx context.Context) error {
	err := retry(ctx,
		func(ctx context.Context) error {
			return p.ping(ctx)
		})
	return err
}

func (p *PostgresDB) ping(ctx context.Context) error {
	err := p.db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Ошибка соединения с базой: %v \n", err)
		return err
	}
	return nil
}

func (p *PostgresDB) Updates(ctx context.Context, mems []j.Metrics) error {
	err := retry(ctx,
		func(ctx context.Context) error {
			return p.updates(ctx, mems)
		})
	return err
}

func (p *PostgresDB) updates(ctx context.Context, mems []j.Metrics) error {
	if len(mems) == 0 {
		return nil
	}
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtGauge, err := tx.PrepareContext(ctx,
		`INSERT INTO gauge (name, value) VALUES($1,$2)
		ON CONFLICT (name) DO UPDATE SET value = EXCLUDED.value`)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()
	stmtCount, err := tx.PrepareContext(ctx,
		`INSERT INTO counter AS c (name, delta) VALUES($1,$2)
		ON CONFLICT (name) DO UPDATE SET delta = c.delta + EXCLUDED.delta`)
	if err != nil {
		return err
	}
	defer stmtCount.Close()

	for _, mem := range mems {
		if mem.ID == "" {
			return fmt.Errorf("no name")
		}

		switch mem.MType {
		case "gauge":
			if mem.Value == nil {
				return fmt.Errorf("no value")
			}
			_, err := stmtGauge.ExecContext(ctx, mem.ID, *mem.Value)
			if err != nil {
				return err
			}

		case "counter":
			if mem.Delta == nil {
				return fmt.Errorf("no delta")
			}
			_, err := stmtCount.ExecContext(ctx, mem.ID, *mem.Delta)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("bad type")
		}
	}
	return tx.Commit()
}

func retry(ctx context.Context, f func(ctx context.Context) error) error {
	var err error
	delay := [3]int{1, 3, 5}
	err = f(ctx)
	if err != nil {
		for _, t := range delay {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				time.Sleep(time.Duration(t) * time.Second)
				err = f(ctx)
				if err == nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return err
}
