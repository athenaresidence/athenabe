package route

import (
	"github.com/gocroot/lite/controller"
	"github.com/gofiber/fiber/v2"
)

func Web(app *fiber.App) {
	app.Get("/", controller.Homepage)
	app.Post("/webhook/paperka", controller.WebHookPaperka)
	app.Post("/webhook/athena/:waphonenumber", controller.BotAthena)
	//rate
	app.Get("/rate/selfie/:idselfie", controller.GetCekInSelfieData)
	app.Post("/rate/selfie", controller.PostRateSelfie)

}
