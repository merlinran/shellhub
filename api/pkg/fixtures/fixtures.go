package fixtures

import (
	"path/filepath"
	"runtime"

	"github.com/shellhub-io/mongotest"
)

const (
	FixtureAnnouncements    = "announcements"     // Check "fixtures.data.announcement" for fixture info
	FixtureConnectedDevices = "connected_devices" // Check "fixtures.data.connected_devices" for fixture info
	FixtureDevices          = "devices"           // Check "fixtures.data.device" for fixture info
	FixtureSessions         = "sessions"          // Check "fixtures.data.session" for fixture info
	FixtureActiveSessions   = "active_sessions"   // Check "fixtures.data.active_session" for fixture info
	FixtureRecordedSessions = "recorded_sessions" // Check "fixtures.data.recorded_session" for fixture info
	FixtureFirewallRules    = "firewall_rules"    // Check "fixtures.data.firewall_rule" for fixture info
	FixturePublicKeys       = "public_keys"       // Check "fixtures.data.public_key" for fixture info
	FixturePrivateKeys      = "private_keys"      // Check "fixtures.data.private_key" for fixture info
	FixtureLicenses         = "licenses"          // Check "fixtures.data.license" for fixture info
	FixtureUsers            = "users"             // Check "fixtures.data.user" for fixture iefo
	FixtureNamespaces       = "namespaces"        // Check "fixtures.data.namespace" for fixture info
	FixtureRecoveryTokens   = "recovery_tokens"   // Check "fixtures.data.recovery_tokens" for fixture info
)

// Init configures the mongotest for the provided host's database. It is necessary
// before using any fixtures and panics if any errors arise.
func Init(host, database string) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to retrieve the fixtures path at runtime")
	}

	mongotest.Configure(mongotest.Config{
		URL:            "mongodb://" + host,
		Database:       database,
		FixtureRootDir: filepath.Join(filepath.Dir(file), "data"),
		FixtureFormat:  mongotest.FixtureFormatJSON,
		PreInsertFuncs: setupPreInsertFuncs(),
	})
}

// Apply applies 'n' fixtures in the database.
func Apply(fixtures ...string) error {
	return mongotest.UseFixture(fixtures...)
}

// Teardown resets all applied fixtures.
func Teardown() error {
	return mongotest.DropDatabase()
}
