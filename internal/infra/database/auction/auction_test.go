package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestWatchAuction(t *testing.T) {
	mongoT := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	auction := auction_entity.Auction{
		Id:          "1",
		ProductName: "Amazing Gadget",
		Category:    "Electronics",
		Description: "A cutting-edge gadget that makes your life easier.",
		Condition:   auction_entity.New,
		Status:      auction_entity.Active,
		Timestamp:   time.Now(),
	}

	durationInterval := time.Duration(time.Second)

	mongoT.Run("test watch auction", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "n", Value: 1},
			{Key: "acknowledged", Value: true},
		})

		repo := NewAuctionRepository(mt.DB)
		repo.auctionInterval = durationInterval
		fmt.Printf("repo.auctionInterval: %v\n", repo.auctionInterval)

		go repo.watchCloseAuction(context.Background(), &auction)

		eventsBeforeInterval := mt.GetAllStartedEvents()
		if len(eventsBeforeInterval) != 0 {
			mt.Error("esperado nenhum evento iniciado antes do intervalo do leilão")
		}

		time.Sleep(durationInterval + 30*time.Millisecond)

		eventsAfterInterval := mt.GetAllStartedEvents()
		if len(eventsAfterInterval) == 0 {
			mt.Error("esperado eventos iniciados após o intervalo do leilão")
		}

		updatesArray, ok := mt.GetStartedEvent().Command.Lookup("updates").ArrayOK()
		if !ok {
			mt.Fatal("esperado que updates seja um array")
		}

		firstUpdateElement, err := updatesArray.IndexErr(0)
		if err != nil {
			mt.Fatalf("esperado que o array tenha pelo menos um elemento: %v", err)
		}

		updateDocument, ok := firstUpdateElement.Value().Document().Lookup("u").Document().Lookup("$set").DocumentOK()
		if !ok {
			mt.Fatal("esperado que $set seja um documento")
		}

		capturedStatus := updateDocument.Lookup("status").Int32()
		if auction_entity.AuctionStatus(capturedStatus) != auction_entity.Completed {
			mt.Errorf("esperado status %v, mas obteve %v", auction_entity.Completed, auction_entity.AuctionStatus(capturedStatus))
		}
	})
}
