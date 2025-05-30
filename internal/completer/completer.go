package completer

import (
	"fmt"
	"sort"
	"strings"
)

// CmdNode represents a node in the command tree
type CmdNode struct {
	Children map[string]*CmdNode
	IsValue  bool // Whether this is a value node
}

var rootCmdNode = &CmdNode{
	Children: map[string]*CmdNode{
		"set": {
			Children: map[string]*CmdNode{
				"dns": {IsValue: true},
				"interfaces": {
					Children: map[string]*CmdNode{
						"eth0": {
							Children: map[string]*CmdNode{
								"address": {IsValue: true},
								"mac":     {IsValue: true},
							},
						},
						"eth1": {
							Children: map[string]*CmdNode{
								"address": {IsValue: true},
								"mac":     {IsValue: true},
							},
						},
					},
				},
			},
		},
		"add": {
			Children: map[string]*CmdNode{
				"dns": {IsValue: true},
			},
		},
		"show": {
			Children: map[string]*CmdNode{
				"dns":        {},
				"config":     {},
				"interfaces": {},
				"version":    {},
			},
		},
		"save": {},
		"exit": {},
		"help": {},
		"?":    {},
	},
}

// CLICompleter implements readline.AutoCompleter for tab completion
type CLICompleter struct{}

// Do implements readline.AutoCompleter interface
func (c *CLICompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Find the start of the current word
	start := pos
	for start > 0 && line[start-1] != ' ' {
		start--
	}
	prefix := string(line[start:pos])
	tokens := strings.Fields(string(line[:start]))

	// If ending with space, get next level candidates without prefix
	if pos > 0 && line[pos-1] == ' ' {
		tokens = append(tokens, "")
		prefix = ""
		start = pos
	}

	// Show candidate list with ?
	if prefix == "?" {
		candidates := getCompletionsStrict(tokens, "")
		if len(candidates) > 0 {
			fmt.Println()
			for _, cand := range candidates {
				fmt.Println("  " + cand)
			}
		}
		return nil, pos
	}

	// Get completion candidates
	completions := getCompletionsStrict(tokens, prefix)
	if len(completions) == 0 {
		return nil, pos
	}

	// Sort candidates
	sort.Strings(completions)

	// Return results
	var result [][]rune
	for _, comp := range completions {
		// Add space to all completion candidates
		comp = comp + " "
		// Return the completion word itself (readline will remove the prefix part)
		result = append(result, []rune(comp[len(prefix):]))
	}
	return result, start
}

// getCompletionsStrict: traverse the tree with tokens except prefix, only return prefix matches
func getCompletionsStrict(tokens []string, prefix string) []string {
	node := rootCmdNode
	for _, t := range tokens {
		if t == "" {
			// Empty token doesn't descend the tree (right after space)
			break
		}
		if node.Children == nil {
			return nil
		}
		child, ok := node.Children[t]
		if !ok {
			return nil
		}
		node = child
	}

	// Don't complete if there are no child nodes
	if node.Children == nil {
		return nil
	}

	// Collect candidates
	var res []string
	for k := range node.Children {
		// Only add candidates that match the prefix
		if strings.HasPrefix(k, prefix) {
			res = append(res, k)
		}
	}

	// If there's only one candidate, treat it as an exact match
	if len(res) == 1 {
		return res
	}

	// If there are multiple candidates, return only prefix matches
	return res
}
