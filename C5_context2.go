// context package

package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("cwyang", "foo123")
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string) // type assertion of interface type
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"Handling response for %v (auth: %v)\n",
		UserID(ctx),
		AuthToken(ctx),
	)
}
