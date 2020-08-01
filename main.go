package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sampalm/starwars/config"
	c "github.com/sampalm/starwars/controllers"
)

func init() {
	config.Conn()
}

func main() {
	router := NewRouter()
	router.Run()
}

// NewRouter Retorna um router do tipo *gin.Engine e todas as routes disponiveis
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/planetas/", c.ListPlanetas)
	r.GET("/planetas/:query", c.QueryPlaneta)
	r.POST("/planetas/create", c.CreatePlaneta)
	r.DELETE("/planetas/remove/:id", c.RemovePlaneta)
	return r
}
