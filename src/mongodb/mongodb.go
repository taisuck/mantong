// Package mongodb 는 MongoDB 연결과 컬렉션 조회를 담당한다.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/taisuck/mantong/src/config"
)

// Connect 는 MongoDB 에 연결하고 Ping 으로 연결을 확인한 뒤 대상 Database 를 반환한다.
func Connect(ctx context.Context, cfg config.MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	opts := options.Client().
		ApplyURI(cfg.URI).
		SetConnectTimeout(cfg.Timeout).
		SetServerSelectionTimeout(cfg.Timeout)
	if cfg.User != "" {
		// URI 에 넣지 않고 Credential 로 전달하므로 비밀번호의 특수문자를 escape 할 필요가 없다.
		opts.SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		})
	}

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("mongodb 연결 실패: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, fmt.Errorf("mongodb ping 실패: %w", err)
	}

	return client, client.Database(cfg.Database), nil
}

// Disconnect 는 timeout 내에 MongoDB 연결을 종료한다.
func Disconnect(client *mongo.Client, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return client.Disconnect(ctx)
}
