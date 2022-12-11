package redis

import (
	"net/url"
	"strconv"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/wire"
	"github.com/knadh/koanf"
	"github.com/rueian/rueidis"
	"github.com/rueian/rueidis/rueidislock"
)

type Options struct {
	Url        string
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

func New(o *Options) (rueidis.Client, error) {
	opt, err := NewRedisOption(o)
	if err != nil {
		return nil, err
	}
	c, err := rueidis.NewClient(opt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewLock(o *Options) (rueidislock.Locker, error) {
	opt, err := NewRedisOption(o)
	if err != nil {
		return nil, err
	}
	return rueidislock.NewLocker(rueidislock.LockerOption{
		ClientOption: opt,
		KeyPrefix:    o.LockPrefix,
	})
}

var ProviderSet = wire.NewSet(New, NewOptions, NewLock)
