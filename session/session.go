package session

import (
	"encoding/gob"
	"kurohelper-api/dto"

	"github.com/gofiber/fiber/v3"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"
)

const (
	storageKey = "user"
	localsKey  = "sessionUser"
)

func init() {
	gob.Register(dto.SessionUser{})
}

func SetUser(c fiber.Ctx, user dto.SessionUser) bool {
	sess := fibersession.FromContext(c)
	if sess == nil {
		return false
	}
	sess.Set(storageKey, user)
	c.Locals(localsKey, user)
	return true
}

func LoadUser(c fiber.Ctx) *dto.SessionUser {
	if user := UserFromLocals(c); user != nil {
		return user
	}

	sess := fibersession.FromContext(c)
	if sess == nil {
		return nil
	}

	raw := sess.Get(storageKey)
	if raw == nil {
		return nil
	}

	user, ok := raw.(dto.SessionUser)
	if !ok {
		return nil
	}

	c.Locals(localsKey, user)
	return &user
}

func UserFromLocals(c fiber.Ctx) *dto.SessionUser {
	raw := c.Locals(localsKey)
	if raw == nil {
		return nil
	}

	user, ok := raw.(dto.SessionUser)
	if !ok {
		return nil
	}

	return &user
}
