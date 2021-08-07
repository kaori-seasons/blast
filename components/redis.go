package components

import (
	"github.com/go-redis/redis/v7"
	"github.com/complone/blast/monitor"
	"time"
)

type RedisDB struct {
	hosts                []string
	db                   int
	default_expiration_s int
	retry                int
	rety_inter_ms        int

	instaces []*redis.Client
	cur_idx  int

	monitorKey string
}

func NewRedisDB(hosts []string, db int, expiration int, apiName string) *RedisDB {
	if len(hosts) == 0 {
		panic("hosts of redis must not null")
	}
	r := &RedisDB{
		hosts,
		db,
		expiration,
		3,
		100,
		make([]*redis.Client, len(hosts)),
		0,
		apiName,
	}

	for i, item := range r.hosts {
		c := redis.NewClient(&redis.Options{
			Addr: item,
			DB:   r.db,
		})
		r.instaces[i] = c
	}

	return r
}

func (r *RedisDB) fetchConnection() (*redis.Client, error) {
	var err error
	for idx := r.cur_idx; idx < len(r.instaces); idx++ {
		if _, err = r.instaces[idx].Ping().Result(); err == nil {
			r.cur_idx = idx
			return r.instaces[idx], nil
		}
	}

	for idx := 0; idx < r.cur_idx; idx++ {
		if _, err = r.instaces[idx].Ping().Result(); err == nil {
			r.cur_idx = idx
			return r.instaces[idx], nil
		}
	}

	return nil, err
}

func (r *RedisDB) Set(k string, v string) error {
	monitor.RedisTotalRequests.WithLabelValues(monitor.GetClutser(), r.monitorKey, "set").Inc()
	now := time.Now()
	defer func(begin time.Time) {
		monitor.RedisRequestLatency.WithLabelValues(monitor.GetClutser(), r.monitorKey, "set").Observe(float64(time.Since(begin).Milliseconds()))
	}(now)

	var err error
	for i := 0; i < r.retry; i++ {
		var cli *redis.Client
		cli, err = r.fetchConnection()
		if err != nil {
			goto CONTINUE
		}

		_, err = cli.Set(k, v, time.Duration(r.default_expiration_s)*time.Second).Result()
		if err == nil {
			return nil
		}

	CONTINUE:
		r.sleep(i)
	}

	return err
}

func (r *RedisDB) MSet(values map[string]interface{}) error {
	monitor.RedisTotalRequests.WithLabelValues(monitor.GetClutser(), r.monitorKey, "mset").Inc()
	now := time.Now()
	defer func(begin time.Time) {
		monitor.RedisRequestLatency.WithLabelValues(monitor.GetClutser(), r.monitorKey, "mset").Observe(float64(time.Since(begin).Milliseconds()))
	}(now)

	if len(values) == 0 {
		return nil
	}

	var err error
	for i := 0; i < r.retry; i++ {
		var cli *redis.Client
		var p redis.Pipeliner
		cli, err = r.fetchConnection()
		if err != nil {
			goto CONTINUE
		}

		p = cli.Pipeline()
		defer p.Close()

		for k, v := range values {
			_, err := p.Set(k, v, time.Duration(r.default_expiration_s)*time.Second).Result()
			if err != nil {
				goto CONTINUE
			}
		}

		_, err = p.Exec()
		if err == nil {
			return nil
		}

	CONTINUE:
		r.sleep(i)
	}

	return err
}

func (r *RedisDB) Get(k string) (string, error) {
	monitor.RedisTotalRequests.WithLabelValues(monitor.GetClutser(), r.monitorKey, "get").Inc()
	now := time.Now()
	defer func(begin time.Time) {
		monitor.RedisRequestLatency.WithLabelValues(monitor.GetClutser(), r.monitorKey, "get").Observe(float64(time.Since(begin).Milliseconds()))
	}(now)

	var err error
	var value string
	for i := 0; i < r.retry; i++ {
		var cli *redis.Client
		cli, err = r.fetchConnection()
		if err != nil {
			goto CONTINUE
		}

		value, err = cli.Get(k).Result()
		if err == nil {
			return value, nil
		}

	CONTINUE:
		r.sleep(i)
	}

	return value, err
}

// when internal error occurs, return nil, err
// when key not exist, return len(object) <= len(keys), nil
func (r *RedisDB) MGet(keys []string) ([]interface{}, error) {
	monitor.RedisTotalRequests.WithLabelValues(monitor.GetClutser(), r.monitorKey, "mget").Inc()
	now := time.Now()
	defer func(begin time.Time) {
		monitor.RedisRequestLatency.WithLabelValues(monitor.GetClutser(), r.monitorKey, "mget").Observe(float64(time.Since(begin).Milliseconds()))
	}(now)

	var err error
	var rst []interface{}
	for i := 0; i < r.retry; i++ {
		var cli *redis.Client
		cli, err = r.fetchConnection()
		if err != nil {
			goto CONTINUE
		}
		rst, err = cli.MGet(keys...).Result()
		if err == nil {
			return rst, nil
		}
	CONTINUE:
		r.sleep(i)
	}

	return rst, err
}

func (r *RedisDB) sleep(i int) {
	if i < r.retry-1 {
		time.Sleep(time.Duration(r.rety_inter_ms) * time.Millisecond)
	}
}
