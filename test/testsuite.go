package test

import (
	"database/sql"
	"event-booking-api/app/router"
	"event-booking-api/config"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

type ApiTestSuite struct {
	suite.Suite
	dbClient       *sql.DB
	terminateMysql func()
	app            *gin.Engine
	adminToken     string
	user1Token     string
	user2Token     string
}

func (suite *ApiTestSuite) SetupTest() {
	var dsn string
	dsn, suite.dbClient, suite.terminateMysql = setupMysqlContainer()

	os.Setenv("DB_DSN", dsn)
	os.Setenv("JWT_SECRET_KEY", "supersecret")
	os.Setenv("LOG_LEVEL", "DEBUG")

	config.InitLog()
	init := config.Init()
	suite.app = router.Init(init)

	suite.adminToken, _ = generateToken(1, "admin@example.com", 1)
	suite.user1Token, _ = generateToken(2, "user1@example.com", 2)
	suite.user2Token, _ = generateToken(3, "user2@example.com", 2)
}

func (suite *ApiTestSuite) TearDownTest() {
	if suite.dbClient != nil {
		suite.dbClient.Close()
	}

	if suite.terminateMysql != nil {
		suite.terminateMysql()
	}

	os.Unsetenv("DB_DSN")
	os.Unsetenv("JWT_SECRET_KEY")
	os.Unsetenv("LOG_LEVEL")
}

func generateToken(userId int, email string, roleId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userId,
		"email":   email,
		"role_id": roleId,
		"exp":     time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
