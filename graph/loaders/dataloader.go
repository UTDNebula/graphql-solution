package loaders

import (
	"context"
	"net/http"
	"time"

	"github.com/KBui4/graphql-workshop/configs"
	"github.com/KBui4/graphql-workshop/graph/model"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type Loaders struct {
	UserLoader    *dataloadgen.Loader[bson.ObjectID, *model.User]
	ProductLoader *dataloadgen.Loader[bson.ObjectID, *model.Product]
	ReviewLoader  *dataloadgen.Loader[bson.ObjectID, *model.Review]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders() *Loaders {
	// Define the dataloaders
	userReader := &userReader{collection: configs.GetCollection("users")}
	productReader := &productReader{collection: configs.GetCollection("products")}
	reviewReader := &reviewReader{collection: configs.GetCollection("reviews")}

	return &Loaders{
		UserLoader:   dataloadgen.NewLoader(userReader.getUsers, dataloadgen.WithWait(time.Millisecond)),
		ProductLoader: dataloadgen.NewLoader(productReader.getProducts, dataloadgen.WithWait(time.Millisecond)),
		ReviewLoader: dataloadgen.NewLoader(reviewReader.getReviews, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(next http.Handler) http.Handler {
	// Return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders()
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}