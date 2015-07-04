package tournaments

import (
	"encoding/json"
	"errors"
	"fmt"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/m4rw3r/uuid"
	"os"
	"time"
)

type RedisTournamentStorage struct {
	pool *redigo.Pool
}

func (rts *RedisTournamentStorage) Store(t *Tournament) error {
	conn := rts.pool.Get()
	defer conn.Close()
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if _, err = conn.Do("SADD", "tournaments", t.UUID); err != nil {
		return err
	}
	if _, err = conn.Do("SET", fmt.Sprintf("tournament:%s", t.UUID), b); err != nil {
		return err
	}
	if _, err = conn.Do("SADD", "seasons", t.Info.Season); err != nil {
		return err
	}
	if _, err = conn.Do("SADD", fmt.Sprintf("season:%d:tournaments", t.Info.Season), t.UUID); err != nil {
		return err
	}
	return nil
}

func (rts *RedisTournamentStorage) Load(uuid uuid.UUID) (*Tournament, error) {
	conn := rts.pool.Get()
	defer conn.Close()
	b, err := redigo.Bytes(conn.Do("GET", fmt.Sprintf("tournament:%s", uuid)))
	if err != nil {
		return nil, err
	}
	t := new(Tournament)
	if err := json.Unmarshal(b, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (rts *RedisTournamentStorage) Delete(uuid uuid.UUID) error {
	// FIXME: Not implemented yet
	return errors.New("Not implemented yet")
}

func (rts *RedisTournamentStorage) LoadAll() ([]*Tournament, error) {
	var tournaments []*Tournament
	conn := rts.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", "tournaments"))
	if err != nil {
		return nil, err
	}
	for _, tournament := range b {
		uuid, _ := uuid.FromString(tournament)
		t, err := rts.Load(uuid)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

func (rts *RedisTournamentStorage) LoadBySeason(season int) ([]*Tournament, error) {
	var tournaments []*Tournament
	conn := rts.pool.Get()
	defer conn.Close()
	b, err := redigo.Strings(conn.Do("SMEMBERS", fmt.Sprintf("season:%d:tournaments", season)))
	if err != nil {
		return nil, err
	}
	for _, tournament := range b {
		uuid, _ := uuid.FromString(tournament)
		t, err := rts.Load(uuid)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, t)
	}
	return tournaments, nil
}

func NewRedisTournamentStorage() *RedisTournamentStorage {
	rts := new(RedisTournamentStorage)
	rts.pool = &redigo.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			return redigo.Dial("tcp", os.Getenv("CKPT_REDIS"))
		},
	}
	return rts
}
