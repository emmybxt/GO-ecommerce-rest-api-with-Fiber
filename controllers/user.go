package controllers

import (
	"context"
	"e-commerce-fiber/database"
	"e-commerce-fiber/models"
	"e-commerce-fiber/utils"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func HomePage(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"Message": "Welcome to the e-commerce Api service"})
}

func RegisterUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)

	defer cancel()

	var user models.User

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	var userAddress models.Address
	userAddress.ZipCode = gofakeit.Zip()
	userAddress.City = gofakeit.City()
	userAddress.State = gofakeit.State()
	userAddress.Country = gofakeit.Country()
	userAddress.Street = gofakeit.Street()
	userAddress.HouseNumber = gofakeit.StreetNumber()
	user.Address = userAddress
	user.Orders = make([]models.Order, 0)
	user.UserCart = make([]models.ProductsToOrder, 0)

	bodyBytes := c.Body()


	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request payload")
	}

	if user.UserType == "ADMIN" {
		filter := bson.M{"userType": "ADMIN"}

		if _, err := userCollection.FindOne(ctx, filter).DecodeBytes(); err == nil {

			return utils.ErrorResponse(c, 400, "Admin User already exists")

		}
	}

	//check if email exists

	filter := bson.M{"email": user.Email}

	if _, err := userCollection.FindOne(ctx, filter).DecodeBytes(); err == nil {

		return utils.ErrorResponse(c, 400, "Admin email already exists")

	}

	//hash password

	password, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Fatal("Failed hashing Password:", err)
	}
	user.Password = password

	fmt.Print(user)
	if _, err := userCollection.InsertOne(ctx, user); err == nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 400, "Error creating new user account")

	}

	//signtoken

	signedToken, err := utils.CreateJwtToken(user.ID, user.Email, user.UserType)

	fmt.Println(signedToken)

	if err != nil {
		fmt.Println("Error creating signed token")
	}

	// add token to cookie session
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	// set cookie
	c.Cookie(cookie)

	return utils.SuccessMessage(c, "user signed up successfully", user)

}

func Login(c *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required, password"`
	}

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)

	defer cancel()

	var req loginRequest

	bodyBytes := c.Body()

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request payload")
	}

	password := req.Password
	email := req.Email

	fmt.Println(email)
	var existingUser bson.Raw

	if err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser); err != nil {
		fmt.Println(err)

		return utils.ErrorResponse(c, 400, "User does not exist")
	}

	//check password

	isValid := utils.VerifyPassword(password, existingUser.Lookup("password").StringValue())

	if !isValid {

		return utils.ErrorResponse(c, 400, "password is incorrect")
	}

	signJwtToken, err := utils.CreateJwtToken(existingUser.Lookup("_id").ObjectID(), existingUser.Lookup("email").StringValue(), existingUser.Lookup("userType").StringValue())

	if err != nil {
		fmt.Println(err)

		return utils.ErrorResponse(c, 500, "error creating user token")
	}

	// add token to cookie session
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    signJwtToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	// set cookie
	c.Cookie(cookie)

	// data := struct {
	// 	token string
	// 	existingUser string
	// } {
	// 	token: signJwtToken,
	// 	existingUser: existingUser,
	// }

	return utils.SuccessMessage(c, "User signed in successfully", existingUser)

}

