package persons

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/alukart32/effective-mobile-test-task/internal/person/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type pgxDB struct {
	pool *pgxpool.Pool
}

func (p *pgxDB) Save(ctx context.Context, person model.Person) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("pgxDB.Save: %w", err)
		}
	}()
	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return err
	}
	defer func() {
		err = p.finishTx(ctx, tx, err)
	}()

	const query = `INSERT INTO
	persons(id, name, surname, patronymic, nation, gender, age)
	VALUES($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.Exec(ctx, query,
		person.Id,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Nation,
		person.Gender,
		person.Age,
	)

	var pgErr *pgconn.PgError
	if err != nil && errors.As(err, &pgErr) {
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.SQLState()) &&
			pgErr.SQLState() == pgerrcode.UniqueViolation {
			err = fmt.Errorf("pgxDB.Save: %s unique violation", pgErr.ColumnName)
		}
		if pgerrcode.IsIntegrityConstraintViolation(pgErr.SQLState()) &&
			pgErr.SQLState() == pgerrcode.CheckViolation {
			err = fmt.Errorf("pgxDB.Save: %s check violation", pgErr.ColumnName)
		}
	}
	return err
}

func (p *pgxDB) FindById(ctx context.Context, id string) (_ model.Person, err error) {
	const query = `SELECT * FROM persons WHERE id = $1`
	var record record
	row := p.pool.QueryRow(ctx, query, id)
	err = row.Scan(
		&record.Id,
		&record.Name,
		&record.Surname,
		&record.Patronymic,
		&record.Nation,
		&record.Gender,
		&record.Age,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return model.Person{}, fmt.Errorf("pgxDB.FindById: %w", err)
	}
	return record.ToModel(), nil
}

func (p *pgxDB) Collect(ctx context.Context, filter model.PersonFilter, limit, offset int) (_ []model.Person, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("pgxDB.Collect: %w", err)
		}
	}()

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, err
	}
	defer func() error {
		return p.finishTx(ctx, tx, err)
	}()

	query, args := p.getCollectQuery(limit, offset, filter)
	var rows pgx.Rows
	if len(args) > 0 {
		rows, err = tx.Query(ctx, query, args...)
	} else {
		rows, err = tx.Query(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]record, 0)
	for rows.Next() {
		var r record
		err = rows.Scan(
			&r.Id,
			&r.Name,
			&r.Surname,
			&r.Patronymic,
			&r.Nation,
			&r.Gender,
			&r.Age)
		if err != nil {
			return nil, err
		}
		records = append(records, r)
	}
	if rows.Err() != nil {
		return nil, err
	}

	persons := make([]model.Person, len(records))
	for i, r := range records {
		persons[i] = r.ToModel()
	}
	return persons, err
}

func (p *pgxDB) getCollectQuery(limit, offset int, filter model.PersonFilter) (string, []any) {
	var (
		sb         strings.Builder
		fieldOrder int
		args       []any
	)
	addFilters := func() {
		if filter.IsEmpty() {
			return
		}
		sb.WriteString(" WHERE ")

		if filter.OlderThan != 0 {
			fieldOrder++
			sb.WriteString(fmt.Sprintf("age > $%d", fieldOrder))
			args = append(args, filter.OlderThan)
		}
		if filter.YoungerThan != 0 {
			if fieldOrder > 0 {
				sb.WriteString(" AND ")
			}
			fieldOrder++
			sb.WriteString(fmt.Sprintf("age < $%d", fieldOrder))
			args = append(args, filter.YoungerThan)
		}
		if len(filter.Gender) != 0 {
			if fieldOrder > 0 {
				sb.WriteString(" AND ")
			}
			fieldOrder++
			sb.WriteString(fmt.Sprintf("gender = $%d", fieldOrder))
			args = append(args, filter.Gender)
		}
		if len(filter.Nations) > 0 {
			if fieldOrder > 0 {
				sb.WriteString(" AND ")
			}
			if len(filter.Nations) == 1 {
				fieldOrder++
				sb.WriteString(fmt.Sprintf("nation = $%d", fieldOrder))
				args = append(args, filter.Nations[0])
			} else {
				sb.WriteString("nation IN ( ")
				for i, n := range filter.Nations {
					fieldOrder++
					if i < len(filter.Nations)-1 {
						sb.WriteString(fmt.Sprintf("$%d,", fieldOrder))
					} else {
						sb.WriteString(fmt.Sprintf("$%d", fieldOrder))
					}
					args = append(args, n)
				}
				sb.WriteString(" )")
			}
		}
	}
	sb.WriteString("SELECT p.id, p.name, p.surname, p.patronymic, p.nation, p.gender, p.age")
	sb.WriteString(" FROM persons AS p")

	if limit > 0 {
		if offset > 0 {
			sb.WriteString(" JOIN (SELECT id FROM persons ")
			if !filter.IsEmpty() {
				addFilters()
			}
			sb.WriteString(" ORDER BY id")
			sb.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", fieldOrder+1, fieldOrder+2))
			args = append(args, limit, offset)
			fieldOrder += 2

			sb.WriteString(" ) as tmp ON tmp.id = p.id")
		} else {
			if !filter.IsEmpty() {
				addFilters()
			}
			fieldOrder++
			sb.WriteString(fmt.Sprintf(" LIMIT $%d", fieldOrder))
			args = append(args, limit)
		}
	} else if !filter.IsEmpty() {
		addFilters()
	}
	return sb.String(), args
}

func (p *pgxDB) Update(ctx context.Context, id string, meta model.PersonalMetaData) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("pgxDB.Update: %w", err)
		}
	}()

	tx, err := p.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       pgx.RepeatableRead,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return err
	}
	defer func() {
		err = p.finishTx(ctx, tx, err)
	}()
	query, args := p.getUpdateQuery(id, meta)
	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (p *pgxDB) getUpdateQuery(id string, meta model.PersonalMetaData) (string, []any) {
	var (
		sb         strings.Builder
		fieldOrder int
		args       []any
	)
	sb.WriteString("UPDATE persons SET ")
	if len(meta.Nation) != 0 {
		fieldOrder++

		sb.WriteString(fmt.Sprintf("nation = $%v", fieldOrder))
		args = append(args, meta.Nation)
	}
	if len(meta.Gender) != 0 {
		if fieldOrder > 0 {
			sb.WriteString(", ")
		}
		fieldOrder++
		sb.WriteString(fmt.Sprintf("gender = $%v", fieldOrder))
		args = append(args, meta.Gender)
	}
	if meta.Age > 0 {
		if fieldOrder > 0 {
			sb.WriteString(", ")
		}
		fieldOrder++

		sb.WriteString(fmt.Sprintf("age = $%v", fieldOrder))
		args = append(args, meta.Age)
	}

	sb.WriteString(fmt.Sprintf(" WHERE id = $%v", fieldOrder+1))
	args = append(args, id)

	return sb.String(), args
}

func (p *pgxDB) Delete(ctx context.Context, id string) (err error) {
	const query = `DELETE FROM persons WHERE id = $1`
	_, err = p.pool.Exec(ctx, query, id)
	if err != nil {
		err = fmt.Errorf("pgxDB.Delete: %w", err)
	}
	return err
}

// finishTx rollbacks transaction if error is provided.
// If err is nil transaction is committed.
func (p *pgxDB) finishTx(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}
	if commitErr := tx.Commit(ctx); commitErr != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

type cachedStorage struct {
	db    *pgxDB
	cache *redis.Client

	mtx *sync.RWMutex
}

func CachedStorage(db *pgxpool.Pool, cache *redis.Client) (*cachedStorage, error) {
	if cache == nil {
		return nil, fmt.Errorf("init cached storage: redis client is nil")
	}
	if db == nil {
		return nil, fmt.Errorf("init cached storage: postgres pool is nil")
	}

	return &cachedStorage{
		cache: cache,
		db:    &pgxDB{db},
		mtx:   &sync.RWMutex{},
	}, nil
}

func (s *cachedStorage) Save(ctx context.Context, person model.Person) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if _, err := s.cache.Pipelined(ctx, func(rdb redis.Pipeliner) error {
		key := "person:" + person.Id
		rdb.HSet(ctx, key, "id", person.Id)
		rdb.HSet(ctx, key, "name", person.Name)
		rdb.HSet(ctx, key, "surname", person.Surname)
		rdb.HSet(ctx, key, "patronymic", person.Patronymic)
		rdb.HSet(ctx, key, "nation", person.Nation)
		rdb.HSet(ctx, key, "gender", person.Gender)
		rdb.HSet(ctx, key, "age", person.Age)
		return nil
	}); err != nil {
		return err
	}

	var err error
	if err = s.db.Save(ctx, person); err != nil {
		if cacheErr := s.cache.Del(ctx, "person:"+person.Id).Err(); cacheErr != nil {
			return fmt.Errorf("redis: %w", cacheErr)
		}
	}
	return err
}

func (s *cachedStorage) FindById(ctx context.Context, id string) (model.Person, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	var record record
	err := s.cache.HGetAll(ctx, "person:"+id).Scan(&record)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return model.Person{}, err
		}

		p, err := s.db.FindById(ctx, id)
		if err != nil {
			return model.Person{}, err
		}
		if _, err = s.cache.Pipelined(ctx, func(rdb redis.Pipeliner) error {
			key := "person:" + p.Id
			rdb.HSet(ctx, key, "id", p.Id)
			rdb.HSet(ctx, key, "name", p.Name)
			rdb.HSet(ctx, key, "surname", p.Surname)
			rdb.HSet(ctx, key, "patronymic", p.Patronymic)
			rdb.HSet(ctx, key, "nation", p.Nation)
			rdb.HSet(ctx, key, "gender", p.Gender)
			rdb.HSet(ctx, key, "age", p.Age)
			return nil
		}); err != nil {
			return model.Person{}, err
		}
		return p, nil
	}

	return record.ToModel(), nil
}

func (s *cachedStorage) Update(ctx context.Context, id string, meta model.PersonalMetaData) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if err := s.db.Update(ctx, id, meta); err != nil {
		return err
	}
	if err := s.cache.HSet(ctx, "person:"+id, s.getUpdateArgs(meta)).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	return nil
}

func (s *cachedStorage) getUpdateArgs(meta model.PersonalMetaData) map[string]any {
	values := make(map[string]any)
	if len(meta.Nation) != 0 {
		values["nation"] = meta.Nation
	}
	if len(meta.Gender) != 0 {
		values["gender"] = meta.Gender
	}
	if meta.Age > 0 {
		values["age"] = meta.Age
	}
	return values
}

func (s *cachedStorage) Collect(
	ctx context.Context,
	filter model.PersonFilter,
	limit, offset int,
) ([]model.Person, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.db.Collect(ctx, filter, limit, offset)
}

func (s *cachedStorage) Delete(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if err := s.cache.Del(ctx, "person:"+id).Err(); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	return s.db.Delete(ctx, id)
}
