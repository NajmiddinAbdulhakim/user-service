package postgres

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	pb "github.com/template-service/genproto"
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
	id := uuid.New()
	time := time.Now()

	var usr pb.User
	query := `INSERT INTO users (
        id, first_name, last_name, user_name, email, phone_number, bio, status, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
        RETURNING id, first_name, last_name, user_name, email, phone_number, bio, status, created_at, updated_at`
	err = tx.QueryRow(query, id, user.FirstName, user.LastName, user.UserName, user.Email, pq.Array(user.PhoneNumber), user.Bio, user.Status, time, time).Scan(
		&usr.Id,
		&usr.FirstName,
		&usr.LastName,
		&usr.UserName,
		&usr.Email,
		pq.Array(&usr.PhoneNumber),
		&usr.Bio,
		&usr.Status,
		&usr.CreatedAt,
		&usr.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	queryAddress := `INSERT INTO addresses (user_id, country, city, district, postalcode)
    VALUES ($1, $2, $3, $4, $5) RETURNING country, city, district, postalcode`
	for _, addr := range user.Addresses {
		err = tx.QueryRow(queryAddress, usr.Id, addr.Country, addr.City, addr.District, addr.PostalCode).Scan(
			&addr.Country, &addr.City, &addr.District, &addr.PostalCode,
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

func (r *UserRepo) UpdateUser(user *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	query := `UPDATE users SET user_name = $1 WHERE id = $2`
	_, err = tx.Exec(query, user.NewUserName, user.Id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return &pb.UpdateUserRes{Update: true}, nil
}

func (r *UserRepo) GetUserById(userID string) (*pb.User, error) {
	var usr pb.User
	query := `SELECT id, first_name, last_name, user_name, 
	email, phone_number, bio, status FROM users WHERE id = $1`
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
	email, phone_number, bio, status FROM users`)
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
