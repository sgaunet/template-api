package webserver

// "github.com/go-redis/redis/v7"

func (w *WebServer) initRoutes() {
	w.app.Post("/authors", w.authorsSvc.Create)
	w.app.Get("/authors/:id", w.authorsSvc.Get)
	w.app.Put("/authors/:id", w.authorsSvc.FullUpdate)
	w.app.Delete("/authors/:id", w.authorsSvc.Delete)
	w.app.Get("/authors", w.authorsSvc.List)
}
