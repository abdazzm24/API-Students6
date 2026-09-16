package service

import (
	"strconv"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type UserService struct {
	repo  repository.UserRepository
	perms *helper.PermissionSet
}

func NewUserService(
	repo repository.UserRepository,
	perms *helper.PermissionSet,
) *UserService {

	return &UserService{
		repo:  repo,
		perms: perms,
	}
}

// ========================================================
// GET /users
// ========================================================

func (s *UserService) List(
	c *fiber.Ctx,
) error {

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	users, err := s.repo.List(ctx)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil daftar user",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"daftar user berhasil diambil",
		users,
	)
}

// ========================================================
// GET /users/:id
// ========================================================

func (s *UserService) Get(
	c *fiber.Ctx,
) error {

	current, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			nil,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	user, err := s.repo.FindByID(
		ctx,
		id,
	)

	if err != nil {
		if err == repository.ErrNotFound {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"user tidak ditemukan",
				nil,
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil user",
			nil,
		)
	}

	if !CanAccessUser(
		current,
		id,
		s.perms,
		"user:read:any",
	) {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak berhak mengakses data user lain",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user ditemukan",
		user,
	)
}

// ========================================================
// POST /users
// ========================================================

func (s *UserService) Create(
	c *fiber.Ctx,
) error {

	var req model.CreateUserRequest

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

	registerReq := model.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if errs := ValidateRegister(registerReq); len(errs) > 0 {
		return helper.FailValidation(
			c,
			errs,
		)
	}

	passwordHash, err := helper.HashPassword(
		req.Password,
	)

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
		Password: passwordHash,

		// Role selalu ditentukan server.
		Role: "user",

		IsActive: true,
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	created, err := s.repo.Create(
		ctx,
		user,
	)

	if err != nil {
		if err == repository.ErrDuplicate {
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
			"gagal membuat user",
			nil,
		)
	}

	return helper.Created(
		c,
		"user berhasil dibuat",
		created,
		"/api/v1/users/"+strconv.Itoa(created.ID),
	)
}

// ========================================================
// PUT /users/:id
// ========================================================

func (s *UserService) Replace(
	c *fiber.Ctx,
) error {

	current, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			nil,
		)
	}

	if !CanAccessUser(
		current,
		id,
		s.perms,
		"user:update:any",
	) {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak berhak mengubah user lain",
			nil,
		)
	}

	var req model.ReplaceUserRequest

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

	if req.Username == "" {
		return helper.FailValidation(
			c,
			map[string]string{
				"username": "wajib diisi",
			},
		)
	}

	if !isValidEmail(req.Email) {
		return helper.FailValidation(
			c,
			map[string]string{
				"email": "format email tidak valid",
			},
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	user, err := s.repo.Replace(
		ctx,
		id,
		req,
	)

	if err != nil {
		if err == repository.ErrNotFound {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"user tidak ditemukan",
				nil,
			)
		}

		if err == repository.ErrDuplicate {
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
			"gagal memperbarui user",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user berhasil diperbarui",
		user,
	)
}

// ========================================================
// PATCH /users/:id
// ========================================================

func (s *UserService) Patch(
	c *fiber.Ctx,
) error {

	current, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			nil,
		)
	}

	if !CanAccessUser(
		current,
		id,
		s.perms,
		"user:update:any",
	) {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak berhak mengubah user lain",
			nil,
		)
	}

	var req model.PatchUserRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	if req.Username != nil {
		value := strings.TrimSpace(*req.Username)

		if value == "" {
			return helper.FailValidation(
				c,
				map[string]string{
					"username": "tidak boleh kosong",
				},
			)
		}

		req.Username = &value
	}

	if req.Email != nil {
		value := strings.TrimSpace(*req.Email)

		if !isValidEmail(value) {
			return helper.FailValidation(
				c,
				map[string]string{
					"email": "format email tidak valid",
				},
			)
		}

		req.Email = &value
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	user, err := s.repo.Patch(
		ctx,
		id,
		req,
	)

	if err != nil {
		if err == repository.ErrNotFound {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"user tidak ditemukan",
				nil,
			)
		}

		if err == repository.ErrDuplicate {
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
			"gagal memperbarui user",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"user berhasil diperbarui",
		user,
	)
}

// ========================================================
// DELETE /users/:id
// ========================================================

func (s *UserService) Delete(
	c *fiber.Ctx,
) error {

	current, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			nil,
		)
	}

	// Admin tetap tidak boleh menghapus dirinya sendiri.
	if current.UserID == id {
		return helper.Fail(
			c,
			fiber.StatusForbidden,
			"tidak boleh menghapus akun sendiri",
			nil,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	if err := s.repo.Delete(
		ctx,
		id,
	); err != nil {

		if err == repository.ErrNotFound {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"user tidak ditemukan",
				nil,
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal menghapus user",
			nil,
		)
	}

	return helper.NoContent(c)
}

// ========================================================
// PATCH /users/:id/role
// ========================================================

func (s *UserService) AssignRole(
	c *fiber.Ctx,
) error {

	current, ok := helper.CurrentUser(c)

	if !ok {
		return helper.Fail(
			c,
			fiber.StatusUnauthorized,
			"belum terautentikasi",
			nil,
		)
	}

	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			nil,
		)
	}

	var req model.AssignRoleRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body harus berupa JSON yang valid",
			nil,
		)
	}

	if errs := ValidateAssignRole(
		current,
		id,
		req,
		s.perms,
	); len(errs) > 0 {

		return helper.FailValidation(
			c,
			errs,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	user, err := s.repo.UpdateRole(
		ctx,
		id,
		strings.TrimSpace(req.Role),
	)

	if err != nil {
		if err == repository.ErrNotFound {
			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"user tidak ditemukan",
				nil,
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengubah role user",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"role user berhasil diubah",
		user,
	)
}