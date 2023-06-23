package controllers

import (
	"context"
	"e-commerce-fiber/models"
	"e-commerce-fiber/utils"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func UpdateAddress (c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	idLocal := c.Locals("id").(string)

	userId, err := primitive.ObjectIDFromHex(idLocal)

	if err != nil {
		return utils.ErrorResponse(c,400, "invalid User Id")
	}

	var address models.Address

	bodyBytes := c.Body()


	if err := json.Unmarshal(bodyBytes, &address); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid request payload")
	}

	if address.City == "" {
		return utils.ErrorResponse(c, 400, "Empty Address field")
	}

	//find user by id & update address 

	filter := bson.M{"_id": userId}
	update := bson.M{"$set": bson.M{"address": address}}

	if _,err := userCollection.UpdateOne(ctx, filter, update); err != nil {
		return utils.ErrorResponse(c, 400,"Error updating address")
	}

	return utils.SuccessMessage(c, "User address updated successfully", address)

}