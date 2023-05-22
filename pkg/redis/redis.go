package redis

import (
	"net/url"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"github.com/redis/rueidis"
	"github.com/redis/rueidis/rueidislock"
)

type Options struct {
	Url        string //revive:disable-line
	LockPrefix string
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("redis", o); err != nil {
		return nil, err
	}
	return o, nil
}

func NewRedisOption(opt *Options) (rueidis.ClientOption, error) {
	u, err := url.Parse(opt.Url)
	if err != nil {
		return rueidis.ClientOption{}, err
	}
	username := u.User.Username()
	if username == "" {
		username = "default"
	}
	password, _ := u.User.Password()
	db, err := strconv.Atoi(strings.TrimPrefix(u.Path, "/"))
	if err != nil {
		return rueidis.ClientOption{}, err
	}
	return rueidis.ClientOption{
		Username:          username,
		Password:          password,
		SelectDB:          db,
		InitAddress:       []string{u.Host},
		CacheSizeEachConn: 1 * bytefmt.GIGABYTE,
	}, nil
}

func New(opt rueidis.ClientOption) (rueidis.Client, error) {
	c, err := rueidis.NewClient(opt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewLock(opt rueidis.ClientOption, o *Options) (rueidislock.Locker, error) {
	return rueidislock.NewLocker(rueidislock.LockerOption{
		ClientOption: opt,
		KeyPrefix:    o.LockPrefix,
	})
}

var ProviderSet = wire.NewSet(New, NewRedisOption, NewOptions, NewLock)
