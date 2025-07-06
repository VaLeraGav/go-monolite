package server

import (
	"go-monolite/internal/module/auth"
	"go-monolite/internal/module/category"
	"go-monolite/internal/module/price"
	"go-monolite/internal/module/product"
	"go-monolite/internal/module/property"
	"go-monolite/internal/module/storage"
	"go-monolite/internal/module/user"
	"go-monolite/pkg/logger"
	"go-monolite/pkg/middleware/cors"
	middlewareLogger "go-monolite/pkg/middleware/logger"
	"go-monolite/pkg/middleware/request_id"
	"go-monolite/pkg/middleware/timemiddleware"

	_ "go-monolite/docs"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) ConfigureRouting() {
	s.router.Use(request_id.RequestID)
	s.router.Use(middlewareLogger.New(logger.GetZerologLogger()))
	s.router.Use(timemiddleware.Handler)
	s.router.Use(cors.Handler)

	// Swagger UI
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	s.router.Route("/api", func(r chi.Router) {
		r.Route("/product", product.NewHandler(s.store).Init)
		r.Route("/category", category.NewHandler(s.store).Init)
		r.Route("/property", property.NewHandler(s.store).Init)
		r.Route("/storage", storage.NewHandler(s.store).Init)
		r.Route("/price", price.NewHandler(s.store).Init)

		r.Route("/auth", auth.NewHandler(s.store).Init)
		r.Route("/user", user.NewHandler(s.store).Init)
	})
}
