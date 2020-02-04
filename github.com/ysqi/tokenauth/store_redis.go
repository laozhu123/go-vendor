// Copyright 2016 Author YuShuangqi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tokenauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"s4s/common/lib/redis"
	"strconv"
)

// Store implement by redis
type RedisStore struct {
	Alias  string
	db     *redis.RedisManager
	dbPath string
}

// Save audience into store.
// Returns error if error occured during execution.
func (store *RedisStore) SaveAudience(audience *Audience) error {

	if audience == nil || len(audience.ID) == 0 {
		return errors.New("audience id is empty.")
	}

	_, err := json.Marshal(audience)
	if err != nil {
		return err
	}

	return nil
}

// Delete audience and  all tokens of audience.
func (store *RedisStore) DeleteAudience(audienceID string) error {
	if len(audienceID) == 0 {
		return errors.New("audienceID is emtpty.")
	}
	return nil
}

// Get audience info or returns error.
func (store *RedisStore) GetAudience(audienceID string) (audience *Audience, err error) {

	if len(audienceID) == 0 {
		return nil, errors.New("audienceID is emtpty.")
	}

	return

}

// Save token to store. return error when save fail.
// Save token json to store and save the relation of token with client if not single model.
// The first , token must not empty and effectiveness.
// Does not consider concurrency.
func (store *RedisStore) SaveToken(token *Token) error {
	if token == nil || len(token.Value) == 0 {
		return errors.New("token tokenString is empty.")
	}
	if len(token.ClientID) == 0 && len(token.SingleID) == 0 {
		return errors.New("token clientid and singleid,It can't be empty")
	}
	if token.Expired() {
		return errors.New("token is expired,not need save.")
	}

	recordKey := fmt.Sprintf("audience_%s_%s", token.ClientID, token.SingleID)

	// delete old auth token
	mp, _ := store.db.HgetAll(recordKey)
	for k := range mp {
		store.db.Del(k)
	}

	//first to get token byte data
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	store.db.Hset(recordKey, token.Value, time.Now().Format("20160102150405"))
	_, err = store.db.SetEx(token.Value, tokenBytes, token.DeadLine-time.Now().Unix())

	return err
}

func (store *RedisStore) FlushToken(token *Token) error {
	if token == nil || len(token.Value) == 0 {
		return errors.New("token tokenString is empty.")
	}
	if len(token.ClientID) == 0 && len(token.SingleID) == 0 {
		return errors.New("token clientid and singleid,It can't be empty")
	}
	if token.Expired() {
		return errors.New("token is expired,not need save.")
	}

	//first to get token byte data
	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return err
	}

	recordKey := fmt.Sprintf("audience_%s_%s", token.ClientID, token.SingleID)
	store.db.Hset(recordKey, token.Value, time.Now().Format("20160102150405"))
	_, err = store.db.SetEx(token.Value, tokenBytes, token.DeadLine-time.Now().Unix())

	return err
}

//Get token info if find in store,or return error
func (store *RedisStore) GetToken(tokenString string) (token *Token, err error) {
	if len(tokenString) == 0 {
		return nil, errors.New("tokenString is empty.")
	}

	reply, err := store.db.Get(tokenString)
	if err != nil {
		return
	}

	data, ok := reply.([]byte)
	if !ok {
		return nil, errors.New("get token failed")
	}

	token = new(Token)
	err = json.Unmarshal(data, token)

	return
}

// Delete token
// Returns error if delete token fail.
func (store *RedisStore) DeleteToken(tokenString string) error {

	if len(tokenString) == 0 {
		return errors.New("incompatible tokenString")
	}

	_, err := store.db.Del(tokenString)

	return err
}

// Close redis db
func (store *RedisStore) Close() error {
	return nil
}

// Delete token if token expired
func (store *RedisStore) DeleteExpired() {}

// Init and Open redis .
// config is json string.
// e.g:
//  {"host":"127.0.0.1:6379", "auth":"123678", "pool_size"="10"}
func (store *RedisStore) Open(config string) error {

	if len(config) == 0 {
		return errors.New("redisStore: redis db store config is empty")
	}

	var cf map[string]string

	err := json.Unmarshal([]byte(config), &cf)
	if err != nil {
		return fmt.Errorf("redisStore: unmarshal %p fail:%s", config, err.Error())
	}

	host, ok := cf["host"]
	if !ok {
		return errors.New("redisStore: redis db store config has no host key.")
	}

	auth := cf["auth"]

	pool_size, _ := strconv.ParseInt(cf["pool_size"], 10, 64)
	if pool_size <= 0 {
		pool_size = 10
	}

	timeoutS, _ := strconv.ParseInt(cf["timeout"], 10, 64)
	timeout := time.Second / 10
	if timeoutS > 0 {
		timeout = time.Second * time.Duration(timeoutS)
	}

	store.db, err = redis.NewRedisManager(host, auth, int(pool_size), timeout)
	return err
}

// new redis store instance.
func NewRedisStore() *RedisStore {
	return &RedisStore{Alias: "RedisStore"}
}

func init() {
	RegStore("default", NewRedisStore())
}
