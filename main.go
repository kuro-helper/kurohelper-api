package main

import (
	"fmt"
	"kurohelper-api/middlware"
	"kurohelper-api/router"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/lmittmann/tint"
	slogmulti "github.com/samber/slog-multi"

	kurohelperdb "kurohelperservice/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Error(".env file load error", "err", err)
		os.Exit(1)
	}

	stdHandler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Stamp,
	})

	logDir := strings.TrimSpace(os.Getenv("LOG_PATH"))
	if logDir != "" {
		logDir = filepath.Clean(logDir)
		info, err := os.Stat(logDir)
		if err != nil || !info.IsDir() {
			slog.SetDefault(slog.New(stdHandler))
			if err != nil {
				slog.Warn("LOG_PATH stat failed; stdout only", "path", logDir, "err", err)
			} else {
				slog.Warn("LOG_PATH is not a directory; stdout only", "path", logDir)
			}
			return
		}

		logFile, err := os.OpenFile(filepath.Join(logDir, "kurohelper-api-nocolor.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			slog.SetDefault(slog.New(stdHandler))
			slog.Warn("log file open failed; stdout only", "path", logDir, "err", err)
			return
		}

		fileHandler := slog.NewTextHandler(logFile, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					return slog.String(slog.TimeKey, a.Value.Time().Format(time.Stamp))
				}
				return a
			},
		})
		slog.SetDefault(slog.New(slogmulti.Fanout(stdHandler, fileHandler)))
		return
	}

	slog.SetDefault(slog.New(stdHandler))
}

func main() {
	config := kurohelperdb.Config{
		DBOwner:    os.Getenv("DB_OWNER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
	}

	kurohelperdb.InitDsn(config)
	// kurohelperdb.Migration(kurohelperdb.Dbs) // 選填

	initTokenCache()

	allowOrigins := parseAllowOrigins(os.Getenv("ALLOW_CORS"))

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// api route group
	apiGroup := app.Group("/api") // main api route group

	// site route group
	// routes.TokenRouter(apiGroup)
	router.UserDataRouter(apiGroup)
	router.UserRouter(apiGroup)
	// routes.SearchRouter(apiGroup)

	addr := fmt.Sprintf("127.0.0.1:%s", os.Getenv("PRODUCTION_PORT"))
	slog.Info("fiber open...")
	if err := app.Listen(addr); err != nil {
		slog.Error("fiber listen", "addr", addr, "err", err)
		os.Exit(1)
	}
}

// parseAllowOrigins splits ALLOW_CORS (comma-separated) into a slice for cors.Config.
func parseAllowOrigins(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var out []string
	for part := range strings.SplitSeq(raw, ",") {
		if o := strings.TrimSpace(part); o != "" {
			out = append(out, o)
		}
	}
	return out
}

func initTokenCache() {
	webAPIToken, err := kurohelperdb.GetWebAPIToken(kurohelperdb.Dbs)
	if err != nil {
		slog.Error("init token cache", "err", err)
		os.Exit(1)
	}

	for _, t := range webAPIToken {
		middlware.VaildToken[t.ID] = struct{}{}
	}
}
