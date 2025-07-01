// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package version

var (
	// BuildTime is the timestamp at which the binary was compiled at.
	BuildTime string

	// BuildCommit is filled in by the compiler and describe the git reference
	// information at build time.
	BuildCommit string

	// version is the main version number that is being run at the moment. It
	// must conform to the format expected by https://semver.org/.
	version = "0.0.1"

	// versionPrerelease is a pre-release marker for the version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta.1", "rc1.1", etc.
	versionPrerelease = "alpha.1"
)

// Get constructs and returns the Attila version identifier which will include
// any pre-release information.
func Get() string {
	v := version
	if versionPrerelease != "" {
		v += "-" + versionPrerelease
	}
	return v
}
