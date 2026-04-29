package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/posener/complete/v2"
	"github.com/posener/complete/v2/install"
	"github.com/posener/complete/v2/predict"
)

func predictIdentities(prefix string) []string {
	identities := getAllIdentities()
	var suggestions []string

	for _, identity := range identities {
		if identity.Nickname != "" {
			suggestions = append(suggestions, identity.Nickname)
		}
		suggestions = append(suggestions, identity.Name)
		suggestions = append(suggestions, identity.Email)
	}

	return suggestions
}

func setupCompletion() {
	cmd := &complete.Command{
		Sub: map[string]*complete.Command{
			"list":     {},
			"current":  {},
			"switch":   {Args: complete.PredictFunc(predictIdentities)},
			"use":      {Args: complete.PredictFunc(predictIdentities)},
			"add":      {},
			"delete":   {Args: complete.PredictFunc(predictIdentities)},
			"nickname": {Args: complete.PredictFunc(predictIdentities)},
			"repo": {
				Sub: map[string]*complete.Command{
					"current": {},
					"use":     {Args: complete.PredictFunc(predictIdentities)},
					"add":     {},
				},
			},
			"completion": {
				Sub: map[string]*complete.Command{
					"upgrade": {},
				},
				Args:  predict.Set{"bash", "zsh", "fish"},
				Flags: map[string]complete.Predictor{"r": predict.Nothing},
			},
			"help": {},
		},
		Flags: map[string]complete.Predictor{
			"h":    predict.Nothing,
			"help": predict.Nothing,
		},
	}

	cmd.Complete("gitid")
}

func handleCLICommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	command := args[0]
	switch command {
	case "list":
		return listIdentitiesCLI()
	case "current":
		return getCurrentIdentityCLI()
	case "switch", "use":
		if len(args) < 2 {
			return fmt.Errorf("usage: gitid %s <identifier>", command)
		}
		return switchIdentityCLI(args[1])
	case "add":
		if len(args) < 3 {
			return fmt.Errorf("usage: gitid add <name> <email> [nickname]")
		}
		nickname := ""
		if len(args) > 3 {
			nickname = args[3]
		}
		return addIdentityCLI(args[1], args[2], nickname)
	case "delete", "remove": // also accept 'remove' as an alias for 'delete'
		if len(args) < 2 {
			return fmt.Errorf("usage: gitid delete <identifier>")
		}
		return deleteIdentityCLI(args[1])
	case "nickname":
		if len(args) < 3 {
			return fmt.Errorf("usage: gitid nickname <identifier> <nickname>")
		}
		return setNicknameCLI(args[1], args[2])
	case "completion":
		return completionCLI(args[1:])
	case "repo":
		if len(args) < 2 {
			return fmt.Errorf("usage: gitid repo <current|use|add> [identifier]")
		}
		return repoCLI(args[1:])
	case "help", "--help", "-h":
		showHelp()
		return nil
	default:
		return fmt.Errorf("unknown command: %s\nRun 'gitid help' for usage information", command)
	}
}

func listIdentitiesCLI() error {
	identities := getAllIdentities()
	if len(identities) == 0 {
		fmt.Println("No identities configured.")
		return nil
	}

	var localName, localEmail string
	hasLocal := hasLocalIdentity()
	if hasLocal {
		localName, loca