package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenClaims struct {
	Id    primitive.ObjectID `json:"id"`
	Email string             `json:"email"`
	jwt.StandardClaims
	UserType string `json:"userType`
}

//create a jwt token

func createJwtToken(id primitive.ObjectID, email string, userType string) (tokenString string, err error) {
	tokenClaims := &TokenClaims{
		Id:       id,
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	if signToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", err
	} else {
		return signToken, nil
	}

}
