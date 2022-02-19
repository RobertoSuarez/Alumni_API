module github.com/RobertoSuarez/apialumni

// +heroku goVersion go1.16
go 1.16

require (
	github.com/aws/aws-sdk-go v1.43.2
	github.com/gofiber/fiber/v2 v2.24.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.3.0
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/viper v1.10.1
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	gorm.io/driver/postgres v1.2.3
	gorm.io/gorm v1.22.5
)
