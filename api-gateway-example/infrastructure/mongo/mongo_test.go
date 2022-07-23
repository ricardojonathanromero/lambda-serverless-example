package mongo_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/infrastructure/mongo"
	"github.com/ricardojonathanromero/lambda-serverless-example/api-gateway-example/utils"
)

var _ = Describe("unit tests", func() {
	Context("tests", func() {
		_, err := mongo.NewConn(utils.GetEnv("MONGODB_URI", ""))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("connected!")
		}
	})
})
