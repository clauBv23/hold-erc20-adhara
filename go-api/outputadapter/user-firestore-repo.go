package outputadapter

import (
	"cleanGo/api/usecase/user"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type UserRepo struct {
}

func NewUserFirestoreRepo() user.UserRepository {
	return &UserRepo{}
}

const (
	userCollectionName = "users"
)

func (*UserRepo) SaveUser(user *user.User) (*user.User, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	_, _, err = client.Collection(userCollectionName).Add(ctx, map[string]interface{}{
		"Id":      user.Id,
		"Address": user.Address,
	})

	if err != nil {
		log.Fatalf("Failed adding a new post: %v", err)
		return nil, err
	}
	return user, nil
}

func (*UserRepo) FindAllUsers() ([]user.User, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	var users []user.User
	iter := client.Collection(userCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("Failed to iterate the users list: %v", err)
			return nil, err
		}

		user := user.User{
			Id:      doc.Data()["Id"].(int64),
			Address: doc.Data()["Address"].(string),
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepo) FindUserByAddress(addr string) (*user.User, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}
	defer client.Close()

	iter := client.Collection(userCollectionName).Where("Address", "==", addr).Limit(1).Documents(ctx)
	var foundUser user.User

	// get the first element
	doc, err := iter.Next()
	if err == iterator.Done {
		foundUser = user.User{
			Id:      0,
			Address: "",
		}
		return &foundUser, nil
	}
	if err != nil {
		log.Fatalf("Failed get the user: %v", err)
		return nil, err
	}

	foundUser = user.User{
		Id:      doc.Data()["Id"].(int64),
		Address: doc.Data()["Address"].(string),
	}
	return &foundUser, nil
}
