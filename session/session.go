package session

import (
	"encoding/gob"
	"time"

	"github.com/gofiber/fiber/v3"
	fibersession "github.com/gofiber/fiber/v3/middleware/session"

	"kurohelperservice/db"
)

const (
	storageKey = "user"
	localsKey  = "sessionUser"
)

// User is the logged-in user snapshot stored in session (not an API contract).
type User struct {
	ID          int
	UserName    string
	NickName    string
	DiscordID   string
	Avatar      string
	Description string
	Role        int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func init() {
	gob.Register(User{})
}

func SetUser(c fiber.Ctx, user User) bool {
	sess := fibersession.FromContext(c)
	if sess == nil {
		return false
	}
	sess.Set(storageKey, user)
	c.Locals(localsKey, user)
	return true
}

func LoadUser(c fiber.Ctx) *User {
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

	user, ok := raw.(User)
	if !ok {
		return nil
	}

	c.Locals(localsKey, user)
	return &user
}

func UserFromLocals(c fiber.Ctx) *User {
	raw := c.Locals(localsKey)
	if raw == nil {
		return nil
	}

	user, ok := raw.(User)
	if !ok {
		return nil
	}

	return &user
}

func discordIDOrEmpty(discordID *string) string {
	if discordID == nil {
		return ""
	}
	return *discordID
}

func NewUser(u db.User, userName string) User {
	return User{
		ID:          u.ID,
		UserName:    userName,
		NickName:    u.Name,
		DiscordID:   discordIDOrEmpty(u.DiscordID),
		Avatar:      u.Avatar,
		Description: u.Description,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
