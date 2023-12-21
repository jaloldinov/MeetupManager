package router

import (
	"meetup/internal/auth"
	auth_controller "meetup/internal/controller/http/v1/auth"
	person_controller "meetup/internal/controller/http/v1/person"
	user_controller "meetup/internal/controller/http/v1/user"

	"meetup/internal/pkg/repository/postgres"
	person_repo "meetup/internal/repository/postgres/person"
	user_repo "meetup/internal/repository/postgres/user"

	"github.com/gin-gonic/gin"
)

type Auth interface {
	HasPermission(roles ...string) gin.HandlerFunc
}

type AuthController interface {
	SignIn(c *gin.Context)
}

type Router struct {
	postgresDB *postgres.Database
	auth       *auth.Auth
}

func New(auth *auth.Auth, postgresDB *postgres.Database) *Router {
	return &Router{
		auth:       auth,
		postgresDB: postgresDB,
	}
}

func (r *Router) Init(port string) error {
	router := gin.Default()

	//repository
	userRepo := user_repo.NewRepository(r.postgresDB)
	personRepo := person_repo.NewRepository(r.postgresDB)

	//controller
	authController := auth_controller.NewController(userRepo, r.auth)

	userController := user_controller.NewController(userRepo)
	personController := person_controller.NewController(personRepo)

	// #AUTH
	router.POST("/api/v1/admin/sign-in", authController.SignIn)

	// # ADMIN USER
	router.POST("api/v1/admin/user/create", r.auth.HasPermission("ADMIN"), userController.CreateUser)
	router.GET("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.GetUserById)
	router.GET("/api/v1/admin/user/list", r.auth.HasPermission("ADMIN"), userController.GetUserList)
	router.PUT("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.UpdateUser)
	router.DELETE("/api/v1/admin/user/:id", r.auth.HasPermission("ADMIN"), userController.DeleteUser)
	router.PUT("/api/v1/admin/user/reset/password", r.auth.HasPermission("ADMIN"), userController.ResetUserPassword)

	// # ADMIN PERSON
	router.POST("api/v1/admin/person/create", r.auth.HasPermission("ADMIN"), personController.CreatePerson)
	router.GET("/api/v1/admin/person/:id", r.auth.HasPermission("ADMIN"), personController.GetPersonById)
	router.GET("/api/v1/admin/person/list", r.auth.HasPermission("ADMIN"), personController.GetPersonList)
	router.PUT("/api/v1/admin/person/:id", r.auth.HasPermission("ADMIN"), personController.UpdatePerson)
	router.DELETE("/api/v1/admin/person/:id", r.auth.HasPermission("ADMIN"), personController.DeletePerson)

	return router.Run(port)
}
