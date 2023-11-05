package resources

import (
	"net/http"

	"github.com/NimbusX-CMS/NimbusX/api/internal/db"
	"github.com/gin-gonic/gin"
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

func (s *Server) GetWebhooks(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PostWebhooks(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (s *Server) PutWebhooksName(c *gin.Context, name string) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
