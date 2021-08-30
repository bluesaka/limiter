/**
令牌桶算法，使用redis做分布式限流

refer to go-zero https://github.com/tal-tech/go-zero
*/

package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	xrate "golang.org/x/time/rate"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
	// KEYS[1] as tokenKey
	// KEYS[2] as timestampKey
	script = `local rate = tonumber(ARGV[1])
local capacity = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])
local fill_time = capacity/rate
local ttl = math.floor(fill_time*2)
local last_tokens = tonumber(redis.call("get", KEYS[1]))
if last_tokens == nil then
    last_tokens = capacity
end

local last_refreshed = tonumber(redis.call("get", KEYS[2]))
if last_refreshed == nil then
    last_refreshed = 0
end

local delta = math.max(0, now-last_refreshed)
local filled_tokens = math.min(capacity, last_tokens+(delta*rate))
local allowed = filled_tokens >= requested
local new_tokens = filled_tokens
if allowed then
    new_tokens = filled_tokens - requested
end

redis.call("setex", KEYS[1], ttl, new_tokens)
redis.call("setex", KEYS[2], ttl, now)

return allowed`

	tokenFormat     = "{%s}.tokens"
	timestampFormat = "{%s}.ts"
	pingInterval    = time.Millisecond * 100
)

const (
	defaultRate             = 5
	defaultCapacity         = 10
	defaultRedisAddr        = "localhost:6379"
	defaultRedisPassword    = ""
	defaultRedisDB          = 0
	defaultRedisDialTimeout = 5 * time.Second
	defaultRedisPoolSize    = 3
)

var defaultTokenLimit = TokenLimit{
	rate:     defaultRate,
	capacity: defaultCapacity,
}

// TokenLimit controls how frequently events are allowed to happen with in one second.
type TokenLimit struct {
	rate           int // 每秒放入令牌数
	capacity       int // 令牌桶容量
	redis          *redis.Client
	tokenKey       string
	timestampKey   string
	redisAlive     uint32
	mu             sync.Mutex
	rateLimiter    *xrate.Limiter
	monitorStarted bool
}

type Option func(*TokenLimit)

// NewTokenLimiter returns a PeriodLimit object.
// opts can be used to customize the PeriodLimit.
func NewTokenLimiter(key string, opts ...Option) *TokenLimit {
	p := &defaultTokenLimit
	for _, o := range opts {
		o(p)
	}

	if p.redis == nil {
		client := redis.NewClient(&redis.Options{
			Addr:        defaultRedisAddr,
			Password:    defaultRedisPassword,
			DB:          defaultRedisDB,
			DialTimeout: defaultRedisDialTimeout,
			PoolSize:    defaultRedisPoolSize,
		})
		if err := client.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
		p.redis = client
	}

	p.tokenKey = fmt.Sprintf(tokenFormat, key)
	p.timestampKey = fmt.Sprintf(timestampFormat, key)
	p.redisAlive = 1
	p.rateLimiter = xrate.NewLimiter(xrate.Every(time.Second/time.Duration(p.rate)), p.capacity)
	return p
}

// AllowN reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate.
// Otherwise use Reserve or Wait.
func (tl *TokenLimit) AllowN(now time.Time, n int) bool {
	if atomic.LoadUint32(&tl.redisAlive) == 0 {
		return tl.rateLimiter.AllowN(now, n)
	}
	resp, err := tl.redis.Eval(
		context.Background(),
		script,
		[]string{tl.tokenKey, tl.timestampKey},
		[]string{
			strconv.Itoa(tl.rate),
			strconv.Itoa(tl.capacity),
			strconv.FormatInt(now.Unix(), 10),
			strconv.Itoa(n),
		},
	).Result()

	if err == redis.Nil {
		return false
	} else if err != nil {
		tl.startMonitor()
		return tl.rateLimiter.AllowN(now, n)
	}

	code, ok := resp.(int64)
	if !ok {
		tl.startMonitor()
		return tl.rateLimiter.AllowN(now, n)
	}

	// redis allowed == true
	// Lua boolean true -> r integer reply with value of 1
	return code == 1
}

func (tl *TokenLimit) startMonitor() {
	tl.mu.Lock()
	tl.mu.Unlock()

	if tl.monitorStarted {
		return
	}

	tl.monitorStarted = true
	atomic.StoreUint32(&tl.redisAlive, 0)
	go tl.waitForRedis()
}

func (tl *TokenLimit) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		tl.mu.Lock()
		tl.monitorStarted = false
		tl.mu.Unlock()
	}()

	for range ticker.C {
		if err := tl.redis.Ping(context.Background()).Err(); err == nil {
			return
		}
	}
}

// WithRate 设置令牌桶每秒放入token速率
func WithRate(rate int) Option {
	return func(tl *TokenLimit) {
		tl.rate = rate
	}
}

// WithCapacity 设置令牌桶容量
func WithCapacity(capacity int) Option {
	return func(tl *TokenLimit) {
		tl.capacity = capacity
	}
}
