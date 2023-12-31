package controller

import (
	"fmt"
	"net/http"
	"notes/lib"
	"notes/model"
	"notes/repository"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthController struct {
	Repos repository.AuthRepository
}

type MetaAndAssets struct {
	Title  string `json:"title"`
	Style  string `json:"style"`
	Script string `json:"script"`
	Popper string `json:"popper"`
}

const (
	BootstrapCSS = "https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css"
	BootstrapJS  = "https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"
	PopperJS     = "https://cdn.jsdelivr.net/npm/@popperjs/core@2.11.8/dist/umd/popper.min.js"
)

var metaAndAssets MetaAndAssets = MetaAndAssets{
	Style:  BootstrapCSS,
	Script: BootstrapJS,
	Popper: PopperJS,
}

func NewAuthController(repos repository.AuthRepository) *AuthController {
	return &AuthController{
		Repos: repos,
	}
}

func (controller *AuthController) Routes(r *gin.Engine) {
	r.GET("/signup", controller.SignUpView)
	r.GET("/me", controller.NotesView)
	r.POST("/signup", controller.SignUp)
}

func (controller *AuthController) SignUpView(ctx *gin.Context) {
	metaAndAssets.Title = "Sign Up - Notes"

	ctx.HTML(http.StatusOK, "signup.html", metaAndAssets)
}

func (controller *AuthController) NotesView(ctx *gin.Context) {
	metaAndAssets.Title = "My Notes"

	s := sessions.Default(ctx)
	s.Get("user_email")
	fmt.Println("User Email:", s)

	ctx.HTML(http.StatusOK, "me.html", metaAndAssets)
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	var payload model.AuthRequest

	err := ctx.ShouldBind(&payload)

	if err != nil {
		fmt.Println("Request should be in JSON format!", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Payload should be in JSON",
		})
		return
	}

	// Check if user is already registered
	user, _ := controller.Repos.FindUser(payload.Email)

	if user != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failure",
			"message": "Email sudah terdaftar",
		})

		return
	}

	// If not yet, then user can register

	// We need to hash password before insert into database.
	hash := lib.HashPassword(payload.Password)

	// Generate unique ID
	id := uuid.New().String()

	err = controller.Repos.CreateUser(&model.User{
		ID:       id,
		Email:    payload.Email,
		Password: hash,
		Name:     payload.Name,
		Role:     0,
	})

	// Internal logging
	if err != nil {
		fmt.Println("Failed to create user", err)
		return
	}

	sessions := sessions.Default(ctx)
	sessions.Set("user_email", payload.Email)

	// if success, we need redirect it into Home.
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Berhasil mendaftar user baru",
	})
}

func (controller *AuthController) SetupSession(redis redis.Store, app *gin.Engine) {
	app.Use(sessions.Sessions("auth_session", redis))
}
