package tools

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"turtle/core/lgr"
)

type YieldType[T any] func(func(T) bool)

func Yield[T any](rangeFunc func() []T) YieldType[T] {
	// Readme
	// input func should looks like
	// cursorFunc := func() []int {
	// 	return []int{2,4}
	// }
	return func(yield func(T) bool) {
		defer func() {
			if r := recover(); r != nil {
				lgr.Error("Error in Yield: ", r)
			}
		}()
		for _, val := range rangeFunc() {
			if !yield(val) {
				return
			}
		}
	}
}
func QueryYield[T any](mongo_cursor *mongo.Cursor, ctx context.Context) YieldType[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	return func(yield func(T) bool) {
		defer Recover("Failed to yield")

		for mongo_cursor.Next(ctx) {
			var result T
			if err := mongo_cursor.Decode(&result); err != nil {
				mongo_cursor.Close(ctx)
				return
			}
			if !yield(result) {
				return
			}
		}
	}
}
