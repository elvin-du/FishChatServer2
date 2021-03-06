package dao

import (
	"github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
	"github.com/oikomi/FishChatServer2/common/dao/xredis"
	"github.com/oikomi/FishChatServer2/server/register/conf"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

const (
	_key     = "re_"
	_online  = 1
	_offline = 0
)

func key(uid int64) string {
	return _key + strconv.FormatInt(uid, 10)
}

type Dao struct {
	redis *xredis.Pool
}

func NewDao() (dao *Dao) {
	dao = &Dao{
		redis: xredis.NewPool(conf.Conf.Redis.Redis),
	}
	return
}

func (dao *Dao) Token(ctx context.Context, uid int64) (res string, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", key(uid)))
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) SetToken(ctx context.Context, uid int64, token string) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key(uid), int(time.Duration(conf.Conf.Redis.Expire)/time.Second), token)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) RegisterAccess(ctx context.Context, uid int64, accessAddr string) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SET", key(uid), accessAddr)
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) RouterAccess(ctx context.Context, uid int64) (res string, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.String(conn.Do("GET", key(uid)))
	if err != nil {
		glog.Error(err)
	}
	return
}

func (dao *Dao) GetOnline(ctx context.Context, uid int64) (res int, err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	res, err = redis.Int(conn.Do("GET", key(uid)))
	if err != nil {
		res = _offline
		glog.Error(err)
	}
	return
}

func (dao *Dao) SetOnline(ctx context.Context, uid int64) (err error) {
	conn := dao.redis.Get(ctx)
	defer conn.Close()
	_, err = conn.Do("SETEX", key(uid), int(time.Duration(conf.Conf.Redis.Expire)/time.Second), _online)
	if err != nil {
		glog.Error(err)
	}
	return
}
