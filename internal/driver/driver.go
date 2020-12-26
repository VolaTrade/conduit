package driver

ithub.com/volatrade/conduit/internal/streamprocessor"
	logger "github.com/volatrade/currie-logs"
	stats "github.com/volatrade/k-stats"
)

var Module = wire.NewSet(
	New,
)

type (
	ConduitDriver struct {
		sp      sproc.StreamProcessor
		kstats  *stats.Stats
		session session.Session
		logger  *logger.Logger
	}
)
xpxp