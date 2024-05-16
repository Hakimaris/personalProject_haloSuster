package helpers

import (
	"HaloSuster/models"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// "HaloSuster/db"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Id  string `json:"id"`
	Nip int64  `json:"nip"`
}

func SignUserJWT(user models.UserModel) string {
	// expiredIn := 28800 // 8 hours
	exp := time.Now().Add(time.Hour * 8)
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(exp),
			Issuer:    "HaloSuster",
		},
		Id:  user.ID,
		Nip: user.NIP,
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	jwtSecret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return ""
	}
	return signedToken
}

func ParseToken(jwtToken string) (string, string, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return "", "", err
	}
	parsedToken, OK := token.Claims.(jwt.MapClaims)
	if !OK {
		return "", "", errors.New("unable to parse claims")
	}
	id := fmt.Sprint(parsedToken["id"])                                   // changed "Id" to "id"
	nip := strconv.FormatFloat(parsedToken["nip"].(float64), 'f', -1, 64) // Convert NIP to string without scientific notation // changed "Nip" to "nip"
	return id, nip, nil
}

func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	// Get the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).SendString("Missing Authorization header")
	}

	// Extract the JWT token from the Authorization header
	tokenStr, err := getBearerToken(authHeader)

	if err != nil {
		return c.Status(401).SendString("Invalid Authorization header format")
	}

	// Parse and validate the JWT token, and extract the Nip
	id, nip, err := ParseToken(tokenStr)
	fmt.Println(tokenStr)
	fmt.Println(id)
	fmt.Println(nip)
	if err != nil {
		return c.Status(401).SendString("Invalid JWT token")
	}

	// Store the Nip in the request context
	c.Locals("userNip", nip)
	c.Locals("userId", id)

	// Continue with the next middleware function or the request handler
	return c.Next()
}
