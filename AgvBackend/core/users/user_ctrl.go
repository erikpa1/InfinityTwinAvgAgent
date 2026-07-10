package users

import (
	"turtle/core/pipeline"
	"turtle/db"
	"turtle/tools"

	"go.mongodb.org/mongo-driver/bson"
)

const CT_USERS = "users"

func COUUser(ctx *pipeline.GinPipeline, user *User) {

	if user.Uid.IsZero() {
		hashPass, err := tools.HashPassword(user.Password)
		if ctx.SetError(err) {
			user.Password = hashPass
			db.InsertEntity(CT_USERS, user)
		}
	} else {
		db.COU(CT_USERS, user.Uid, user)
	}
}

func QueryUsers(user *User, query bson.M) []*User {
	return db.QueryEntities[User](CT_USERS, user.FillOrgQuery(query))
}
