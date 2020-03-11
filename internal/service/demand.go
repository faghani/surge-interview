package service

import (
	"encoding/json"
	"fmt"
	"github.com/faghani/surge-interview/internal/database"
	redispkg "github.com/garyburd/redigo/redis"
	natspkg "github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type NotificationMessage struct {
	RelationId string `json:"relation_id"`
	Coefficient float64 `json:"coefficient"`
}

type DemandService struct {
	sync.Mutex
	Redis      redispkg.Conn
	Nats       *natspkg.Conn
	Db         *database.Database
	Timestep   time.Duration
	Expiration time.Duration
}

func (d *DemandService) Increment(relationId string) (err error) {
	now := time.Now()

	d.Lock()
	key := d.key(relationId, now)
	d.Redis.Send("MULTI")
	d.Redis.Send("SET", key, 0, "PX", d.Expiration, "NX")
	d.Redis.Send("INCR", key)
	_, err = d.Redis.Do("EXEC")
	d.Unlock()

	d.Process(key)

	return nil
}

func (d *DemandService) Process(key string) {
	d.Lock()
	defer d.Unlock()

	var ruleList []database.Rule
	if err := d.Db.FetchRules(&ruleList); err != nil {
		log.WithError(err).Error("failed to get rules from db")
		return
	}

	count, _ := redispkg.Int(d.Redis.Do("GET", key))

	// Find
	var cf float64
	cf = 1
	for _, rule := range ruleList {
		if count >= rule.Threshold {
			cf = rule.Coefficient
		}
	}

	relationId := strings.Split(key, ":")[0]
	cacheKey := relationId + ":last_cf"
	lastNotifiedCoefficient, _ := redispkg.Float64(d.Redis.Do("GET", cacheKey))
	if lastNotifiedCoefficient != cf {
		err := d.Notify(relationId, cf)
		if err != nil {
			log.WithError(err).Error("failed to notify")
		}

		d.Redis.Do("SET", cacheKey, cf)
	}
}

func (d *DemandService) Notify(relationId string, coefficient float64) error {
	data, _ := json.Marshal(NotificationMessage{
		RelationId: relationId,
		Coefficient: coefficient,
	})

	err := d.Nats.Publish("notification.surge.coefficient", data)
	if err != nil {
		return err
	}

	return nil
}

func (d *DemandService) key(r string, t time.Time) string {
	tn := t.UnixNano()
	return fmt.Sprintf("%s:%d", r, d.normalizeTimeInt64(tn))
}

func (d *DemandService) normalizeTimeInt64(t int64) int64 {
	return t - (t % int64(d.Timestep.Nanoseconds()))
}