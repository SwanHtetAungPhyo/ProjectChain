package http

import (
	"github.com/SwanHtetAungPhyo/common/model"
	_ "github.com/SwanHtetAungPhyo/ledchain/cmd/http/docs"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"time"
)

var log = logrus.New()

func Start(httpPort string) {
	app := fiber.New()

	log.Info("App started at the port http://localhost:" + httpPort)
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Get("/", func(c *fiber.Ctx) error {
		claims := jwt.MapClaims{
			"name":  "John Doe",
			"admin": true,
			"exp":   time.Now().Add(time.Hour * 72).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.Status(200).JSON(fiber.Map{
			"token": t,
		})

	})
	routeSetup(app)
	if err := app.Listen(":" + httpPort); err != nil {
		log.Fatal(err)
	}

}

func routeSetup(app *fiber.App) {
	// Correct usage of JWT middleware
	chainGroup := app.Group("/chain")
	chainGroup.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte("secret")},
	}))

	// Ensure you're passing the right handler functions
	handlerr := HandlerImpl{}
	chainGroup.Post("/wallet", handlerr.CreateAccount)
	chainGroup.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    model.SwanDAG,
		})
	})
	chainGroup.Post("/trans", handlerr.ExecuteTransaction)
}
