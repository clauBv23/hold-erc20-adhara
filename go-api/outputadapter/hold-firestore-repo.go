package outputadapter

import (
	"cleanGo/api/usecase/hold"
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type holdRepo struct {
}

func NewHoldFirestoreRepo() hold.HoldRepository {
	return &holdRepo{}
}

const (
	projectId          = "go-api-bac51"
	holdCollectionName = "holds"
)

func (*holdRepo) SaveHold(hold *hold.Hold) (*hold.Hold, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		//log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	_, err = client.Collection(holdCollectionName).Doc(hold.SId).Set(ctx, map[string]interface{}{
		"Id":     hold.Id,
		"SId":    hold.SId,
		"Amount": hold.Amount,
		"User":   hold.User,
		"Status": hold.Status,
	})

	if err != nil {
		//log.Fatalf("Failed adding a new post: %v", err)
		return nil, err
	}

	return hold, nil
}

func (*holdRepo) UpdateHoldStatus(hold *hold.Hold) (*hold.Hold, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		//log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	_, err = client.Collection(holdCollectionName).Doc(hold.SId).Update(ctx, []firestore.Update{
		{
			Path:  "Status",
			Value: "EXEC",
		},
	})

	if err != nil {
		//log.Fatalf("Failed adding a new post: %v", err)
		return nil, err
	}
	return hold, nil
}

func (*holdRepo) FindAllHolds() ([]hold.Hold, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		log.Fatalf("Failed to create a Firestore Client: %v", err)
		return nil, err
	}

	defer client.Close()

	var holds []hold.Hold
	iter := client.Collection(holdCollectionName).Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			log.Fatalf("Failed to iterate the hosts list: %v", err)
			return nil, err
		}

		hold := hold.Hold{
			Id:     doc.Data()["Id"].(int64),
			Amount: doc.Data()["Amount"].(int64),
			User:   doc.Data()["User"].(string),
		}
		holds = append(holds, hold)
	}
	return holds, nil
}

func (r *holdRepo) FindAllHoldsOnCreated() ([]hold.Hold, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	var holds []hold.Hold
	iter := client.Collection(holdCollectionName).Where("Status", "==", "CREATED").Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		hold := hold.Hold{
			Id:     doc.Data()["Id"].(int64),
			SId:    doc.Data()["SId"].(string),
			Amount: doc.Data()["Amount"].(int64),
			User:   doc.Data()["User"].(string),
			Status: doc.Data()["Status"].(string),
		}
		holds = append(holds, hold)
	}
	return holds, nil
}

func (r *holdRepo) FindHoldsFromUser(user string) ([]hold.Hold, error) {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	var holds []hold.Hold
	iter := client.Collection(holdCollectionName).Where("User", "==", user).Documents(ctx)

	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		hold := hold.Hold{
			Id:     doc.Data()["Id"].(int64),
			Amount: doc.Data()["Amount"].(int64),
			User:   doc.Data()["User"].(string),
		}
		holds = append(holds, hold)
	}
	return holds, nil
}
