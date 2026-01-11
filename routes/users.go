package routes

import (
	"apiproject/models"
	"apiproject/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signUp(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.Save()
	if err != nil {
		fmt.Println("Error saving user: " + err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user: " + err.Error()})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateCredentials()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password: " + err.Error()})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token: " + err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
