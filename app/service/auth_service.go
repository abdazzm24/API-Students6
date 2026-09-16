package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

const refreshTokenBytes = 32

type AuthService struct {
	users      repository.UserRepository
	tokens     repository.TokenRepository
	jwt        *helper.JWTManager
	perms      *helper.PermissionSet
	refreshTTL time.Duration
}

func NewAuthService(
	users repository.UserRepository,
	tokens repository.TokenRepository,
	jwtManager *helper.JWTManager,
	perms *helper.PermissionSet,
	refreshTTL time.Duration,
) *AuthService {

	return &AuthService{
		users:      users,
		tokens:     tokens,
		jwt:        jwtManager,
		perms:      perms,
		refreshTTL: refreshTTL,
	}
}

func (s *AuthService) Register(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if errs := ValidateRegister(req); len(errs) > 0 {
		return helper.Fail(
			c,
			fiber.StatusUnprocessableEntity,
			"validasi gagal",
			errs,
		)
	}

	hashed, err := helper.HashPassword(req.Password)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal memproses password",
			nil,
		)
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,

		// Role ditentukan server.
		Role: "user",

		IsActive: true,
	}

	created, err := s.users.Create(ctx, user)

	if err != nil {

		if errors.Is(err, repository.ErrDuplicate) {
			return helper.Fail(
				c,
				fiber.StatusConflict,
				"username atau email sudah digunakan",
				nil,
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mendaftarkan user",
			nil,
		)
	}

	return helper.Created(
		c,
		"pendaftaran berhasil",
		created,
		"/api/v1/auth/me",
	)
}

func (s *AuthService) Login(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	if errs := ValidateLogin(req); len(errs) > 0 {
		return helper.Fail(
			c,
			fiber.StatusUnprocessableEntity,
			"validasi gagal",
			errs,
		)
	}

	username := strings.TrimSpace(req.Username)

	user, err := s.users.FindByUsername(
		ctx,
		username,
	)

	if err != nil {

		// Tetap melakukan bcrypt comparison
		// untuk mencegah user enumeration
		// melalui perbedaan waktu response.
		helper.VerifyDummyPassword(req.Password)

		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"username atau password salah",
			nil,
		)
	}

	if !helper.VerifyPassword(
		user.Password,
		req.Password,
	) {

		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"username atau password salah",
			nil,
		)
	}

	if !user.IsActive {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"akun dinonaktifkan",
			nil,
		)
	}

	pair, err := s.issueTokenPair(
		ctx,
		user,
	)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat token",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"login berhasil",
		pair,
	)
}

func (s *AuthService) Refresh(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	if strings.TrimSpace(req.RefreshToken) == "" {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"refresh_token wajib diisi",
			nil,
		)
	}

	hash := helper.SHA256Hex(
		req.RefreshToken,
	)

	stored, err := s.tokens.FindActive(
		ctx,
		hash,
	)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"refresh token tidak valid atau sudah kedaluwarsa",
			nil,
		)
	}

	user, err := s.users.FindByID(
		ctx,
		stored.UserID,
	)

	if err != nil || !user.IsActive {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"akun tidak dapat dipakai",
			nil,
		)
	}

	// Rotasi refresh token.
	if err := s.tokens.Revoke(
		ctx,
		hash,
	); err != nil {

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal memperbarui token",
			nil,
		)
	}

	pair, err := s.issueTokenPair(
		ctx,
		user,
	)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat token",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"token berhasil diperbarui",
		pair,
	)
}

func (s *AuthService) Logout(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	var req model.RefreshRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	if strings.TrimSpace(req.RefreshToken) != "" {

		_ = s.tokens.Revoke(
			ctx,
			helper.SHA256Hex(
				req.RefreshToken,
			),
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"logout berhasil",
		nil,
	)
}

func (s *AuthService) Me(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	authUser, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	user, err := s.users.FindByID(
		ctx,
		authUser.UserID,
	)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"user tidak ditemukan",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"profil berhasil diambil",
		fiber.Map{
			"user":        user,
			"permissions": s.perms.PermissionsOf(user.Role),
		},
	)
}

func (s *AuthService) issueTokenPair(
	ctx context.Context,
	user model.User,
) (model.TokenPair, error) {

	accessToken, err := s.jwt.GenerateAccess(
		user,
	)

	if err != nil {
		return model.TokenPair{}, err
	}

	refreshToken, err := helper.RandomToken(
		refreshTokenBytes,
	)

	if err != nil {
		return model.TokenPair{}, err
	}

	err = s.tokens.Save(
		ctx,
		model.RefreshToken{
			UserID: user.ID,

			TokenHash: helper.SHA256Hex(
				refreshToken,
			),

			ExpiresAt: time.Now().Add(
				s.refreshTTL,
			),
		},
	)

	if err != nil {
		return model.TokenPair{}, err
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn: int(
			s.jwt.AccessTTL().Seconds(),
		),
	}, nil
}