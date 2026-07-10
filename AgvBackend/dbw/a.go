package dbw

import (
	"turtle/core/lgr"
	"turtle/core/users"

	"go.mongodb.org/mongo-driver/bson"
)

func QueryOrgEntities(user *users.User) {
	query := bson.M{
		"org": user.Org,
	}

	lgr.Error(query)
}

func QueryUserData(user *users.User) {
	query := bson.M{
		"org":  user.Org,
		"user": user.Uid,
	}

	lgr.Error(query)

}
