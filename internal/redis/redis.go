package redis

import (
	"github.com/garyburd/redigo/redis"
)

// Connect to redis
func New(adr string) (redis.Conn, error) {
	rc, err := redis.Dial("tcp", adr)
	if err != nil {
		return nil, err
	}
	return rc, nil
}