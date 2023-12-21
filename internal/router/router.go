package router

import (
	"meetup/internal/auth"
	auth_controller "meetup/internal/controller/http/v1/auth"
	meeting_controller "meetup/internal/controller/http/v1/meeting"
	person_controller "meetup/internal/controller/http/v1/person"
	place_controller "meetup/internal/controller/http/v1/place"
	user_controller "meetup/internal/controller/http/v1/user"

	"meetup/internal/pkg/repository/postgres"
	meeting_repo "meetup/internal/repository/postgres/meeting"
	person_repo "meetup/internal/repository/postgres/person"
	place_repo "meetup/internal/repository/postgres/place"
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
	placeRepo := place_repo.NewRepository(r.postgresDB)
	meetingRepo := meeting_repo.NewRepository(r.postgresDB)

	//controller
	authController := auth_controller.NewController(userRepo, r.auth)

	userController := user_controller.NewController(userRepo)
	personController := person_controller.NewController(personRepo)
	placeController := place_controller.NewController(placeRepo)
	meetingController := meeting_controller.NewController(meetingRepo)

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

	// # ADMIN PLACE
	router.POST("api/v1/admin/place/create", r.auth.HasPermission("ADMIN"), placeController.CreatePlace)
	router.GET("/api/v1/admin/place/:id", r.auth.HasPermission("ADMIN"), placeController.GetPlaceById)
	router.GET("/api/v1/admin/place/list", r.auth.HasPermission("ADMIN"), placeController.GetPlaceList)
	router.PUT("/api/v1/admin/place/:id", r.auth.HasPermission("ADMIN"), placeController.UpdatePlace)
	router.DELETE("/api/v1/admin/place/:id", r.auth.HasPermission("ADMIN"), placeController.DeletePlace)

	// # ADMIN MEETING
	router.POST("api/v1/admin/meeting/create", r.auth.HasPermission("ADMIN"), meetingController.CreateMeeting)
	router.GET("/api/v1/admin/meeting/:id", r.auth.HasPermission("ADMIN"), meetingController.GetMeetingById)
	router.GET("/api/v1/admin/meeting/list", r.auth.HasPermission("ADMIN"), meetingController.GetMeetingList)
	router.PUT("/api/v1/admin/meeting/:id", r.auth.HasPermission("ADMIN"), meetingController.UpdateMeeting)
	router.DELETE("/api/v1/admin/meeting/:id", r.auth.HasPermission("ADMIN"), meetingController.DeleteMeeting)

	return router.Run(port)
}
