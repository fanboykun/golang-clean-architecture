package middleware

import (
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func NewAuth(userUserCase *usecase.UserUseCase) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := SplitBearerToken(ctx.Get("Authorization", "NOT_FOUND"))
		if err != nil {
			return err
		}
		request := &model.VerifyUserRequest{Token: token}
		userUserCase.Log.Debugf("Authorization : %s", request.Token)

		auth, err := userUserCase.Verify(ctx.UserContext(), request)
		if err != nil {
			userUserCase.Log.Warnf("Failed find user by token : %+v", err)
			return fiber.ErrUnauthorized
		}

		userUserCase.Log.Debugf("User : %+v", auth.ID)
		ctx.Locals("auth", auth)
		return ctx.Next()
	}
}

func GetUser(ctx *fiber.Ctx) *model.Auth {
	return ctx.Locals("auth").(*model.Auth)
}

func SplitBearerToken(token string) (string, error) {
	if token == "NOT_FOUND" {
		return "", fiber.ErrUnauthorized
	}
	if len(token) < 7 || token[:7] != "Bearer " {
		return "", fiber.ErrUnauthorized
	}
	splitted := token[7:]
	return splitted, nil
}
