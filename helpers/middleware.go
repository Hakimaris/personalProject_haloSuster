package helpers

import (
	"HaloSuster/models"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"HaloSuster/db"

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
	id := fmt.Sprint(parsedToken["Id"])
	nip := fmt.Sprint(parsedToken["Nip"])
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

func AdminAuthMiddleware(c *fiber.Ctx) error {
	token, err := getBearerToken(c.Get("Authorization"))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	id, nip, err := ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// if id.Role != "admin" {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"message": "access not allowed",
	// 	})
	// }

	// find user
	conn := db.CreateConn()
	res, err := conn.Exec("SELECT 1 FROM public.user WHERE nip = $1 LIMIT 1", nip)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "server error",
		})
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	c.Locals("userNip", nip)
	return c.Next()
}
