package auction

import (
	"context"
	"os"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// fakeCollection implements the subset of mongo Collection methods used by the repository.
type fakeCollection struct {
	inserted atomic.Value // stores the last inserted document
	updated  atomic.Value // stores a tuple {filter, update}
}

func (f *fakeCollection) InsertOne(ctx context.Context, document interface{}, _ ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	f.inserted.Store(document)
	return &mongo.InsertOneResult{InsertedID: "fake-id"}, nil
}

func (f *fakeCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, _ ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	f.updated.Store([2]interface{}{filter, update})
	return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (f *fakeCollection) FindOne(ctx context.Context, filter interface{}, _ ...*options.FindOneOptions) *mongo.SingleResult {
	// Not needed for this unit test; return nil result.
	return &mongo.SingleResult{}
}

func (f *fakeCollection) Find(ctx context.Context, filter interface{}, _ ...*options.FindOptions) (*mongo.Cursor, error) {
	// Not needed for this unit test; return nil
	return nil, nil
}

func TestCreateAuction_automaticallyClosesAfterInterval(t *testing.T) {
	// set a short auction interval
	os.Setenv("AUCTION_INTERVAL", "200ms")
	defer os.Unsetenv("AUCTION_INTERVAL")

	fake := &fakeCollection{}

	// build a repo with the fake collection
	repo := &AuctionRepository{Collection: fake}

	// create a sample auction entity
	auctionEntity, _ := auction_entity.CreateAuction("product", "category", "description more than ten", auction_entity.New)

	// call CreateAuction which should start a goroutine that updates after the interval
	if err := repo.CreateAuction(context.Background(), auctionEntity); err != nil {
		t.Fatalf("unexpected error calling CreateAuction: %v", err)
	}

	// wait longer than the interval
	time.Sleep(400 * time.Millisecond)

	// check that UpdateOne was called
	val := fake.updated.Load()
	if val == nil {
		t.Fatalf("expected UpdateOne to be called, but it was not")
	}

	tuple := val.([2]interface{})
	filter := tuple[0]
	update := tuple[1]

	// filter should contain the auction id
	mFilter, ok := filter.(bson.M)
	if !ok {
		t.Fatalf("expected filter to be bson.M, got %T", filter)
	}

	if mFilter["_id"] != auctionEntity.Id {
		t.Fatalf("expected filter _id %v, got %v", auctionEntity.Id, mFilter["_id"])
	}

	// update should set status to Completed
	mUpdate, ok := update.(bson.M)
	if !ok {
		// sometimes UpdateOne was called with a nested bson.M ($set)
		if reflect.TypeOf(update).Kind() == reflect.Map {
			// try to inspect via reflect
		}
	}

	// expect update to be bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	top, ok := mUpdate["$set"].(bson.M)
	if !ok {
		t.Fatalf("expected update to be a $set map, got %#v", mUpdate)
	}

	if top["status"] != auction_entity.Completed {
		t.Fatalf("expected status to be Completed (%v), got %v", auction_entity.Completed, top["status"])
	}
}
