package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/mattellis91/libretasks-server/config"
	"github.com/mattellis91/libretasks-server/models"
	"github.com/mattellis91/libretasks-server/responses"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")
var validate = validator.New()

func CreateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		c.Header("Content-Type", "application/json")
		
		var user models.User
		defer cancel()

		//request body validation
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		newUserResponse := CreateUser(ctx, user)
		
		c.JSON(newUserResponse.Status, newUserResponse)

		return
		
	}
}

func CreateUser(ctx context.Context, user models.User) responses.UserResponse {
	fmt.Println("%v", user)

	if validationErr := validate.Struct(&user); validationErr != nil {
		return responses.UserResponse{
			Status: http.StatusBadRequest,
			Message: "Error validating user data",
			Data: map[string]interface{}{"data": validationErr.Error()},
		}
	}

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
	}

	result, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {
		return responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		}
	}

	return responses.UserResponse {
		Status:  http.StatusCreated,
		Message: "success",
		Data:    map[string]interface{}{"data": result},
	}

}

func GetAUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.UserResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"data": user}},
		)
	}
}
