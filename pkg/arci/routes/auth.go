package routes

import (
	"log"
	"os"
	"time"

	"arci.it/pkg/arci/db"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type RegisterData struct {
	Email    string `json:"email" binding:"required,email"`
	ShowName string `json:"name" binding:"required"`
}

type LoginData struct {
	Email string `json:"email" binding:"required,email"`
}

type Claims struct {
	MemberID int    `json:"member_id"`
	ShowName string `json:"showname"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	value, exists := os.LookupEnv("JWT_SECRET")
	if !exists {
		log.Fatal("JWT_SECRET environment variable is missing")
	}
	jwtSecret = []byte(value)
}

func generateToken(member *db.Member) (string, error) {
	claims := Claims{
		MemberID: member.ID,
		ShowName: member.ShowName,
		IsAdmin:  member.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			// TODO: Ricordarsi di controllare sull'app gli errori di token scaduto
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func Register(c *gin.Context) {
	var registerData RegisterData
	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	member, err := db.Register(registerData.Email, registerData.ShowName, false)
	if err != nil {
		if err.Error() == "email already registered" || err.Error() == "showname already taken" {
			c.JSON(409, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := generateToken(member)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(201, gin.H{
		"ok":    true,
		"token": token,
		"member": gin.H{
			"id":       member.ID,
			"showname": member.ShowName,
			"is_admin": member.IsAdmin,
		},
	})
}

func Login(c *gin.Context) {
	var loginData LoginData
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	member, err := db.Login(loginData.Email)
	if err != nil {
		if err.Error() == "email not found" {
			c.JSON(401, gin.H{
				"error": "Email inesistente!",
			})
			return
		}
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate JWT token
	token, err := generateToken(member)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "failed to generate token",
		})
		return
	}

	c.JSON(200, gin.H{
		"ok":    true,
		"token": token,
		"member": gin.H{
			"id":       member.ID,
			"showname": member.ShowName,
			"is_admin": member.IsAdmin,
		},
	})
}
