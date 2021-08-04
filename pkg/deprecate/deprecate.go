package deprecate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arr-ai/wbnf/parser"
	"github.com/sirupsen/logrus"

	"github.com/arr-ai/arrai/pkg/buildinfo"
)

type sourceContextCache struct {
	sync.RWMutex
	m map[string]struct{}
}

func (d *sourceContextCache) encountered(scanner parser.Scanner) bool {
	s := scanner.String()
	d.RLock()
	if _, has := d.m[s]; has {
		defer d.RUnlock()
		return true
	}
	d.RUnlock()

	d.Lock()
	defer d.Unlock()
	d.m[s] = struct{}{}
	return false
}

func newDeprecatorMap() *sourceContextCache {
	return &sourceContextCache{m: make(map[string]struct{})}
}

// Deprecator contains information required to do deprecation checks and deprecation event trigger
type Deprecator struct {
	cache                             *sourceContextCache
	featureDesc                       string
	weakWarning, strongWarning, crash time.Time
}

var (
	sleepDurationSync sync.Once
	sleepDuration     time.Duration
)

const (
	delayDurationStr = "5s"
	inputDateFormat  = "2006-01-02"
)

func delayDuration() time.Duration {
	sleepDurationSync.Do(func() {
		d, err := time.ParseDuration(delayDurationStr)
		if err != nil {
			panic(err)
		}
		sleepDuration = d
	})
	return sleepDuration
}

// MustNewDeprecator returns a deprecator and panics if deprecator creation fails.
func MustNewDeprecator(featureDesc, weakWarningDate, strongWarningDate, crashDate string) *Deprecator {
	d, err := NewDeprecator(featureDesc, weakWarningDate, strongWarningDate, crashDate)
	if err != nil {
		panic(err)
	}
	return d
}

// NewDeprecator takes brief feature description and three dates with different level of deprecations. It returns a
// deprecator. Dates are considered to be in UTC timezone.
func NewDeprecator(featureDesc, weakWarningDate, strongWarningDate, crashDate string) (*Deprecator, error) {
	weak, err := time.Parse(inputDateFormat, weakWarningDate)
	if err != nil {
		return nil, err
	}
	strong, err := time.Parse(inputDateFormat, strongWarningDate)
	if err != nil {
		return nil, err
	}
	crash, err := time.Parse(inputDateFormat, crashDate)
	if err != nil {
		return nil, err
	}

	if !weak.Before(strong) || !strong.Before(crash) {
		return nil, fmt.Errorf(
			"weak, strong, and crash versions are not sorted: %q, %q, %q",
			weakWarningDate,
			strongWarningDate,
			crashDate,
		)
	}

	return &Deprecator{
		featureDesc:   featureDesc,
		weakWarning:   weak,
		strongWarning: strong,
		crash:         crash,
		cache:         newDeprecatorMap(),
	}, nil
}

// Deprecate does extracts build information from context once and does deprecation checks and trigger deprecation
// events.
func (d *Deprecator) Deprecate(ctx context.Context, scanner parser.Scanner) error {
	if d.cache.encountered(scanner) {
		return nil
	}

	buildDate, err := buildinfo.BuildDateFrom(ctx)
	if err != nil {
		return err
	}

	switch {
	case buildDate == nil:
		logrus.Warnf("%s is being deprecated\n%s", d.featureDesc, scanner.Context(parser.DefaultLimit))
		return nil
	case d.crash.Before(*buildDate):
		// doesn't need to show source context, expression error will handle that.
		return fmt.Errorf("%s is deprecated", d.featureDesc)
	case d.strongWarning.Before(*buildDate):
		logrus.Warnf(
			"%s is being deprecated (pausing %s...)\n%s",
			d.featureDesc, delayDurationStr, scanner.Context(parser.DefaultLimit),
		)
		time.Sleep(delayDuration())
		return nil
	case d.weakWarning.Before(*buildDate):
		logrus.Warnf("%s is being deprecated\n%s", d.featureDesc, scanner.Context(parser.DefaultLimit))
	}
	return nil
}
