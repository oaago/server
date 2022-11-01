package limiter

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oaago/cloud/logx"
	"github.com/oaago/cloud/redis"
)

// to be compatible with aliyun redis, we cannot use `local key = KEYS[1]` to reuse the key
const periodScript = `local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call("INCRBY", KEYS[1], 1)
if current == 1 then
    redis.call("expire", KEYS[1], window)
end
if current < limit then
    return 1
elseif current == limit then
    return 2
else
    return 0
end`

const (
	// Unknown means not initialized state.
	Unknown = iota
	// Allowed means allowed state.
	Allowed
	// HitQuota means this request exactly hit the quota.
	HitQuota
	// OverQuota means passed the quota.
	OverQuota

	internalOverQuota = 0
	internalAllowed   = 1
	internalHitQuota  = 2
)

// ErrUnknownCode is an error that represents unknown status code.
var ErrUnknownCode = errors.New("unknown status code")

type (
	// PeriodOption defines the method to customize a PeriodLimit.
	PeriodOption func(l *PeriodLimit)

	// A PeriodLimit is used to limit requests during a period of time.
	PeriodLimit struct {
		period     int
		quota      int
		limitStore *redis.Cli
		keyPrefix  string
		align      bool
	}
)

func NewRateLimiterUrl(c *gin.Context) {
	// seconds 窗口时常  请求限流quota 没有redis不能使用
	if redis.RedisClient.Client != nil {
		key := "oaaPeriodLimit:url:" + c.Request.Host + ":" + c.Request.Method + ":" + c.Request.URL.String()
		val, err := redis.RedisClient.Client.Get(key).Result()
		if err != nil {
			logx.Logger.Error(err.Error())
			c.Next()
		}
		if len(val) > 0 {
			//格式 时间/速率
			arg := strings.Split(val, ":")
			if len(arg) == 2 {
				period, _ := strconv.Atoi(arg[1])
				quota, _ := strconv.Atoi(arg[0])
				// quota/period 量/时间
				limit := newPeriodLimit(period, quota, redis.RedisClient, key)
				take, _ := limit.Take(key)
				switch take {
				case OverQuota:
					logx.Logger.Info("OverQuota key: %v")
					c.Data(500, "application/json; charset=utf-8", []byte(c.Request.URL.String()+"请求超频:OverQuota"))
					logx.Logger.Info("OverQuota key:"+key, "请求超频 稍后重试")
					c.Abort()
				case Allowed:
					c.Next()
				case HitQuota:
					logx.Logger.Info("HitQuota key: %v")
					c.Data(500, "application/json; charset=utf-8", []byte(c.Request.URL.String()+"请求超频:HitQuota"))
					logx.Logger.Info("HitQuota key:"+key, "请求超频 稍后重试")
					c.Abort()
				default:
					logx.Logger.Info("DefaultQuota key: %v")
					c.Next()
				}
			}
		}
		c.Next()
	}
}

// NewPeriodLimit returns a PeriodLimit with given parameters.
func newPeriodLimit(period, quota int, limitStore *redis.Cli, keyPrefix string,
	opts ...PeriodOption) *PeriodLimit {
	limiter := &PeriodLimit{
		period:     period,
		quota:      quota,
		limitStore: limitStore,
		keyPrefix:  keyPrefix,
	}

	for _, opt := range opts {
		opt(limiter)
	}

	return limiter
}

// Take requests a permit, it returns the permit state.
func (h *PeriodLimit) Take(key string) (int, error) {
	code, err := h.limitStore.Client.Eval(periodScript, []string{h.keyPrefix + key}, []string{
		strconv.Itoa(h.quota),
		strconv.Itoa(h.calcExpireSeconds()),
	}).Int64()
	if err != nil {
		return Unknown, err
	}
	switch code {
	case internalOverQuota:
		return OverQuota, nil
	case internalAllowed:
		return Allowed, nil
	case internalHitQuota:
		return HitQuota, nil
	default:
		return Unknown, ErrUnknownCode
	}
}

func (h *PeriodLimit) calcExpireSeconds() int {
	if h.align {
		now := time.Now()
		_, offset := now.Zone()
		unix := now.Unix() + int64(offset)
		return h.period - int(unix%int64(h.period))
	}

	return h.period
}

// Align returns a func to customize a PeriodLimit with alignment.
// For example, if we want to limit end users with 5 sms verification messages every day,
// we need to align with the local timezone and the start of the day.
func Align() PeriodOption {
	return func(l *PeriodLimit) {
		l.align = true
	}
}
