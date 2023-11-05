package resources

import (
	"fmt"
	"net/http"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"github.com/gin-gonic/gin"
)

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
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
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

func (s *Server) getSpaceAccessByIDs(c *gin.Context, userID int, spaceID int) (models.SpaceAccess, bool) {
	spaceAccessFromDB, err := s.DB.GetSpaceAccess(userID, spaceID)
	if err != nil {
		fmt.Println("Error getting spaceAccess by id's:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.SpaceAccess{}, false
	}
	if spaceAccessFromDB == (models.SpaceAccess{}) {
		fmt.Println("Error spaceAccess with userID and spaceID not found:", userID, spaceID)
		c.JSON(http.StatusNotFound, models.Error{Error: models.ErrorSpaceAccessWithIdsNotFound})
		return models.SpaceAccess{}, false
	}
	return spaceAccessFromDB, true
}
