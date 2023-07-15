package controllers

import (
	"context"
	"net/http"
	"time"
	"fmt"

	"github.com/mattellis91/libretasks-server/config"
	"github.com/mattellis91/libretasks-server/models"
	"github.com/mattellis91/libretasks-server/responses"
	"github.com/mattellis91/libretasks-server/util"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

var userCollection *mongo.Collection = config.GetCollection(config.DB, "user")
var validate = validator.New()

func SignUpHandler() gin.HandlerFunc {
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
		
	}
}

func CreateUser(ctx context.Context, user models.User) responses.UserResponse {
	
	if validationErr := validate.Struct(&user); validationErr != nil {
		return responses.UserResponse{
			Status: http.StatusBadRequest,
			Message: "Error validating user data",
			Data: map[string]interface{}{"data": validationErr.Error()},
		}
	}

	fmt.Println(user.Password)

	newUser := models.User{
		ID:       primitive.NewObjectID(),
		Username: user.Username,
		Email:    user.Email,
		Password: util.GetHashedValue([]byte(user.Password)),
	}

	_, err := userCollection.InsertOne(ctx, newUser)

	if err != nil {

		if mongo.IsDuplicateKeyError(err) {
			return responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": "User with that email already exists."},
			}
		}

		return responses.UserResponse{
			Status:  http.StatusInternalServerError,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		}
	}

	newUser.Password = ""

	return responses.UserResponse {
		Status:  http.StatusCreated,
		Message: "success",
		Data:    map[string]interface{}{"data": newUser},
	}

}


func LoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		type LoginData struct {
			Email string `json:"email`
			Password string `json:"password`
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		c.Header("Content-Type", "application/json")
		defer cancel()

		var loginData LoginData
		if err := c.BindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": err.Error()}},
			)
			return
		}

		fmt.Println("%s", loginData.Email)
		fmt.Println("%s", loginData.Password)

		var foundUser models.User 
		err := userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&foundUser)

		if err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": "Failed to login"}},
			)
			return;
		}

		fmt.Println(loginData.Password)
		fmt.Println(foundUser.Password)

		fmt.Println(util.HashedValueMatches("$2a$04$3tJg.QD3wCquGB0OEaOrKO9r4TROB3EkV5vPXqThEYCZos4vcvTyy", []byte("password")))

		if !util.HashedValueMatches(foundUser.Password, []byte(loginData.Password)) {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": "Auth failed"}},
			)
			return;
		}

		foundUser.Password = ""

		c.JSON(http.StatusOK, responses.UserResponse{
			Status: http.StatusOK, Message: "",
			Data: map[string]interface{}{"data": foundUser}},
		)

	}
}


