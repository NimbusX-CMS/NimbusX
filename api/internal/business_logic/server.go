package business_logic

import (
	"fmt"
	"github.com/NimbusX-CMS/NimbusX/api/internal/db"
	"github.com/NimbusX-CMS/NimbusX/api/internal/error_msg"
	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	DB db.DataBase
}

func (s *Server) GetLogin(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PostLogin(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PostPasswordToken(c *gin.Context, token string) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PostSpace(c *gin.Context) {
	var space models.Space
	if err := c.ShouldBindJSON(&space); err != nil {
		fmt.Println("Error binding space JSON:", err)
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: err.Error()})
		return
	}

	createdSpace, err := s.DB.CreateSpace(space)
	if err != nil {
		fmt.Println("Error creating space:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, createdSpace)
}

func (s *Server) DeleteSpaceSpaceId(c *gin.Context, spaceId int) {
	_, success := s.getSpaceByID(c, spaceId)
	if !success {
		return
	}

	err := s.DB.DeleteSpace(spaceId)
	if err != nil {
		fmt.Println("Error deleting space:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (s *Server) GetSpaceSpaceId(c *gin.Context, spaceId int) {
	space, success := s.getSpaceByID(c, spaceId)
	if !success {
		return
	}
	c.JSON(http.StatusOK, space)
}

func (s *Server) PutSpaceSpaceId(c *gin.Context, spaceId int) {
	var space models.Space
	if err := c.ShouldBindJSON(&space); err != nil {
		fmt.Println("Error binding space JSON:", err)
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: err.Error()})
		return
	}

	space.ID = spaceId

	_, success := s.getSpaceByID(c, spaceId)
	if !success {
		return
	}

	updatedSpace, err := s.DB.UpdateSpace(space)
	if err != nil {
		fmt.Println("Error updating space:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, updatedSpace)
}

func (s *Server) GetSpaces(c *gin.Context) {
	spaces, err := s.DB.GetSpaces()
	if err != nil {
		fmt.Println("Error getting spaces:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, spaces)
}

func (s *Server) PostUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Println("Error binding user JSON:", err)
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: err.Error()})
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
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: error_msg.ErrorEmailAlreadyInUse})
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
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: err.Error()})
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
		c.JSON(http.StatusNotFound, error_msg.Error{Error: error_msg.ErrorUserWithIdNotFound})
		return models.User{}, false
	}
	return userFromDB, true
}

func (s *Server) getSpaceByID(c *gin.Context, spaceID int) (models.Space, bool) {
	spaceFromDB, err := s.DB.GetSpace(spaceID)
	if err != nil {
		fmt.Println("Error getting space by id:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.Space{}, false
	}
	if spaceFromDB.IsEmpty() {
		fmt.Println("Error space with id not found:", spaceID)
		c.JSON(http.StatusNotFound, error_msg.Error{Error: error_msg.ErrorSpaceWithIdNotFound})
		return models.Space{}, false
	}
	return spaceFromDB, true
}

func (s *Server) getSpaceAccessByIDs(c *gin.Context, userID int, spaceID int) (models.SpaceAccess, bool) {
	spaceAccessFromDB, err := s.DB.GetSpaceAccess(userID, spaceID)
	if err != nil {
		fmt.Println("Error getting spaceAccess by id's:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.SpaceAccess{}, false
	}
	if spaceAccessFromDB == (models.SpaceAccess{}) {
		fmt.Println("Error spaceAccess with userID and spaceID not found:", userID, spaceID)
		c.JSON(http.StatusNotFound, error_msg.Error{Error: error_msg.ErrorSpaceAccessWithIdsNotFound})
		return models.SpaceAccess{}, false
	}
	return spaceAccessFromDB, true
}

func (s *Server) GetUserUserIdSpaces(c *gin.Context, userId int) {
	_, success := s.getUserByID(c, userId)
	if !success {
		return
	}
	spaceAccesses, err := s.DB.GetSpaceAccesses(userId)
	if err != nil {
		fmt.Println("Error getting space accesses:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, spaceAccesses)
}

func (s *Server) DeleteUserUserIdSpaceSpaceId(c *gin.Context, userId int, spaceId int) {
	spaceAccess, success := s.getSpaceAccessByIDs(c, userId, spaceId)
	if !success {
		return
	}

	err := s.DB.DeleteSpaceAccess(userId, spaceId)
	if err != nil {
		fmt.Println("Error deleting space access:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, spaceAccess)
}

func (s *Server) PatchUserUserIdSpaces(c *gin.Context, userId int) {
	var spaceAccess models.SpaceAccess
	if err := c.ShouldBindJSON(&spaceAccess); err != nil {
		fmt.Println("Error binding spaceAccess JSON:", err)
		c.JSON(http.StatusBadRequest, error_msg.Error{Error: err.Error()})
		return
	}
	spaceAccess.UserID = userId
	_, success := s.getUserByID(c, userId)
	if !success {
		return
	}
	getSpaceAccess, err := s.DB.GetSpaceAccess(spaceAccess.UserID, spaceAccess.SpaceID)
	if err != nil {
		fmt.Println("Error getting spaceAccess:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if getSpaceAccess == (models.SpaceAccess{}) {
		access, err := s.DB.CreateSpaceAccess(spaceAccess)
		if err != nil {
			fmt.Println("Error creating spaceAccess:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, access)
		return
	}

	access, err := s.DB.UpdateSpaceAccess(spaceAccess)
	if err != nil {
		fmt.Println("Error creating spaceAccess:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, access)
}

func (s *Server) GetWebhooks(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PostWebhooks(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PutWebhooksName(c *gin.Context, name string) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
