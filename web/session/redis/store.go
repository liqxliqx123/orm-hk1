package redis

import (
	"context"
	"errors"
	"fmt"
	"gitee.com/geektime-geekbang/geektime-go/web/session"
	"github.com/go-redis/redis/v9"
	"time"
)

var errSessionNotExist = errors.New("redis-session: session 不存在")
type StoreOption func(store *Store)

type Store struct {
	prefix string
	client redis.Cmdable
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
// 实际上，这里也可以考虑使用 Option 设计模式，允许用户控制过期检查的间隔
func NewStore(client redis.Cmdable, opts...StoreOption) *Store {
	res := &Store{
		client: client,
		prefix: "session",
		expiration: time.Minute * 15,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (m *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	const lua = `
redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
return redis.call("pexpire", KEYS[1], ARGV[3])
`
	key := m.key(id)
	_, err := m.client.Eval(ctx, lua, []string{key}, "_sess_id", id, m.expiration.Milliseconds()).Result()
	if err != nil {
		return nil, err
	}
	return &memorySession{
		key: key,
		id: id,
		client: m.client,
	}, nil
}

func (s *Store) key(id string) string {
	return fmt.Sprintf("%s_%s", s.prefix, id)
}


func (m *Store) Refresh(ctx context.Context, id string) error {
	key := m.key(id)
	affected, err := m.client.Expire(ctx, key, m.expiration).Result()
	if err != nil {
		return err
	}
	if !affected {
		return errSessionNotExist
	}
	return nil
}

func (m *Store) Remove(ctx context.Context, id string) error {
	_, err := m.client.Del(ctx, m.key(id)).Result()
	return err
}

func (m *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := m.key(id)
	// 这里不需要考虑并发的问题，因为在你检测的当下，没有就是没有
	i, err := m.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if i < 0 {
		return nil, errors.New("redis-session: session 不存在")
	}
	return &memorySession{
		id: id,
		key: key,
		client: m.client,
	}, nil
}

type memorySession struct {
	key string
	id string
	client redis.Cmdable
}

func (m *memorySession) Set(ctx context.Context, key string, val string) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := m.client.Eval(ctx, lua, []string{m.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotExist
	}
	return nil
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	return m.client.HGet(ctx, m.key, key).Result()
}

func (m *memorySession) ID() string {
	return m.id
}

