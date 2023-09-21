/*
 * This file is part of PaperCrypt.
 *
 * PaperCrypt lets you prepare encrypted messages for printing on paper.
 * Copyright (C) 2023 TMUniversal <me@tmuniversal.eu>.
 *
 * PaperCrypt is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	_ "embed"
	"strings"

	"github.com/tmuniversal/papercrypt/cmd"
	"github.com/tmuniversal/papercrypt/internal"
)

// LicenseText is the license of the application as a string
//
//go:embed COPYING
var LicenseText string

// WordList is the eff.org large word list as a string
//
//go:embed "eff.org_files_2016_07_18_eff_large_wordlist.txt"
var WordList string

// version is the current version of the application
var version = "dev"

// date is the date the application was built
var date = "unknown"

// commit is the git commit hash the application was built from
var commit = "HEAD"

// ref is the git ref the application was built from
var ref = "HEAD"

// branch is the git branch the application was built from
var branch = "HEAD"

// tag is the git tag the application was built from
var tag = "HEAD"

// summary is the git describe summary
var summary = "HEAD"

// repo is the git repository url
var repo = "https://github.com/TMUniversal/PaperCrypt"

// dirty is whether the git repository was dirty when the application was built
var dirty = "false"

// GoVersion is the version of the Go compiler used to build the application
var GoVersion = "unknown"

// arch is the os/arch the application was built for
var arch = "unknown"

// os is the os the application was built for
var os = "unknown"

// builtBy is the tool used to build the application
var builtBy = "go build"

func main() {
	cmd.LicenseText = &LicenseText
	cmd.WordListFile = &WordList

	internal.VersionInfo = internal.VersionDetails{
		Version:    strings.TrimSuffix(version, "\n"),
		BuildDate:  strings.TrimSuffix(date, "\n"),
		GitCommit:  strings.TrimSuffix(commit, "\n"),
		GitRef:     strings.TrimSuffix(ref, "\n"),
		GitBranch:  strings.TrimSuffix(branch, "\n"),
		GitTag:     strings.TrimSuffix(tag, "\n"),
		GitRepo:    repo,
		GitIsDirty: dirty == "true",
		GitSummary: strings.TrimSuffix(summary, "\n"),
		GoVersion:  strings.TrimSuffix(GoVersion, "\n"),
		OsArch:     strings.TrimSuffix(arch, "\n"),
		OsType:     strings.TrimSuffix(os, "\n"),
		BuiltBy:    strings.TrimSuffix(builtBy, "\n"),
	}

	cmd.Execute()
}
