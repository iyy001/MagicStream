package utils

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/database"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var UserCollection *mongo.Collection = database.OpenCollection("users")
var SECRET_KEY string = os.Getenv("SECRET_KEY")
var REFRESH_SECRET_KEY string = os.Getenv("REFRESH_TOKEN_SECRET_KEY")

func GenerateAllTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "MagicStream",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			//ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Second)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(REFRESH_SECRET_KEY))
	fmt.Println("refreshToken", signedRefreshToken)

	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}
func GetAccessToken(c *gin.Context) (string, error) {
	// authHeader := c.Request.Header.Get("Authorization")
	// if authHeader == "" {
	// 	return "", errors.New("Authorization header is required")
	// }
	// tokenString := authHeader[len("Bearer "):]
	// if tokenString == "" {
	// 	return "", errors.New("Bearer token is required")
	// }
	tokenString, err := c.Cookie("access_token")
	if err != nil {

		return "", err
	}

	return tokenString, nil
}
func ValidateToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
func ValidateRefreshToken(tokenString string) (*SignedDetails, error) {
	claims := &SignedDetails{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {

		return []byte(REFRESH_SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		fmt.Println("refresh token has expired")
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}
func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")

	if !exists {
		return "", errors.New("role does not exist in this context")
	}

	id, ok := userId.(string)
	if !ok {
		return "", errors.New("unable to retrieve role")
	}

	return id, nil

}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")

	if !exists {
		return "", errors.New("role does not exist in this context")
	}

	memberRole, ok := role.(string)
	if !ok {
		return "", errors.New("unable to retrieve role")
	}

	return memberRole, nil

}
func UpdateAllTokens(userId string, token string, refreshToken string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedData := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updatedAt,
		},
	}

	//	opt := &options.UpdateOptions{Upsert: new(bool)}
	//	*opt.Upsert = true

	_, err = UserCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updatedData)
	if err != nil {
		return err
	}
	return nil
}
