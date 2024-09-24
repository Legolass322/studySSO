package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domains/models"
)

func (s *Database) SaveUser(ctx context.Context, login string, passhash []byte) (uid int64, err error) {
	const op = "db.provider.SaveUser"

	log := s.log.With(slog.String("op", op))

	log.Info("Start request")

	stmt, err := s.db.Prepare("INSERT INTO users(login, passhash) VALUES(?, ?)")
	if err != nil {
		log.Error("Cannot prepare statement", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, login, passhash)
	if err != nil {
		// todo: make check for existing user
		log.Error("Cannot exec", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error("Cannot get id", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Request success")

	return id, nil
}

func (s *Database) UserByLogin(ctx context.Context, login string) (models.User, error) {
	const op = "db.provider.UserByLogin"

	log := s.log.With(slog.String("op", op))

	log.Info("Start request")

	stmt, err := s.db.Prepare("SELECT id, login, passhash FROM users WHERE login = ?")
	if err != nil {
		log.Error("Cannot prepare statement", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User
	row := stmt.QueryRowContext(ctx, login)
	if err := row.Scan(&user.Id, &user.Login, &user.Passhash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No such user", slog.String("err", err.Error()))
		} else {
			log.Error("Internal error", slog.String("err", err.Error()))
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Request success")

	return user, nil
}

func (s *Database) AppById(ctx context.Context, id int64) (models.App, error) {
	const op = "db.provider.AppById"

	log := s.log.With(slog.String("op", op))

	log.Info("Start request")

	stmt, err := s.db.Prepare("SELECT id, name FROM apps WHERE id = ?")
	if err != nil {
		log.Error("Cannot prepare statement", slog.String("err", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	var app models.App
	row := stmt.QueryRowContext(ctx, id)
	if err := row.Scan(&app.Id, &app.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Info("No such app", slog.String("err", err.Error()))
		} else {
			log.Error("Internal error", slog.String("err", err.Error()))
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Request success")

	return app, nil
}
