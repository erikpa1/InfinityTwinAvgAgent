package reflective

import (
	"turtle/core/serverKit"
	"turtle/db"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateDeleteIdRoute(r *gin.Engine, container string, path string, authMethod gin.HandlerFunc) {
	r.DELETE(path, authMethod, func(c *gin.Context) {
		uid := serverKit.MongoObjectIdFromQuery(c)
		db.DeleteEntity(container, bson.M{"_id": uid})
	})
}

func CreateQueryRoute(r *gin.Engine, container string, path string, authMethod gin.HandlerFunc) {
	r.DELETE(path, authMethod, func(c *gin.Context) {
		uid := serverKit.MongoObjectIdFromQuery(c)
		db.DeleteEntity(container, bson.M{"_id": uid})
	})
}
