package grbac

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"zlsapp/grbac/meta"
	"zlsapp/grbac/pkg/tree"

	"github.com/sohaha/zlsgo/zlog"
)

var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUndefinedLoader = errors.New("loader undefined")
)

type Engine struct {
	cron         *time.Ticker
	loader       func() (Rules, error)
	tree         *tree.Tree
	logger       *zlog.Logger
	rules        Rules
	loadInterval time.Duration
	matchMode    meta.MatchMode
	rulesLock    sync.RWMutex
	treeLock     sync.RWMutex
}

type Option func(*Engine) error

func WithMatchMode(mode meta.MatchMode) Option {
	return func(c *Engine) error {
		c.matchMode = mode
		return nil
	}
}

func WithFile(name string, loadInterval time.Duration) Option {
	return func(c *Engine) error {
		fd, err := NewFileLoader(name)
		if err != nil {
			return err
		}
		c.loader = fd.Load
		c.loadInterval = loadInterval
		return nil
	}
}

// WithRules is used to load config via user defined rules
func WithRules(rules Rules) Option {
	return func(c *Engine) error {
		fd, err := NewRulesLoader(rules)
		if err != nil {
			return nil
		}
		c.loader = fd.Load
		c.loadInterval = -1
		return nil
	}
}

// WithLoader provides a custom Loader entry that you can use to load arbitrary storage.
func WithLoader(loader func() (Rules, error), loadInterval time.Duration) Option {
	return func(c *Engine) error {
		if loader == nil {
			return ErrUndefinedLoader
		}
		c.loader = loader
		c.loadInterval = loadInterval
		return nil
	}
}

// New is used to initialize an RBAC instance
func New(loaderOptions Option, options ...Option) (*Engine, error) {
	log := zlog.New("[RBAC]")
	log.ResetFlags(zlog.BitLevel)
	log.SetLogLevel(zlog.LogSuccess)
	c := &Engine{
		logger: log,
	}

	opts := append([]Option{loaderOptions}, options...)
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	if c.loader == nil {
		return nil, ErrUndefinedLoader
	}

	err := c.reload()
	if err != nil {
		return nil, err
	}

	go c.runCronTab()

	return c, nil
}

// SetLogger is used to modify the default logger
func (c *Engine) SetLogger(logger *zlog.Logger) {
	if logger != nil {
		c.logger = logger
	}
}

func (c *Engine) reload() error {
	if c.loader == nil {
		return ErrUndefinedLoader
	}

	rules, err := c.loader()
	if err != nil {
		return err
	}

	err = rules.IsValid()
	if err != nil {
		return err
	}

	c.rulesLock.Lock()
	c.rules = rules
	c.rulesLock.Unlock()

	err = c.buildTree()
	if err != nil {
		return err
	}

	return nil
}

func (c *Engine) buildTree() error {
	t := tree.NewTree()
	c.rulesLock.RLock()
	defer c.rulesLock.RUnlock()
	for _, rule := range c.rules {
		t.Insert(rule.GetArguments(), rule)
	}
	c.treeLock.Lock()
	c.tree = t
	c.treeLock.Unlock()
	return nil
}

func (c *Engine) runCronTab() {
	if c.loadInterval < time.Second && c.loadInterval >= 0 {
		c.loadInterval = 5 * time.Second
	}
	if c.loadInterval < 0 {
		return
	}

	ticker := time.NewTicker(c.loadInterval)
	c.cron = ticker

	for range ticker.C {
		err := c.reload()
		if err != nil {
			c.logger.Error("error occurred while loading the configuration:", err)
		}
	}
}

func (c *Engine) find(query *Query) (Rules, error) {
	c.treeLock.RLock()
	defer c.treeLock.RUnlock()
	records, err := c.tree.Query(query.GetArguments())
	if err != nil {
		return nil, err
	}
	var perms Rules
	for _, record := range records {
		perm, ok := record.(*Rule)
		if !ok {
			continue
		}
		perms = append(perms, perm)
	}
	return perms, nil
}

func (c *Engine) IsRequestGranted(r *http.Request, roles []string) (PermissionState, error) {
	query, err := c.NewQueryByRequest(r)
	if err == nil {
		return meta.PermissionUnknown, err
	}

	return query.IsRolesGranted(roles)
}
