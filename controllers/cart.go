package controllers

import (
	"context"
	"e-commerce-fiber/models"
	"e-commerce-fiber/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductStruct struct {
	ID primitive.ObjectID `json:"id"`
}

func AddProductToCart(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)

	defer cancel()

	idLocal := c.Locals("id").(string)

	userId, err := primitive.ObjectIDFromHex(idLocal)

	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid user id")
	}

	var req ProductStruct

	if err := json.Unmarshal(c.Body(), &req); err != nil {
		fmt.Printf("%s+v", &req)
		//omo wetin i dey do
		return utils.ErrorResponse(c, 400, "Error nigga")
	}

	productId := req.ID

	var product models.Product

	err = productCollection.FindOne(ctx, bson.M{"_id": productId}).Decode(&product)
	if err != nil {
		return utils.ErrorResponse(c, 400, "product does not exist")
	}

	fmt.Println(product)

	if product.AvailableQuantity == 0 {
		return utils.ErrorResponse(c, 400, "Product is no longer available")
	}

	var user models.User

	err = userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return utils.ErrorResponse(c, 400, "User does not exist")
	}

	var ProductsToOrder models.ProductsToOrder

	ProductsToOrder.ProductId = productId
	ProductsToOrder.CreatedAt = time.Now()
	ProductsToOrder.UpdatedAt = time.Now()


	if ProductsToOrder.BuyQuantity > product.AvailableQuantity {
		return utils.ErrorResponse(c, 400, "Quantity must be less than product quantity")
	}

	if len(user.UserCart) == 0 {

		_, err := userCollection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$push": bson.M{"userCart": &ProductsToOrder}})
		if err != nil {
			return utils.ErrorResponse(c, 400, "Error adding product to cart")
		}

	} else {
		// check if product is already in cart
		for _, item := range user.UserCart {

			if item.ProductId == productId {
				// update quantity in user cart array inside user document
				if _, err := userCollection.UpdateOne(ctx, bson.M{"_id": userId, "userCart.productId": productId}, bson.M{"$inc": bson.M{"userCart.$.buyQuantity": ProductsToOrder.BuyQuantity}}); err != nil {

					return utils.ErrorResponse(c, 400, "Failed to update product quantity")

				}
			} else {

				_, err = userCollection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$push": bson.M{"userCart": &ProductsToOrder}})
				if err != nil {

					return utils.ErrorResponse(c, 400, "Failed to add product to cart")

				}
			}
		}

	}

	subtractQuantity := product.AvailableQuantity - ProductsToOrder.BuyQuantity

	product.AvailableQuantity = subtractQuantity

	filter := bson.M{"_id": productId}
	update := bson.M{"$set": bson.M{"availableQuantity": subtractQuantity}}
	if _, err := productCollection.UpdateOne(ctx, filter, update); err != nil {

		return utils.ErrorResponse(c, 400, "Failed to update product quantity")

	}

	return utils.SuccessMessage(c, "Product added to cart successfully", ProductsToOrder)

}
