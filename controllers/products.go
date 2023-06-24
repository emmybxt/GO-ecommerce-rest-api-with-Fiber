package controllers

import (
	"context"
	"e-commerce-fiber/database"
	"e-commerce-fiber/models"
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

func GetAllProducts(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*8)

	defer cancel()

	cursor, err := productCollection.Find(ctx, bson.M{})

	if err != nil {
		return utils.ErrorResponse(c, 400, "Error fettching products")

	}

	var products []bson.M

	if err := cursor.All(ctx, &products); err != nil {
		return utils.ErrorResponse(c, 400, "error parsing products")
	}

	productsCount := len(products)

	fmt.Println(productsCount)

	return utils.SuccessMessage(c, "products fetched successfully", products)

}

func GetProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

	defer cancel()

	id := c.Params("id")

	productId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid ObjectId")
	}

	//check if product id exists

	var product models.Product

	if err := productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product); err != nil {
		return utils.ErrorResponse(c, 400, "Product id does not exist")
	}

	return utils.SuccessMessage(c, "product fetched successfully", product)
}
