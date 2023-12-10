package webserver

import (
	"database/sql"
	"fmt"

	// "github.com/go-redis/redis/v7"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sgaunet/template-api/api/authors"
	"github.com/sgaunet/template-api/internal/database"
	"github.com/sgaunet/template-api/pkg/config"
)

type WebServer struct {
	db         *sql.DB
	authorsSvc *authors.Service
	// redisClient *redis.Client
	app *fiber.App
	cfg *config.Config
}

func NewWebServer(cfg *config.Config) (*WebServer, error) {
	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		// StrictRouting: true,
		// ServerHeader:  "Fiber",
		// AppName:       "Test App v1.0.1",
	})

	// Logging Request ID
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		// For more options, see the Config section
		Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
	}))

	// app.Use(limiter.New(limiter.Config{
	// 	Max:        20,
	// 	Expiration: 30 * time.Second,
	// 	KeyGenerator: func(c *fiber.Ctx) string {
	// 		return c.Get("x-forwarded-for")
	// 	},
	// 	// LimitReached: func(c *fiber.Ctx) error {
	// 	// 	return c.SendFile("./toofast.html")
	// 	// },
	// }))
	// app.Static("/bootstrap", "./static/bootstrap")
	w := &WebServer{
		app: app,
		cfg: cfg,
	}
	// w.PublicRoutes()
	err := w.initRepository(cfg.DbDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to init database connection: %w", err)
	}
	// err = w.initRedisConnection(cfg.RedisDSN)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to init redis connection: %w", err)
	// }
	// Add routes
	w.initRoutes()
	return w, nil
}

func (w *WebServer) initRepository(dbUrl string) (err error) {
	pg, err := database.NewPostgres(dbUrl)
	if err != nil {
		return err
	}
	w.db = pg.DB
	queries := database.New(w.db)
	w.authorsSvc = authors.NewService(queries)
	return nil
}

// func (w *WebServer) initRedisConnection(redisdsn string) error {
// 	var err error
// 	d, err := dsn.New(redisdsn)
// 	if err != nil {
// 		return err
// 	}
// 	addr := fmt.Sprintf("%s:%s", d.GetHost(), d.GetPort("6379"))
// 	// fmt.Println("conn to ", addr)
// 	w.redisClient = redis.NewClient(&redis.Options{
// 		Addr: addr,
// 	})
// 	_, err = w.redisClient.Ping().Result()
// 	return err
// }

func (w *WebServer) Start() error {
	return w.app.Listen(":3000")
}

func (w *WebServer) Shutdown() error {
	return w.app.Shutdown()
}
