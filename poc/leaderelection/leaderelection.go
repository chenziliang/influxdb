package leaderelection

import (
	"github.com/coreos/etcd/clientv3"
	"errors"
	"github.com/coreos/etcd/clientv3/concurrency"
	"context"
	"github.com/golang/glog"
	"time"
)

const (
	DefaultSessionTTL = 15 * time.Second
)

// NewLeaderElector creates a LeaderElector from a LeaderElectionConfig
func NewLeaderElector(lec LeaderElectionConfig) (*LeaderElector, error) {
	if lec.Client == nil {
		return nil, errors.New("Client must not be nil.")
	}
	if lec.SessionTTL == 0 {
		lec.SessionTTL = DefaultSessionTTL
	}
	return &LeaderElector{
		config: lec,
	}, nil
}

type LeaderElectionConfig struct {
	Client     *clientv3.Client
	SessionTTL time.Duration
	Election   string
	Identity   string
	// Callbacks are callbacks that are triggered during certain lifecycle
	// events of the LeaderElector
	Callbacks  LeaderCallbacks
}

// LeaderCallbacks are callbacks that are triggered during certain
// lifecycle events of the LeaderElector. These are invoked asynchronously.
//
// possible future callbacks:
//  * OnChallenge()
type LeaderCallbacks struct {
	// OnStartedLeading is called when a LeaderElector client starts leading
	OnStartedLeading func()
	// OnStoppedLeading is called when a LeaderElector client stops leading
	OnStoppedLeading func()
	// OnNewLeader is called when the client observes a leader that is
	// not the previously observed leader. This includes the first observed
	// leader when the client starts.
	OnNewLeader      func(identity string)
}

// LeaderElector is a leader election client.
type LeaderElector struct {
	config         LeaderElectionConfig
}

func Startup(lec LeaderElectionConfig) {
	le, err := NewLeaderElector(lec)
	if err != nil {
		panic(err)
	}
	le.elect()
}

func (le *LeaderElector) elect() {
	// create context
	ctx, cancel := context.WithCancel(context.Background())

	s, err := concurrency.NewSession(
		le.config.Client,
		concurrency.WithTTL(
			int(le.config.SessionTTL.Seconds()),
		),
	)
	if err != nil {
		glog.Error(err)
		return
	}

	e := concurrency.NewElection(s, le.config.Election)

	// register listeners
	le.registerSessionExpirationListener(s, cancel)
	le.registerLeaderChangeListener(ctx, e)

	if err := e.Campaign(ctx, le.config.Identity); err != nil {
		glog.Error(err)
		return
	}

	go le.config.Callbacks.OnStartedLeading()
}

func (le *LeaderElector) registerSessionExpirationListener(s *concurrency.Session, cancel context.CancelFunc) {
	go func() {
		<-s.Done()
		// session expiration
		cancel()
		le.config.Callbacks.OnStoppedLeading()
		le.elect()
	}()
}

func (le *LeaderElector) registerLeaderChangeListener(ctx context.Context, e *concurrency.Election) {
	go func() {
		reportedLeader := ""
		for ol := range e.Observe(ctx) {
			observedLeader := string(ol.Kvs[0].Value)
			if observedLeader == le.config.Identity {
				continue
			}
			if observedLeader == reportedLeader {
				continue
			}
			reportedLeader = observedLeader
			go le.config.Callbacks.OnNewLeader(observedLeader)
		}
	}()
}