package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sso/internal/domains/models"
	"sso/internal/lib/jwt"
	"sso/storage"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
	appProvider  AppProvider
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		login string,
		passhash []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	UserByLogin(ctx context.Context, login string) (models.User, error)
}

type AppProvider interface {
	AppById(ctx context.Context, appId int64) (models.App, error)
}

// Return Auth service layer
func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *Auth {
	return &Auth{
		log,
		userSaver,
		userProvider,
		appProvider,
	}
}

func (auth *Auth) Login(
	ctx context.Context,
	login string,
	password string,
	appId int64,
) (string, error) {
	const op = "service.auth.Login"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("login", login), // todo
		slog.Int64("appId", appId),
	)

	log.Info("Login user")

	user, err := auth.userProvider.UserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("User not found", slog.String("err", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to found user", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInternal)
	}

	if err := bcrypt.CompareHashAndPassword(user.Passhash, []byte(password)); err != nil {
		log.Info("Invalid credentials", slog.String("err", err.Error()))
		return "", fmt.Errorf("%s: %w", op, ErrInternal)
	}

	app, err := auth.appProvider.AppById(ctx, appId)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("App not found", slog.String("err", err.Error()))
			return "", ErrAppNotFound
		}

		log.Error("Failed to found app", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	token, err := jwt.New(user, app, time.Hour) // todo: duration supposed to be provided
	if err != nil {
		log.Error("Failed to generate token", slog.String("err", err.Error()))
		return "", ErrInternal
	}

	return token, nil
}

func (auth *Auth) Register(
	ctx context.Context,
	login string,
	password string,
) (uid int64, err error) {
	const op = "service.auth.Register"

	log := auth.log.With(
		slog.String("op", op),
		slog.String("login", login), // todo
	)

	log.Info("Registering user")

	passhash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to accept password", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	id, err := auth.userSaver.SaveUser(ctx, login, passhash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Info("User already exists", slog.String("err", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}

		log.Error("Failed to save user", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, ErrInternal)
	}

	log.Info("Registered", slog.Int64("id", id))

	return id, nil
}
