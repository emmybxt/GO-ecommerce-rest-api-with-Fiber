package controllers

import (
	"context"
	"e-commerce-fiber/database"
	"e-commerce-fiber/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var productCollection = database.OpenCollection(database.Client, "products")

type createProduct struct {
	ID                primitive.ObjectID `json:"_id,omitempty"`
	CreatedAt         time.Time          `json:"createdAt,omitempty"`
	UpdatedAt         time.Time          `json:"updatedAt,omitempty"`
	Category          string             `json:"category" validate:"required"`
	Name              string             `json:"name" validate:"required"`
	Price             float64            `json:"price" validate:"required"`
	Description       string             `json:"description" validate:"required"`
	AvailableQuantity int16              `json:"availableQuantity" validate:"required"`
	Images            []string           `json:"images" validate:"required"`
}

func CreateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	userType, err := c.Locals("userType").(string)
	if !err {
		return utils.ErrorResponse(c, 400, "Invalid user data")
	}

	if userType != "ADMIN" {
		return utils.ErrorResponse(c, 400, "You are not allowed to create products")
	}

	var req createProduct

	bodyBytes := c.Body()

	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid payload request")
	}

	


	// id := req.ID
	// createdAt := req.CreatedAt
	// updatedAt := req.UpdatedAt
	// category := req.Category
	name := req.Name
	// price := req.Price
	// description := req.Description
	// availableQuantity := req.AvailableQuantity
	// images := req.Images

	req.ID = primitive.NewObjectID()
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()


	fmt.Println(&req)

	//check if product name already exists
	if err := productCollection.FindOne(ctx, bson.M{"name": name}).Decode(&req); err == nil {
		fmt.Println(err)
		return utils.ErrorResponse(c, 400, "Product already exists")
	}

	//create product

	if _, err := productCollection.InsertOne(ctx, req); err != nil {
		return utils.ErrorResponse(c, 400, "Error creating products")
	}

	return utils.SuccessMessage(c, "product created successfully", req)

}
