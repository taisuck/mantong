package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const stockCollection = "stock"

// Stock 은 시계열 데이터를 수집할 주식(stock 컬렉션) document 이다.
type Stock struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`       // 주식 이름
	Code      string        `bson:"code"`       // 주식 코드 (unique)
	CreatedAt time.Time     `bson:"created_at"` // 생성 시간
}

type StockRepository struct {
	coll *mongo.Collection
}

func NewStockRepository(db *mongo.Database) *StockRepository {
	return &StockRepository{coll: db.Collection(stockCollection)}
}

// FindAll 은 stock 컬렉션의 모든 document 를 반환한다.
func (r *StockRepository) FindAll(ctx context.Context) ([]Stock, error) {
	cur, err := r.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("stock 조회 실패: %w", err)
	}

	var stocks []Stock
	if err := cur.All(ctx, &stocks); err != nil {
		return nil, fmt.Errorf("stock 디코딩 실패: %w", err)
	}
	return stocks, nil
}
