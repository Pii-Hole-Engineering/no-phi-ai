package manager

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/client/az"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/dryrun"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/memory"
	"github.com/Pii-Hole-Engineering/no-phi-ai/pkg/scanner/rrr"
)

// commandHelp() method is used to run the "help" (default) command.
func (m *Manager) commandHelp() (e error) {
	fmt.Printf("CLI Help Information for %s app:\n", m.config.App.Name)
	fmt.Println("\tAvailable Commands:")
	printNameAndDescription(
		cfg.CommandRunHelp,
		"Prints (this) help information for the app.",
	)
	printNameAndDescription(
		cfg.CommandRunListOrgRepos,
		"... not yet implemented ...",
	)
	printNameAndDescription(
		cfg.CommandRunScanOrg,
		"... not yet implemented ...",
	)
	printNameAndDescription(
		cfg.CommandRunScanRepos,
		"... work in progress ...",
	)
	printNameAndDescription(
		cfg.CommandRunScanTest,
		"... work in progress ...",
	)
	printNameAndDescription(
		cfg.CommandRunVersion,
		"Prints version information for the app.",
	)
	fmt.Println("\tEnvironment Variables:")
	for _, envVar := range cfg.GetAppEnvVars() {
		printNameAndDescription(envVar, "")
	}
	return
}

// commandListOrgRepos() method is used to run the "list-org-repos" command.
func (m *Manager) commandListOrgRepos() (e error) {
	m.logger.Warn().Msgf("%s commmand is TODO\n", cfg.CommandRunListOrgRepos)
	return
}

// commandScanOrg() method is used to run the "scan-org" command, which
// is applies the "scan-repos" command to all repositories in the organization.
func (m *Manager) commandScanOrg() (e error) {
	m.logger.Warn().Msgf("%s commmand is TODO\n", cfg.CommandRunScanOrg)
	return
}

// commandScanRepos() method is used to run the "scan-repos" command, which
// is used to scan the contents of a single git repository for PHI/PII.
func (m *Manager) commandScanRepos() (e error) {
	m.scanner, e = scanner.NewScanner(m.ctx, &m.config.Git, memory.NewMemoryResultRecordIO(m.ctx))
	if e != nil {
		e = errors.Wrapf(e, "failed to initialize new Scanner for command %s", m.config.Command.Run)
		return
	}

	if len(m.config.Git.Scan.Repositories) == 0 {
		e = errors.New("no repositories specified for scan")
		return
	}

	var ai *az.EntityDetectionAI
	ai, e = az.NewEntityDetectionAI(m.config)
	if e != nil {
		e = errors.Wrapf(e, "failed to initialize new EntityDetectionAI for command %s", m.config.Command.Run)
		return
	}
	az_ai_detector := az.NewAzAiLanguagePhiDetector(ai)
	chan_scan_errors := make(chan error)
	chan_requests := make(chan rrr.Request)
	chan_responses := make(chan rrr.Response)

	go m.scanner.Scan(chan_scan_errors, chan_requests, chan_responses)
	go az_ai_detector.Run(m.ctx, chan_requests, chan_responses)

	// wait for an error to be returned from the scanner
	e = <-chan_scan_errors
	if e != nil {
		e = errors.Wrapf(e, "failed to run command '%s' ", m.config.Command.Run)
		return
	}
	m.logger.Info().Msgf("command '%s' completed successfully", m.config.Command.Run)

	return
}

// commandScanTest() method is used to run the "scan-test" command, which is
// for development use only.
func (m *Manager) commandScanTest() (e error) {
	m.scanner, e = scanner.NewScanner(m.ctx, &m.config.Git, memory.NewMemoryResultRecordIO(m.ctx))
	if e != nil {
		e = errors.Wrapf(e, "failed to initialize new Scanner for command %s", m.config.Command.Run)
		return
	}

	if len(m.config.Git.Scan.Repositories) == 0 {
		e = errors.New("no repositories specified for scan")
		return
	}

	chan_scan_errors := make(chan error)
	chan_requests := make(chan rrr.Request)
	chan_responses := make(chan rrr.Response)
	dry_run_detector := dryrun.NewDryRunPhiDetector()

	go m.scanner.Scan(chan_scan_errors, chan_requests, chan_responses)
	go dry_run_detector.Run(m.ctx, chan_requests, chan_responses)

	// wait for an error to be returned from the scanner
	e = <-chan_scan_errors
	if e != nil {
		e = errors.Wrapf(e, "failed to run command '%s' ", m.config.Command.Run)
		return
	}
	m.logger.Info().Msgf("command '%s' completed successfully", m.config.Command.Run)

	return
}

// commandVersion() method is used to run the "version" command, which prints
// the version information for the app and then exits.
func (m *Manager) commandVersion() (e error) {
	fmt.Printf("%s %s\n", m.config.App.Name, cfg.AppVersion)
	return
}

// printNameAndDescription() helper function is used to print the name and (optional)
// description of something, such as a command or environment variable, to stdout.
func printNameAndDescription(name string, description string) {
	if description != "" {
		fmt.Printf("\t\t- %s : %s\n", name, description)
	} else {
		fmt.Printf("\t\t- %s\n", name)
	}
}
