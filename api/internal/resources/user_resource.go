package resources

import (
	"fmt"
	"net/http"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error binding user JSON:", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}
	userByEmail, err := s.DB.GetUserByEmail(user.Email)
	if err != nil {
		fmt.Println("Error getting user by email:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userByEmail != (models.User{}) {
		fmt.Println("Error email already in use:", userByEmail)
		c.JSON(http.StatusBadRequest, models.Error{Error: models.ErrorEmailAlreadyInUse})
		return
	}

	createdUser, err := s.DB.CreateUser(user)
	if err != nil {
		fmt.Println("Error creating user:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

func (s *Server) DeleteUserUserId(c *gin.Context, userId int) {
	_, success := s.getUserByID(c, userId)
	if !success {
		return
	}

	err := s.DB.DeleteUser(userId)
	if err != nil {
		fmt.Println("Error deleting user:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetUserUserId(c *gin.Context, userId int) {
	user, success := s.getUserByID(c, userId)
	if !success {
		return
	}
	c.JSON(http.StatusOK, user)
}

func (s *Server) PutUserUserId(c *gin.Context, userId int) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error binding user JSON:", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
		return
	}
	user.ID = userId

	_, success := s.getUserByID(c, userId)
	if !success {
		return
	}

	updatedUser, err := s.DB.UpdateUser(user)
	if err != nil {
		fmt.Println("Error updating user:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (s *Server) GetUsers(c *gin.Context) {
	users, err := s.DB.GetUsers()
	if err != nil {
		fmt.Println("Error getting users:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, users)
}

func (s *Server) getUserByID(c *gin.Context, userId int) (models.User, bool) {
	userFromDB, err := s.DB.GetUser(userId)
	if err != nil {
		fmt.Println("Error getting user by id:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.User{}, false
	}
	if userFromDB == (models.User{}) {
		fmt.Println("Error user with id not found:", userId)
		c.JSON(http.StatusNotFound, models.Error{Error: models.ErrorUserWithIdNotFound})
		return models.User{}, false
	}
	return userFromDB, true
}
