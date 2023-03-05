package repositories

import (
	"acourse_tag_cart_bookmark_service/pkg/contracts"
	"acourse_tag_cart_bookmark_service/pkg/models"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type CartDatabaseRepository struct {
	Connection *mongo.Database
	Collection *mongo.Collection
}

func (c CartDatabaseRepository) Create(ctx context.Context, cart *models.Cart) (cartId primitive.ObjectID, err error) {

	var courseId primitive.ObjectID

	//	Use Transaction
	err = c.Connection.Client().UseSession(ctx, func(sessionContext mongo.SessionContext) error {

		// Start Transaction
		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		// Insert Data To the Database & abort if it fails
		insertedData, err := c.Collection.InsertOne(ctx, cart)
		if err != nil {
			sessionContext.AbortTransaction(ctx)
			return err
		}

		courseId = insertedData.InsertedID.(primitive.ObjectID)

		// Commit Data if no error
		err = sessionContext.CommitTransaction(ctx)
		if err != nil {
			return err
		}

		return nil

	})

	if err != nil {
		return primitive.NilObjectID, err
	}

	return courseId, nil
}

func (c CartDatabaseRepository) FetchById(ctx context.Context, id string, exclude []string) (cart models.Cart, err error) {

	//1. Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	//2. Set options
	opts := options.FindOne().SetProjection(excluded)

	//3. Collection result
	var cartRes models.Cart
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return cartRes, err
	}

	//4. Setup Filter
	filter := map[string]interface{}{"_id": objectID, "deleted_at": nil}
	err = c.Collection.FindOne(ctx, filter, opts).Decode(&cartRes)

	if err != nil {
		return cartRes, err
	}

	return cartRes, nil
}

func (c CartDatabaseRepository) FetchByUserId(ctx context.Context, userID string, exclude []string) (cart models.Cart, err error) {

	//Exclude fields
	excluded := make(map[string]int)
	for _, field := range exclude {
		excluded[field] = 0
	}

	opts := options.FindOne().SetProjection(excluded)

	var bookmark models.Cart

	filter := map[string]interface{}{"user_id": userID, "deleted_at": nil}

	err = c.Collection.FindOne(ctx, filter, opts).Decode(&bookmark)
	if err != nil {
		return bookmark, err
	}

	return bookmark, nil

}

func (c CartDatabaseRepository) AddCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error) {

	//1. Filter by id
	filter := bson.D{{"user_id", userID}}

	//2. Convert id string to ObjectID
	coursesObjID := make([]bson.D, 0)
	for _, c := range coursesID {
		cID, err := primitive.ObjectIDFromHex(c)
		if err == nil {
			coursesObjID = append(coursesObjID, bson.D{{"id", cID}})
		}
		//coursesObjID = append(coursesObjID, bson.D{{"id", c.ID}})
	}

	//3. Prepare statement
	statement := bson.M{"$addToSet": bson.M{"courses": bson.M{"$each": coursesObjID}}}

	//4. Update data
	result, err := c.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY ADD COURSE: ", err.Error())
		return false, err
	}

	//5. Check if document exist / matched by the filter statements
	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY ADD COURSE: document not matched")
		return false, mongo.ErrNoDocuments
	}

	////6.
	//if result.ModifiedCount == 0 {
	//	log.Println("BOOKMARK REPOSITORY ADD COURSE: document not matched")
	//	return false, errors.New("document not modified")
	//}

	return true, nil

}

func (c CartDatabaseRepository) RevokeCourse(ctx context.Context, userID string, coursesID []string) (status bool, err error) {

	filter := bson.D{{"user_id", userID}}

	cID := make([]primitive.ObjectID, 0)
	for _, s := range coursesID {
		cID = append(cID, models.GenerateObjectIDFromHex(s))
	}

	statement := bson.M{"$pull": bson.M{"courses": bson.M{"id": bson.M{"$in": cID}}}}

	result, err := c.Collection.UpdateOne(ctx, filter, statement)
	if err != nil {
		log.Println("BOOKMARK REPOSITORY DELETE COURSE: ", err.Error())
		return false, err
	}

	if result.MatchedCount == 0 {
		log.Println("BOOKMARK REPOSITORY DELETE COURSE: document not matched")
		return false, mongo.ErrNoDocuments
	}

	if result.ModifiedCount == 0 {
		log.Println("BOOKMARK REPOSITORY DELETE COURSE: document not modified")
		//return false, errors.New("document not modified")
	}

	return true, nil
}

func (c CartDatabaseRepository) Delete(ctx context.Context, cartID string) (status bool, err error) {

	objectID, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return false, err
	}

	result, err := c.Collection.DeleteOne(ctx, bson.D{{"_id", objectID}})
	if err != nil {
		return false, err
	}

	if result.DeletedCount == 0 {
		return false, errors.New("mongo: no documents in result")
	}

	return true, nil
}

func ConstructCartDBRepository(conn *mongo.Database, coll *mongo.Collection) contracts.CartDBRepository {
	return &CartDatabaseRepository{
		Connection: conn,
		Collection: coll,
	}
}
