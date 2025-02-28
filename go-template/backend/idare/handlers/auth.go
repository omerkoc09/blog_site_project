package handlers

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/gomail.v2"

	"github.com/hayrat/go-template2/backend/common/service"
	"github.com/hayrat/go-template2/backend/idare/viewmodel"
	"github.com/hayrat/go-template2/backend/pkg/app"
	"github.com/hayrat/go-template2/backend/pkg/errorsx"
	"github.com/hayrat/go-template2/backend/pkg/utils"
)

type AuthHandler struct {
	authService  service.IAuthService
	userService  service.IUserService
	tokenService service.TokenService
}

func NewAuthHandler(s service.IAuthService, us service.IUserService, ts service.TokenService) AuthHandler {
	h := AuthHandler{
		authService:  s,
		userService:  us,
		tokenService: ts,
	}

	return h
}

func (h AuthHandler) Login(ctx *app.Ctx) error {
	var vm viewmodel.AuthLoginVM
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	user, err := h.userService.GetByEmail(ctx.Context(), utils.EmailTemizle(vm.Email))
	if err != nil {
		return err
	}

	ok := utils.CheckPasswordHash(strings.TrimSpace(vm.Password), user.Password)
	if !ok {
		return errorsx.UnauthorizedError("Hatalı Email veya Parola")
	}

	refreshTokenID := uuid.New()
	tokens, err := h.authService.GenerateTokenPair(user.ID, refreshTokenID)
	if err != nil {
		return errorsx.InternalError(err)
	}

	err = h.authService.CreateAuthRefreshToken(ctx.Context(), refreshTokenID, user.ID)
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return ctx.SuccessResponse(result)
}

func (s AuthHandler) TokenForForgotPassword(ctx *app.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}

	var req Request
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	user, err := s.userService.GetByEmail(ctx.Context(), req.Email)
	if err != nil {
		return err
	}

	refreshTokenID := uuid.New()
	tokens, err := s.authService.GenerateTokenPair(user.ID, refreshTokenID)
	if err != nil {
		return errorsx.InternalError(err)
	}

	err = s.authService.CreateAuthRefreshToken(ctx.Context(), refreshTokenID, user.ID)
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}

	sendResetEmail(user.Email, result.AccessToken)

	return ctx.SuccessResponse(result)
}

func (h AuthHandler) RefreshToken(ctx *app.Ctx) error {
	var vm viewmodel.AuthRefreshVM
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	refreshTokenID, userID, err := h.authService.ParseRefreshToken(vm.RefreshToken)
	if err != nil {
		return errorsx.UnauthorizedError(err.Error())
	}
	authRefreshToken, err := h.authService.GetAuthRefreshToken(ctx.Context(), refreshTokenID)
	if err != nil {
		return err
	}
	newTokenPair, err := h.authService.GenerateTokenPair(userID, refreshTokenID)
	if err != nil {
		return errorsx.InternalError(err)
	}

	err = h.authService.UpdateAuthRefreshTokenExpires(ctx.Context(), authRefreshToken)
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
	}
	return ctx.SuccessResponse(result)
}

func (h AuthHandler) Logout(ctx *app.Ctx) error {
	err := h.authService.DeleteAuthRefreshToken(ctx.Context(), ctx.GetUserID())
	if err != nil {
		return err
	}

	return nil
}

func sendResetEmail(email, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ofk6409@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset Request")

	resetLink := "http://localhost:5173/reset-password?token=" + token

	// Email body
	body := `
        Hello,
        
        You requested a password reset.
        Click this link to reset your password: ` + resetLink + `
        
        This link will expire in 24 hours.
        
        If you didn't request this, please ignore this email.`

	m.SetBody("text/plain", body)

	// Configure SMTP settings
	dialer := gomail.NewDialer(
		"smtp.gmail.com",      // SMTP server
		587,                   // Port
		"ofk6409@gmail.com",   // Username
		"mmnh aigu hqzn fqoj", // Password
	)

	return dialer.DialAndSend(m)
}

func (h AuthHandler) ResetPassword(ctx *app.Ctx) error {
	var req struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,min=8"`
	}
	if errs := ctx.BodyParseValidate(&req); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	claims, err := h.tokenService.ParseToken(req.Token)
	if err != nil {
		return err
	}

	user, err := h.userService.GetByID(ctx.Context(), int64(claims.UserID))
	if err != nil {
		return err
	}

	// Update password
	user.Password = req.Password
	user.Password, _ = utils.HashPassword(user.Password)

	err = h.userService.Update(ctx.Context(), user)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, "Password updated successfully")
}
