package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wgnalvian.com/payment-server/entity"
	"wgnalvian.com/payment-server/exception"
)

type UserService struct {
	Db *mongo.Database
}

func (u *UserService) GetUserByEmail(email string) (*entity.User, error) {
	col := u.Db.Collection("users")

	var user entity.User

	filter := bson.M{"email": email}

	err := col.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		exception.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (u *UserService) Register(user *entity.User) error {
	col := u.Db.Collection("users")

	_, err := col.InsertOne(context.Background(), user)

	if err != nil {
		exception.LogError(err)
		return err
	}

	return nil
}

func (u *UserService) Login(email string) (*entity.User, error) {
	col := u.Db.Collection("users")

	var user entity.User
	filter := bson.M{"email": email}

	err := col.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		exception.LogError(err)
		return nil, err
	}

	return &user, nil

}

func (u *UserService) Transfer(userId string, userOpId string, value int) error {

	col := u.Db.Collection("users")
	colt := u.Db.Collection("transactions")
	// Update balance user

	filter := bson.M{"_id": userId}

	update := bson.M{
		"$inc": bson.M{"balance": -value},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		exception.LogError(err)
		return err
	}

	// Add log transaction

	transaction := entity.Transaction{
		ID:       uuid.New().String(),
		UserId:   userId,
		Amount:   value,
		Desc:     "Transfer to " + userOpId,
		From:     userOpId,
		Type:     0,
		CreateAt: time.Now(),
	}

	_, err = colt.InsertOne(context.Background(), transaction)

	if err != nil {
		exception.LogError(err)
		return err
	}

	// Update balance userOp
	filter = bson.M{"_id": userOpId}

	update = bson.M{
		"$inc": bson.M{"balance": value},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		exception.LogError(err)
		return err
	}

	// Add log transaction

	transaction = entity.Transaction{
		ID:       uuid.New().String(),
		UserId:   userOpId,
		Amount:   value,
		Desc:     "Transfer from " + userId,
		From:     userId,
		Type:     1,
		CreateAt: time.Now(),
	}

	_, err = colt.InsertOne(context.Background(), transaction)

	if err != nil {
		exception.LogError(err)
		return err
	}

	return nil

}

func (u *UserService) GetUserById(userId string) (*entity.User, error) {
	col := u.Db.Collection("users")

	var user entity.User

	filter := bson.M{"_id": userId}

	err := col.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		exception.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (u *UserService) GetTransactions(userId string) ([]entity.Transaction, error) {
	col := u.Db.Collection("transactions")

	filter := bson.M{"user_id": userId}

	cursor, err := col.Find(context.Background(), filter)

	if err != nil {
		exception.LogError(err)
		return nil, err
	}

	var transactions []entity.Transaction

	for cursor.Next(context.Background()) {
		var transaction entity.Transaction
		err := cursor.Decode(&transaction)

		if err != nil {
			exception.LogError(err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (u *UserService) CheckIfEmailExist(email string) bool {
	col := u.Db.Collection("users")

	var user entity.User
	filter := bson.M{"email": email}

	err := col.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return false
		}

		exception.LogError(err)
		return false
	}

	return true
}

func (u *UserService) TopUp(userId string, amount int) error {

	col := u.Db.Collection("users")
	colt := u.Db.Collection("transactions")

	// Update balance user

	filter := bson.M{"_id": userId}

	update := bson.M{
		"$inc": bson.M{"balance": amount},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		exception.LogError(err)
		return err
	}

	// Add log transaction

	transaction := entity.Transaction{
		ID:       uuid.New().String(),
		UserId:   userId,
		Amount:   amount,
		Desc:     "Top Up",
		From:     "",
		Type:     1,
		CreateAt: time.Now(),
	}

	_, err = colt.InsertOne(context.Background(), transaction)

	if err != nil {
		exception.LogError(err)
		return err
	}

	return nil

}
