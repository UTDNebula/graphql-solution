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
	var dbUsers []*model.DBUser
	if err := cursor.All(ctx, &dbUsers); err != nil {
		return nil, []error{err}
	}

	// $in does not return in order of given userIDs
	// so we order with map.
	userMap := make(map[bson.ObjectID]*model.User, len(dbUsers))
	for _, dbUser := range dbUsers {
		userMap[dbUser.ID] = model.TransformUser(dbUser)
	}

	users := make([]*model.User, len(userIDs))
	errs := make([]error, len(userIDs))
	for i, ID := range userIDs {
		users[i] = userMap[ID]
		if users[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return users, errs
}

// getProducts implements a batch function that can retrieve many products by ID for use in a dataloader
func (p *productReader) getProducts(ctx context.Context, productIDs []bson.ObjectID) ([]*model.Product, []error) {
	cursor, err := p.collection.Find(ctx, bson.M{"_id": bson.M{"$in": productIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var dbProducts []*model.DBProduct
	if err := cursor.All(ctx, &dbProducts); err != nil {
		return nil, []error{err}
	}

	productMap := make(map[bson.ObjectID]*model.Product, len(dbProducts))
	for _, dbProduct := range dbProducts {
		productMap[dbProduct.ID] = model.TransformProduct(dbProduct)
	}

	products := make([]*model.Product, len(productIDs))
	errs := make([]error, len(productIDs))
	for i, ID := range productIDs {
		products[i] = productMap[ID]
		if products[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return products, errs
}

// getReviews implements a batch function that can retrieve many reviews by ID for use in a dataloader
func (r *reviewReader) getReviews(ctx context.Context, reviewIDs []bson.ObjectID) ([]*model.Review, []error) {
	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": reviewIDs}})
	if err != nil {
		return nil, []error{err}
	}
	defer cursor.Close(ctx)

	var dbReviews []*model.DBReview
	if err := cursor.All(ctx, &dbReviews); err != nil {
		return nil, []error{err}
	}

	reviewMap := make(map[bson.ObjectID]*model.Review, len(dbReviews))
	for _, dbReview := range dbReviews {
		reviewMap[dbReview.ID] = model.TransformReview(dbReview)
	}

	reviews := make([]*model.Review, len(reviewIDs))
	errs := make([]error, len(reviewIDs))
	for i, ID := range reviewIDs {
		reviews[i] = reviewMap[ID]
		if reviews[i] == nil {
			errs[i] = mongo.ErrNoDocuments
		}
	}

	return reviews, errs
}

// For convenience in resolvers

func GetUser(ctx context.Context, userRef *model.User) (*model.User, error) {
	userID, err := bson.ObjectIDFromHex(userRef.ID)
	if err != nil {
		return nil, err
	}

	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

func GetUsers(ctx context.Context, usersRef []*model.User) ([]*model.User, error) {
	userIDs := make([]bson.ObjectID, len(usersRef))
	for i, user := range usersRef {
		userID, err := bson.ObjectIDFromHex(user.ID)
		if err != nil {
			return nil, err
		}
		userIDs[i] = userID
	}

	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}

func GetProduct(ctx context.Context, productRef *model.Product) (*model.Product, error) {
	productID, err := bson.ObjectIDFromHex(productRef.ID)
	if err != nil {
		return nil, err
	}

	loaders := For(ctx)
	return loaders.ProductLoader.Load(ctx, productID)
}

func GetProducts(ctx context.Context, productsRef []*model.Product) ([]*model.Product, error) {
	productIDs := make([]bson.ObjectID, len(productsRef))
	for i, product := range productsRef {
		productID, err := bson.ObjectIDFromHex(product.ID)
		if err != nil {
			return nil, err
		}
		productIDs[i] = productID
	}

	loaders := For(ctx)
	return loaders.ProductLoader.LoadAll(ctx, productIDs)
}

func GetReview(ctx context.Context, reviewRef *model.Review) (*model.Review, error) {
	reviewID, err := bson.ObjectIDFromHex(reviewRef.ID)
	if err != nil {
		return nil, err
	}

	loaders := For(ctx)
	return loaders.ReviewLoader.Load(ctx, reviewID)
}

func GetReviews(ctx context.Context, reviewsRef []*model.Review) ([]*model.Review, error) {
	reviewIDs := make([]bson.ObjectID, len(reviewsRef))
	for i, review := range reviewsRef {
		reviewID, err := bson.ObjectIDFromHex(review.ID)
		if err != nil {
			return nil, err
		}
		reviewIDs[i] = reviewID
	}

	loaders := For(ctx)
	return loaders.ReviewLoader.LoadAll(ctx, reviewIDs)
}
