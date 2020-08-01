package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"github.com/sampalm/StarWarsAPI/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var collection *mongo.Collection

func PlanetaCollection(db *mongo.Database) {
	collection = db.Collection("planetas")
}

func getNextSequence() int32 {
	planeta := models.Planeta{}

	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	err := collection.FindOne(context.Background(), bson.M{}, opts).Decode(&planeta)

	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatalf("getNextSequence: Nao foi possivel completar a tarefa. Erro: %s\n", err)
	}

	return atomic.AddInt32(&planeta.ID, 1)
}
func convertQueryToID(s string) int32 {
	var id int32
	if conv, err := strconv.Atoi(s); err == nil {
		id = int32(conv)
	}
	return id
}

func ListPlanetas(c *gin.Context) {
	ctx := context.Background()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Nao foi possivel completar busca. Erro: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Nao foi possivel completar busca",
		})
		return
	}
	defer cursor.Close(ctx)

	planetas := []models.Planeta{}

	for cursor.Next(ctx) {
		var planeta models.Planeta
		if err = cursor.Decode(&planeta); err != nil {
			log.Printf("Erro ao acessar cursor. Erro: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Nao foi possivel completar busca, erro ao acessar cursor",
			})
			return
		}
		planetas = append(planetas, planeta)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Busca completa",
		"data":    planetas,
	})
	return
}

func QueryPlaneta(c *gin.Context) {
	ctx := context.Background()
	p, _ := url.QueryUnescape(c.Param("query"))

	filter := bson.M{
		"$or": []interface{}{
			bson.M{"nome": bson.M{"$regex": p, "$options": "i"}},
			bson.M{"id": convertQueryToID(p)},
		},
	}

	cursor, err := collection.Find(ctx, filter)

	if err != nil {
		log.Printf("Nao foi possivel completar busca. Erro: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Nao foi possivel completar busca",
		})
		return
	}
	defer cursor.Close(ctx)

	planetas := []models.Planeta{}

	for cursor.Next(ctx) {
		var planeta models.Planeta
		if err = cursor.Decode(&planeta); err != nil {
			log.Printf("Erro ao acessar cursor. Erro: %s\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Nao foi possivel completar busca, erro ao acessar cursor",
			})
			return
		}
		planetas = append(planetas, planeta)
	}

	if len(planetas) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Busca completa, mas nenhum planeta encontrado",
			"data":    planetas,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Busca completa.",
		"data":    planetas,
	})
	return
}

func CreatePlaneta(c *gin.Context) {
	var body models.Planeta
	c.BindJSON(&body)

	if len(strings.TrimSpace(body.Clima)) == 0 || len(strings.TrimSpace(body.Terreno)) == 0 {
		log.Printf("Erro dados incompletos\n")
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Erro dados incompletos",
		})
		return
	}

	count, err := SearchPlanetaAPI(body.Nome)

	if err != nil {
		log.Printf("Erro resquest incompleto. Erro: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Erro resquest incompleto ou dados incorretos.",
		})
		return
	}

	body.QtdAparicao = count
	body.ID = getNextSequence()

	_, err = collection.InsertOne(context.Background(), body)

	if err != nil {
		log.Printf("Nao foi possivel inserir novo planeta. Erro: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Nao foi possivel inserir novo planeta",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": fmt.Sprintf("Planeta com ID %d foi criado com sucesso", body.ID)})
	return
}
func RemovePlaneta(c *gin.Context) {
	id := convertQueryToID(c.Param("id"))

	err := collection.FindOneAndDelete(context.Background(), bson.M{"id": id}).Decode(&models.Planeta{})

	if err != nil {
		log.Printf("Nao foi possivel remover planeta com ID %d. Erro: %s\n", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Nao foi possivel remover planeta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": fmt.Sprintf("Planeta ID %d foi removido", id),
	})
	return
}
