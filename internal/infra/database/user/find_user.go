package user

import (
	"context"
	"errors"
	"fmt"
	"fullcycle-auction_go/internal/entity/user_entity"
	"fullcycle-auction_go/internal/internal_error"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (ur *UserRepository) FindUserById(
	ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {

	// Converte o userId para ObjectID
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		log.Printf("Erro na conversão do userId para ObjectID: %v", err)
		return nil, internal_error.NewBadRequestError("Invalid userId format")
	}

	// Filtro de busca pelo ObjectID
	filter := bson.M{"_id": objectId}

	var userEntityMongo UserEntityMongo

	err = ur.Collection.FindOne(ctx, filter).Decode(&userEntityMongo)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("User not found with this id = %s", userId)
			return nil, internal_error.NewNotFoundError(
				fmt.Sprintf("User not found with this id = %s", userId))
		}

		log.Printf("Error trying to find user by userId: %v", err)
		return nil, internal_error.NewInternalServerError("Error trying to find user by userId")
	}

	// Mapeia o resultado do MongoDB para a entidade de usuário do sistema
	userEntity := &user_entity.User{
		Id:   userEntityMongo.Id.Hex(),
		Name: userEntityMongo.Name,
	}

	return userEntity, nil
}
