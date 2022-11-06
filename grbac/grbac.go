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

// defines a set of errors
var (
	ErrInvalidRequest  = errors.New("invalid request")
	ErrUndefinedLoader = errors.New("loader undefined")
)

// Controller defines the structure of the controller
type Controller struct {
	cron         *time.Ticker
	loader       func() (Rules, error)
	loadInterval time.Duration

	rules     Rules
	rulesLock sync.RWMutex

	tree     *tree.Tree
	treeLock sync.RWMutex

	logger    *zlog.Logger
	matchMode meta.MatchMode
}

// ControllerOption provides an interface for user to define controller.
type ControllerOption func(*Controller) error

// WithMatchMode is used to modify the default match mode
func WithMatchMode(mode meta.MatchMode) ControllerOption {
	return func(c *Controller) error {
		c.matchMode = mode
		return nil
	}
}

// WithJSON is used to load configuration via json file
func WithJSON(name string, loadInterval time.Duration) ControllerOption {
	return func(c *Controller) error {
		fd, err := NewJSONLoader(name)
		if err != nil {
			return err
		}
		c.loader = fd.Load
		c.loadInterval = loadInterval
		return nil
	}
}

// WithFile is used to load configuration via loacl file
func WithFile(name string, loadInterval time.Duration) ControllerOption {
	return func(c *Controller) error {
		fd, err := NewFileLoader(name)
		if err != nil {
			return err
		}
		c.loader = fd.Load
		c.loadInterval = loadInterval
		return nil
	}
}

// WithYAML is used to load configuration via yaml file
func WithYAML(name string, loadInterval time.Duration) ControllerOption {
	return func(c *Controller) error {
		fd, err := NewYAMLLoader(name)
		if err != nil {
			return err
		}
		c.loader = fd.Load
		c.loadInterval = loadInterval
		return nil
	}
}

// WithRules is used to load config via user defined rules
func WithRules(rules Rules) ControllerOption {
	return func(c *Controller) error {
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
func WithLoader(loader func() (Rules, error), loadInterval time.Duration) ControllerOption {
	return func(c *Controller) error {
		if loader == nil {
			return ErrUndefinedLoader
		}
		c.loader = loader
		c.loadInterval = loadInterval
		return nil
	}
}

// New is used to initialize an RBAC instance
func New(loaderOptions ControllerOption, options ...ControllerOption) (*Controller, error) {
	log := zlog.New("[RBAC]")
	log.ResetFlags(zlog.BitLevel)
	log.SetLogLevel(zlog.LogSuccess)
	c := &Controller{
		logger: log,
	}

	opts := append([]ControllerOption{loaderOptions}, options...)
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
func (c *Controller) SetLogger(logger *zlog.Logger) {
	if logger != nil {
		c.logger = logger
	}
}

func (c *Controller) reload() error {
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

func (c *Controller) buildTree() error {
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

func (c *Controller) runCronTab() {
	if c.loadInterval < time.Second && c.loadInterval >= 0 {
		c.loadInterval = 5 * time.Second
	}
	if c.loadInterval < 0 {
		c.logger.Warn("grbac abandoned the periodic loader because loadInterval is less than 0")
		return
	}

	ticker := time.NewTicker(c.loadInterval)
	c.cron = ticker

	for range ticker.C {
		c.logger.Debug("grbac loader is scheduled")
		err := c.reload()
		if err != nil {
			c.logger.Error("error occurred while loading the configuration in grbac: ", err)
		}
	}
}

func getQueryByRequest(r *http.Request) *Query {
	if r.URL == nil {
		return nil
	}
	return &Query{
		Path:   r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
	}
}

func (c *Controller) find(query *Query) (Rules, error) {
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

// IsRequestGranted is used to verify whether a request has permission
func (c *Controller) IsRequestGranted(r *http.Request, roles []string) (PermissionState, error) {
	query := getQueryByRequest(r)
	if query == nil {
		return meta.PermissionUnknown, ErrInvalidRequest
	}
	return c.IsQueryGranted(query, roles)
}

// IsQueryGranted allows query permissions with the given Query parameter
func (c *Controller) IsQueryGranted(q *Query, roles []string) (PermissionState, error) {
	rules, err := c.find(q)
	if err != nil {
		return meta.PermissionUnknown, err
	}

	return rules.IsRolesGranted(roles, c.matchMode)
}
