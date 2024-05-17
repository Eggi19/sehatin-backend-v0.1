package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepoOpt struct {
	Db *sql.DB
}

type UserRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*entities.User, error)
	FindOneById(ctx context.Context, userId int64) (*entities.User, error)
	DeleteOne(ctx context.Context, userId int64) error
	FindAll(ctx context.Context, params entities.UserParams) ([]entities.User, int, error)
	UpdateOne(ctx context.Context, user entities.User) error
	CreateOneUser(ctx context.Context, u entities.User) (*entities.User, error)
	UserVerificationToken(ctx context.Context, userId int64, token string, exp time.Time) error
	VerifyUser(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, userId int64, newPassword string) error
	FindPasswordById(ctx context.Context, userId int64) (*entities.User, error)
}

type UserRepositoryPostgres struct {
	db *sql.DB
}

func NewUserRepositoryPostgres(trOpt *UserRepoOpt) UserRepository {
	return &UserRepositoryPostgres{
		db: trOpt.Db,
	}
}

func (r *UserRepositoryPostgres) FindOneByEmail(ctx context.Context, email string) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.IsVerified, &u.IsGoogle, &u.IsOnline, &u.ProfilePicture)
	} else {
		err = r.db.QueryRowContext(ctx, qFindUserByEmail, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password, &u.IsVerified, &u.IsGoogle, &u.IsOnline, &u.ProfilePicture)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryPostgres) FindOneById(ctx context.Context, userId int64) (*entities.User, error) {
	u := entities.User{
		Gender: &entities.Gender{},
	}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindOneUserById, userId).Scan(
			&u.Id, &u.Name, &u.Email, &u.BirthDate, &u.Gender.Id, &u.Gender.Name, &u.IsVerified, &u.IsGoogle, &u.IsOnline, &u.ProfilePicture,
		)
	} else {
		err = r.db.QueryRowContext(ctx, qFindOneUserById, userId).Scan(
			&u.Id, &u.Name, &u.Email, &u.BirthDate, &u.Gender.Id, &u.Gender.Name, &u.IsVerified, &u.IsGoogle, &u.IsOnline, &u.ProfilePicture,
		)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &u, nil
}

func (r *UserRepositoryPostgres) DeleteOne(ctx context.Context, userId int64) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qDeleteUserById)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qDeleteUserById)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *UserRepositoryPostgres) FindAll(ctx context.Context, params entities.UserParams) ([]entities.User, int, error) {
	users := []entities.User{}

	var totalRows int

	var sb strings.Builder
	sb.WriteString(qCountTotalRows)
	sb.WriteString(qUserColl)
	sb.WriteString(qUserCommand)

	var sbTotalRows strings.Builder
	sbTotalRows.WriteString(qCountTotalRows)
	sbTotalRows.WriteString(qUserCommand)

	values := []interface{}{}
	valuesCountTotal := []interface{}{}

	numberOfArgs := 1

	if params.GenderId != 0 {
		sb.WriteString(`AND g.id = `)
		sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		values = append(values, params.GenderId)

		sbTotalRows.WriteString(`AND g.id = `)
		sbTotalRows.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
		valuesCountTotal = append(valuesCountTotal, params.GenderId)

		numberOfArgs++
	}

	if params.Keyword != "" {
		sb.WriteString(`AND u.name ILIKE '%`)
		sb.WriteString(params.Keyword)
		sb.WriteString(`%' `)

		sbTotalRows.WriteString(`AND u.name ILIKE '%`)
		sbTotalRows.WriteString(params.Keyword)
		sbTotalRows.WriteString(`%' `)
	}

	var sortBy string
	switch params.SortBy {
	case "birth-date":
		sortBy = `u.birth_date `
	case "name":
		sortBy = `u.name `
	default:
		sortBy = `u.id `
	}
	sb.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))
	sbTotalRows.WriteString(fmt.Sprintf(`ORDER BY %s `, sortBy))

	if params.Sort == "" {
		params.Sort = `ASC `
	}
	sb.WriteString(fmt.Sprintf(`%s `, params.Sort))
	sbTotalRows.WriteString(fmt.Sprintf(`%s `, params.Sort))

	sb.WriteString(`LIMIT `)
	sb.WriteString(fmt.Sprintf(`$%d `, numberOfArgs))
	values = append(values, params.Limit)
	numberOfArgs++

	sb.WriteString(`OFFSET `)
	sb.WriteString(fmt.Sprintf(`$%d`, numberOfArgs))
	values = append(values, params.Limit*(params.Page-1))
	numberOfArgs++

	rows, err := r.db.QueryContext(ctx, sb.String(), values...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		u := entities.User{}
		u.Gender = &entities.Gender{}
		err := rows.Scan(&totalRows,
			&u.Id, &u.Name, &u.Email, &u.BirthDate, &u.Gender.Id, &u.Gender.Name, &u.IsVerified, &u.IsGoogle, &u.IsOnline, &u.ProfilePicture,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, 0, custom_errors.NotFound(err)
			}
			return nil, 0, err
		}
		users = append(users, u)
	}

	if totalRows == 0 {
		err := r.db.QueryRowContext(ctx, sbTotalRows.String(), valuesCountTotal...).Scan(&totalRows)
		if err != nil && err != sql.ErrNoRows {
			return nil, 0, err
		}
	}

	return users, totalRows, nil
}

func (r *UserRepositoryPostgres) UpdateOne(ctx context.Context, user entities.User) error {
	values := []interface{}{}
	values = append(values, user.Id)
	values = append(values, user.Name)
	values = append(values, user.BirthDate)
	values = append(values, user.Gender.Id)

	var err error
	var stmt *sql.Stmt

	if user.ProfilePicture.Valid {
		values = append(values, user.ProfilePicture)
		tx := extractTx(ctx)
		if tx != nil {
			stmt, err = tx.PrepareContext(ctx, qUpdateOneUser)
		} else {
			stmt, err = r.db.PrepareContext(ctx, qUpdateOneUser)
		}
	} else {
		tx := extractTx(ctx)
		if tx != nil {
			stmt, err = tx.PrepareContext(ctx, qUpdateOneUserWithOutProfile)
		} else {
			stmt, err = r.db.PrepareContext(ctx, qUpdateOneUserWithOutProfile)
		}
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *UserRepositoryPostgres) CreateOneUser(ctx context.Context, u entities.User) (*entities.User, error) {
	newU := entities.User{}

	values := []interface{}{}
	values = append(values, u.Name)
	values = append(values, u.Email)
	values = append(values, u.BirthDate)

	q := qCreateOneUser
	if u.Gender.Id > 0 {
		q = strings.Replace(q, "password", "gender_id, password", 1)
		values = append(values, u.Gender.Id)
	}

	if u.Password.Valid {
		values = append(values, u.Password)
	}

	if u.IsGoogle {
		q = strings.Replace(q, "password", "profile_picture, is_google, is_verified", 1)
		values = append(values, u.ProfilePicture)
		values = append(values, u.IsGoogle)
		values = append(values, u.IsVerified)
	}

	valuesStr := ""
	for i := 1; i <= len(values); i++ {
		valuesStr += "$" + strconv.Itoa(i)
		if i < len(values) {
			valuesStr += ", "
		}
	}

	q = strings.Replace(q, "--values--", valuesStr, 1)

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, q, values...).Scan(&newU.Id)
	} else {
		err = r.db.QueryRowContext(ctx, q, values...).Scan(&newU.Id)
	}

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.ViolatesUniqueConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.UserEmailNotUniqueErrMsg)
		}
		if errors.As(err, &pgErr) && pgErr.Code == constants.VioletesForeignKeyConstraintPgErrCode {
			return nil, custom_errors.BadRequest(err, constants.GenderNotFoundErrMsg)
		}
		return nil, err
	}

	return &newU, nil
}

func (r *UserRepositoryPostgres) UserVerificationToken(ctx context.Context, userId int64, token string, exp time.Time) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qCreateUserVerificationToken, token, exp, userId)
	} else {
		_, err = r.db.ExecContext(ctx, qCreateUserVerificationToken, token, exp, userId)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPostgres) VerifyUser(ctx context.Context, email string) error {
	var err error

	tx := extractTx(ctx)
	if tx != nil {
		_, err = tx.ExecContext(ctx, qVerifyUser, email)
	} else {
		_, err = r.db.ExecContext(ctx, qVerifyUser, email)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepositoryPostgres) UpdatePassword(ctx context.Context, userId int64, newPassword string) error {
	var err error
	var stmt *sql.Stmt

	tx := extractTx(ctx)
	if tx != nil {
		stmt, err = tx.PrepareContext(ctx, qUpadateUserPassword)
	} else {
		stmt, err = r.db.PrepareContext(ctx, qUpadateUserPassword)
	}

	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, userId, newPassword)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return custom_errors.NotFound(sql.ErrNoRows)
	}

	return nil
}

func (r *UserRepositoryPostgres) FindPasswordById(ctx context.Context, userId int64) (*entities.User, error) {
	u := entities.User{}

	var err error

	tx := extractTx(ctx)
	if tx != nil {
		err = tx.QueryRowContext(ctx, qFindUserPasswordById, userId).Scan(
			&u.Id, &u.Password)
	} else {
		err = r.db.QueryRowContext(ctx, qFindUserPasswordById, userId).Scan(
			&u.Id, &u.Password)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.NotFound(err)
		}
		return nil, err
	}

	return &u, nil
}
