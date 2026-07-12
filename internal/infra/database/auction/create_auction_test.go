package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestCreateAuction_ShouldCloseAuctionAutomatically(t *testing.T) {
	ctx := context.Background()

	os.Setenv("AUCTION_DURATION", "200ms")
	defer os.Unsetenv("AUCTION_DURATION")

	mongoContainer, err := mongodb.Run(ctx, "mongo:7")
	if err != nil {
		t.Fatalf("error starting mongodb container: %v", err)
	}

	defer func() {
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Fatalf("error terminating mongodb container: %v", err)
		}
	}()

	connectionString, err := mongoContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("error getting mongodb connection string: %v", err)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		t.Fatalf("error connecting to mongodb: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("error disconnecting from mongodb: %v", err)
		}
	}()

	database := client.Database("auction_test")
	defer database.Drop(ctx)

	repository := NewAuctionRepository(database)

	auctionEntity, internalErr := auction_entity.CreateAuction(
		"Product test",
		"Category test",
		"Description test with enough characters",
		auction_entity.New,
	)
	if internalErr != nil {
		t.Fatalf("error creating auction entity: %v", internalErr)
	}

	internalErr = repository.CreateAuction(ctx, auctionEntity)
	if internalErr != nil {
		t.Fatalf("error creating auction: %v", internalErr)
	}

	createdAuction, internalErr := repository.FindAuctionById(ctx, auctionEntity.Id)
	if internalErr != nil {
		t.Fatalf("error finding created auction: %v", internalErr)
	}

	if createdAuction.Status != auction_entity.Active {
		t.Fatalf("expected auction status Active, got %v", createdAuction.Status)
	}

	time.Sleep(300 * time.Millisecond)

	closedAuction, internalErr := repository.FindAuctionById(ctx, auctionEntity.Id)
	if internalErr != nil {
		t.Fatalf("error finding closed auction: %v", internalErr)
	}

	if closedAuction.Status != auction_entity.Completed {
		t.Fatalf("expected auction status Completed, got %v", closedAuction.Status)
	}
}
