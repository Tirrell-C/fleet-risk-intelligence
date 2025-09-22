module github.com/Tirrell-C/fleet-risk-intelligence/services/telemetry-ingest

go 1.23

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/joho/godotenv v1.4.0
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.4
	gorm.io/driver/mysql v1.5.2
	gorm.io/gorm v1.25.5
)

replace github.com/Tirrell-C/fleet-risk-intelligence => ../../