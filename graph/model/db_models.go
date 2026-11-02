package model

import "go.mongodb.org/mongo-driver/v2/bson"

// Database models for MongoDB

// GraphQL generates the ID fields as strings (see models_gen.go)
// For any objects that reference another with ID, when retrieving data
// we need to convert from the database model to the GraphQL model.

type DBUser struct {
	ID            bson.ObjectID            `bson:"_id"`
	Name          string                   `bson:"name"`
	Hobbies       []string                 `bson:"hobbies"`
	Products      []bson.ObjectID          `bson:"products"`
	Reviews       []bson.ObjectID          `bson:"reviews"`
}

type DBProduct struct {
	ID            bson.ObjectID            `bson:"_id"`
	Seller        bson.ObjectID            `bson:"seller"`
	Price         float64                  `bson:"price"`
	Name          string                   `bson:"name"`
	Description   string                   `bson:"description"`
	Reviews       []bson.ObjectID          `bson:"reviews"`
}

// Challenge 2 Part A
type DBReview struct {
	ID            bson.ObjectID            `bson:"_id"`
	Rating        int32                      `bson:"rating"`
	Title         string                   `bson:"title"`
	Description   string                   `bson:"description"`
	Author        bson.ObjectID            `bson:"author"`
	Product       bson.ObjectID            `bson:"product"`
}


// Transform functions
// These functions convert from the database model to the GraphQL model.

func TransformUser(dbUser *DBUser) *User {
	if dbUser == nil {
		return nil
	}

	products := make([]*Product, len(dbUser.Products))
	for i, productID := range dbUser.Products {
		products[i] = &Product{ID: productID.Hex()}
	}

	reviews := make([]*Review, len(dbUser.Reviews))
	for i, reviewID := range dbUser.Reviews {
		reviews[i] = &Review{ID: reviewID.Hex()}
	}

	return &User{
		ID:       dbUser.ID.Hex(),
		Name:     dbUser.Name,
		Hobbies:  dbUser.Hobbies,
		Products: products,
		Reviews:  reviews,
	}
}

func TransformProduct(dbProduct *DBProduct) *Product {
	if dbProduct == nil {
		return nil
	}

	userRef := &User{
		ID: dbProduct.Seller.Hex(),
	}

	reviews := make([]*Review, len(dbProduct.Reviews))
	for i, reviewID := range dbProduct.Reviews {
		reviews[i] = &Review{ID: reviewID.Hex()}
	}

	return &Product{
		ID:          dbProduct.ID.Hex(),
		Seller:      userRef,
		Price:       dbProduct.Price,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		Reviews:     reviews,
	}
}

// Challenge 2 Part B
func TransformReview(dbReview *DBReview) *Review {
	if dbReview == nil {
		return nil
	}

	authorRef := &User{
		ID: dbReview.Author.Hex(),
	}

	productRef := &Product{
		ID: dbReview.Product.Hex(),
	}

	return &Review{
		ID:          dbReview.ID.Hex(),
		Title:       dbReview.Title,
		Description: dbReview.Description,
		Author:      authorRef,
		Product:     productRef,
	}
}
