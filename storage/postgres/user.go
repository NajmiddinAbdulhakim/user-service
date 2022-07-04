package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	// "github.com/huandu/go-sqlbuilder"
	pb "github.com/NajmiddinAbdulhakim/user-service/genproto"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(user *pb.User) (*pb.User, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	time := time.Now()

	var usr pb.User
	query := `INSERT INTO users (
        id, first_name, last_name, user_name, email, phone_number, bio, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
        RETURNING id, first_name, last_name, user_name, email, phone_number, bio, status, created_at`
	err = tx.QueryRow(query, user.Id, user.FirstName, user.LastName, user.UserName, user.Email, pq.Array(user.PhoneNumber), user.Bio, user.Status, time).Scan(
		&usr.Id,
		&usr.FirstName,
		&usr.LastName,
		&usr.UserName,
		&usr.Email,
		pq.Array(&usr.PhoneNumber),
		&usr.Bio,
		&usr.Status,
		&usr.CreatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf(`ka;;a ? %v`, err)
	}

	queryAddress := `INSERT INTO addresses (id, user_id, country, city, district, postalcode)
    VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, country, city, district, postalcode`
	for _, addr := range user.Addresses {
		addr_id, _ := uuid.NewV4()
		err = tx.QueryRow(queryAddress, addr_id, usr.Id, addr.Country, addr.City, addr.District, addr.PostalCode).Scan(
			&addr.Id, &addr.Country, &addr.City, &addr.District, &addr.PostalCode,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		usr.Addresses = append(usr.Addresses, addr)
	}
	tx.Commit()
	return &usr, nil
}

func (r *UserRepo) UpdateUser(user *pb.UpdateUserReq) (bool, error) {
	time := time.Now()

	query := `UPDATE users SET first_name = $1, last_name = $2, 
	user_name = $3, email = $4, phone_number = $5, bio = $6, 
	status = $7, updated_at = $8 WHERE id = $9`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UserName,
		user.Email, pq.Array(user.PhoneNumber), user.Bio, user.Status, time, user.Id)
	if err != nil {
		return false, fmt.Errorf(`kalla > %v`, err)
	}
	queryA := `UPDATE addresses SET country = $1, city = $2, district = $3, postalcode = $4 WHERE user_id = $5 AND id = $6`
	for _, addr := range user.Addresses {
		_, err := r.db.Exec(queryA, addr.Country, addr.City, addr.District, addr.PostalCode, user.Id, addr.Id)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func (r *UserRepo) LoginUser(email string) (*pb.User, error) {
	var user *pb.User
	query := `SELECT id, first_name, last_name, username, email, password, 
	phone_number, bio, status FROM users WHERE email = $1`
	err := r.db.QueryRow(query,email).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.Bio,
		&user.Status,
	)
	if err != nil {
		return nil, err
	}
	if user.Id == "" {
		return nil,nil
	}
	return user, nil
}

func (r *UserRepo) GetUserById(userID string) (*pb.User, error) {
	var usr pb.User
	query := `SELECT id, first_name, last_name, user_name, 
	email, phone_number, bio, status FROM users WHERE id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&usr.Id,
			&usr.FirstName,
			&usr.LastName,
			&usr.UserName,
			&usr.Email,
			pq.Array(&usr.PhoneNumber),
			&usr.Bio,
			&usr.Status,
		)
		if err != nil {
			return nil, fmt.Errorf(`error getting user scan by id >%v`, err)
		}

		var add []*pb.Address
		query := `SELECT country, city, district, postalcode FROM addresses WHERE user_id = $1`
		rows, err = r.db.Query(query, usr.Id)
		if err != nil {
			return nil, fmt.Errorf(`error getting user address by id >%v`, err)
		}

		for rows.Next() {
			var adrs pb.Address
			err = rows.Scan(
				&adrs.Country,
				&adrs.City,
				&adrs.District,
				&adrs.PostalCode,
			)
			if err != nil {
				return nil, fmt.Errorf(`error getting user address scan by id >%v`, err)
			}

			add = append(add, &adrs)
		}
		usr.Addresses = add
	}
	return &usr, nil
}

func (r *UserRepo) GetAllUsers() ([]*pb.User, error) {
	var users []*pb.User

	rows, err := r.db.Query(`SELECT id, first_name, last_name, user_name, 
	email, phone_number, bio, status FROM users WHERE deleted_at IS NULL `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr pb.User
		err := rows.Scan(
			&usr.Id,
			&usr.FirstName,
			&usr.LastName,
			&usr.UserName,
			&usr.Email,
			pq.Array(&usr.PhoneNumber),
			&usr.Bio,
			&usr.Status,
		)
		if err != nil {
			return nil, err
		}
		query := `SELECT country, city, district, postalcode FROM addresses`
		rows, err := r.db.Query(query)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var adrs pb.Address
			err = rows.Scan(
				&adrs.Country,
				&adrs.City,
				&adrs.District,
				&adrs.PostalCode,
			)
			if err != nil {
				return nil, err
			}

			usr.Addresses = append(usr.Addresses, &adrs)
		}

		users = append(users, &usr)
	}
	return users, nil
}

func (r *UserRepo) DeleteUser(userID string) (bool, error) {
	time := time.Now()

	query := `UPDATE users SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time, userID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *UserRepo) GetListUsers(page, limit int64) ([]*pb.User, int64, error) {
	var users []*pb.User
	offset := (page - 1) * limit

	query := `SELECT id, first_name, last_name, user_name, 
	email, phone_number, bio, status, updated_at FROM users order by first_name OFFSET $1 LIMIT $2`
	rows, err := r.db.Query(query, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		var user pb.User
		var upd_at sql.NullString
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.UserName,
			&user.Email,
			pq.Array(&user.PhoneNumber),
			&user.Bio,
			&user.Status,
			&upd_at,
			
		)
		user.UpdatedAt = upd_at.String
		if err != nil {
			return nil, 0, err
		}

		queryAddr := `SELECT country, city, district, postalcode FROM addresses WHERE user_id = $1`
		rowss, err := r.db.Query(queryAddr, user.Id)
		if err != nil {
			return nil, 0, fmt.Errorf(`error getting user address by id >%v`, err)
		}

		for rowss.Next() {
			var adrs pb.Address
			err = rowss.Scan(
				&adrs.Country,
				&adrs.City,
				&adrs.District,
				&adrs.PostalCode,
			)
			if err != nil {
				return nil, 0, fmt.Errorf(`error getting user address scan by id >%v`, err)
			}
			user.Addresses = append(user.Addresses, &adrs)
		}
		users = append(users, &user)
	}
	var count int64
	err = r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(
		&count,
	)
	if err != nil {
		return nil, 0, err
	}
	return users, count, nil
}

func (r *UserRepo) CheckUnique(field, value string) (bool, error) {
	var exists int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users WHERE $1 = $2`, field, value).Scan(
		&exists,
	)
	if err != nil {
		return false, err
	}
	if exists > 0 {
		return true, nil
	}
	return false, nil
}

// func (r *UserRepo) GetFilteredUser(userID string)([]*pb.User) {
// 	sql := sqlbuilder.Select("id","first_name","last_name").From("users"), 
// 		Where("id = $1")
	
// 	var user *pb.User
// 	err := r.db.QueryRow(sql, userID).Scan(
// 		&user.Id,
// 		&user.FirstName,
// 		&user.LastName
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user,err
// }