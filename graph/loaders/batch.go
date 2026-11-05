package loaders

import (
	"context"

	"github.com/KBui4/graphql-workshop/graph/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// File for all the batch functions

type userReader struct {
	collection *mongo.Collection
}

type productReader struct {
	collection *mongo.Collection
}

type reviewReader struct {
	collection *mongo.Collection
}

// getUsers implements a batch function that can retrieve many users by ID for use in a dataloader
func (u *userReader) getUsers(ctx context.Context, userIDs []bson.ObjectID) ([]*model.User, []error) {
	// Find documents matching userIDs
	cursor, err := u.collection.Find(ctx, bson.M{"_id": bson.M{"$in": userIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	// Decode
	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given userIDs
	// so we order with map.
	userMap := make(map[bson.ObjectID]*model.User, len(users))
	for _, user := range users {
		userMap[user.ID] = user
	}

	orderedUsers := make([]*model.User, len(userIDs))
	errs := make([]error, len(userIDs))
	for i, ID := range userIDs {
		orderedUsers[i] = userMap[ID]
		if orderedUsers[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return orderedUsers, errs
}

// getProducts implements a batch function that can retrieve many products by ID for use in a dataloader
func (p *productReader) getProducts(ctx context.Context, productIDs []bson.ObjectID) ([]*model.Product, []error) {
	cursor, err := p.collection.Find(ctx, bson.M{"_id": bson.M{"$in": productIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var products []*model.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, []error{err}
	}

	productMap := make(map[bson.ObjectID]*model.Product, len(products))
	for _, product := range products {
		productMap[product.ID] = product
	}

	orderedProducts := make([]*model.Product, len(productIDs))
	errs := make([]error, len(productIDs))
	for i, ID := range productIDs {
		orderedProducts[i] = productMap[ID]
		if orderedProducts[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return orderedProducts, errs
}

// getReviews implements a batch function that can retrieve many reviews by ID for use in a dataloader
func (r *reviewReader) getReviews(ctx context.Context, reviewIDs []bson.ObjectID) ([]*model.Review, []error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": reviewIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var reviews []*model.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, []error{err}
	}

	reviewMap := make(map[bson.ObjectID]*model.Review, len(reviews))
	for _, review := range reviews {
		reviewMap[review.ID] = review
	}

	orderedReviews := make([]*model.Review, len(reviewIDs))
	errs := make([]error, len(reviewIDs))
	for i, ID := range reviewIDs {
		orderedReviews[i] = reviewMap[ID]
		if orderedReviews[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return orderedReviews, errs
}

// For convenience in resolvers

func GetUser(ctx context.Context, userRef bson.ObjectID) (*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userRef)
}

func GetUsers(ctx context.Context, usersRef []bson.ObjectID) ([]*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, usersRef)
}

func GetProduct(ctx context.Context, productRef bson.ObjectID) (*model.Product, error) {
	loaders := For(ctx)
	return loaders.ProductLoader.Load(ctx, productRef)
}

func GetProducts(ctx context.Context, productsRef []bson.ObjectID) ([]*model.Product, error) {
	loaders := For(ctx)
	return loaders.ProductLoader.LoadAll(ctx, productsRef)
}

func GetReview(ctx context.Context, reviewRef bson.ObjectID) (*model.Review, error) {
	loaders := For(ctx)
	return loaders.ReviewLoader.Load(ctx, reviewRef)
}

func GetReviews(ctx context.Context, reviewsRef []bson.ObjectID) ([]*model.Review, error) {
	loaders := For(ctx)
	return loaders.ReviewLoader.LoadAll(ctx, reviewsRef)
}
