/**
固定窗口计数器算法，使用redis做分布式限流

refer to go-zero https://github.com/tal-tech/go-zero
*/

package redis

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

const (
	// in order to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
	periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("pexpire", KEYS[1], window)
    return 1
elseif current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 3
end
`
)

const (
	defaultPeriod           = 2000 // 时间窗口2000ms
	defaultQuota            = 10   // 最大请求数
	defaultRedisAddr        = "localhost:6379"
	defaultRedisPassword    = ""
	defaultRedisDB          = 0
	defaultRedisDialTimeout = 5 * time.Second
	defaultRedisPoolSize    = 3
)

const (
	LimitUnknown = iota
	LimitAllowed
	LimitHitQuota
	LimitOverQuota
)

var (
	ErrUnknownRedisCode = errors.New("unknown redis status code")

	defaultLimit = PeriodLimit{
		period: defaultPeriod,
		quota:  defaultQuota,
	}
)

// PeriodLimit
type PeriodLimit struct {
	key    string
	period int // 统计周期 ms
	quota  int // 最大请求数
	redis  *redis.Client
}

// NewPeriodLimit returns a PeriodLimit object.
// opts can be used to customize the PeriodLimit.
func NewPeriodLimit(key string, opts ...Option) *PeriodLimit {
	p := &defaultLimit
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

	p.key = key
	return p
}

// Exec exec a request and returns the state
func (p *PeriodLimit) Exec() (int, error) {
	resp, err := p.redis.Eval(context.Background(), periodScript, []string{p.key}, []string{
		strconv.Itoa(p.quota),
		strconv.Itoa(p.period),
	}).Result()
	if err != nil {
		return LimitUnknown, err
	}

	code, ok := resp.(int64)
	if !ok {
		return LimitUnknown, ErrUnknownRedisCode
	}

	switch code {
	case LimitAllowed:
		return LimitAllowed, nil
	case LimitHitQuota:
		return LimitHitQuota, nil
	case LimitOverQuota:
		return LimitOverQuota, nil
	default:
		return 0, ErrUnknownRedisCode
	}
}

type Option func(*PeriodLimit)

// WithInterval returns a function to set the interval of a PeriodLimit
func WithInterval(period int) Option {
	return func(p *PeriodLimit) {
		p.period = period
	}
}

// WithLimit returns a function to set the quota of a PeriodLimit
func WithLimit(quota int) Option {
	return func(p *PeriodLimit) {
		p.quota = quota
	}
}

// WithRedis returns a function to set the redis of a PeriodLimit
func WithRedis(client *redis.Client) Option {
	return func(p *PeriodLimit) {
		p.redis = client
	}
}
