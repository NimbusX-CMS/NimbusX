package resources

import (
	"fmt"
	"net/http"

	"github.com/NimbusX-CMS/NimbusX/api/internal/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) PostSpace(c *gin.Context) {
	var space models.Space
	if err := c.ShouldBindJSON(&space); err != nil {
		fmt.Println("Error binding space JSON:", err)
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
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
		c.JSON(http.StatusBadRequest, models.Error{Error: err.Error()})
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

func (s *Server) getSpaceByID(c *gin.Context, spaceID int) (models.Space, bool) {
	spaceFromDB, err := s.DB.GetSpace(spaceID)
	if err != nil {
		fmt.Println("Error getting space by id:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return models.Space{}, false
	}
	if spaceFromDB.IsEmpty() {
		fmt.Println("Error space with id not found:", spaceID)
		c.JSON(http.StatusNotFound, models.Error{Error: models.ErrorSpaceWithIdNotFound})
		return models.Space{}, false
	}
	return spaceFromDB, true
}
