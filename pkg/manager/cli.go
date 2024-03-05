package manager

import "github.com/Pii-Hole-Engineering/no-phi-ai/pkg/cfg"

// runCLI() method is used to run the command specified in m.config.Command.Run var.
func (m *Manager) runCLI() (e error) {
	switch m.config.Command.Run {
	case cfg.CommandRunHelp:
		e = m.commandHelp()
		return
	case cfg.CommandRunListOrgRepos:
		e = m.commandListOrgRepos()
		return
	case cfg.CommandRunScanOrg:
		e = m.commandScanOrg()
		return
	case cfg.CommandRunScanRepos:
		e = m.commandScanRepos()
		return
	case cfg.CommandRunScanTest:
		e = m.commandScanTest()
		return
	case cfg.CommandRunVersion:
		e = m.commandVersion()
		return
	default:
		m.logger.Fatal().Msgf("invalid command : %s", m.config.Command.Run)
		return
	}
}
