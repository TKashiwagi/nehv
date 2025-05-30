package cmd

import (
	"fmt"
	"sort"
	"strings"
)

// VyOS風コマンド階層ツリー（ネスト構造）
type CmdNode struct {
	Children map[string]*CmdNode
	IsValue  bool // 値ノードか
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

// cliCompleter implements readline.AutoCompleter for tab completion
type cliCompleter struct{}

func (c *cliCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// 現在の単語の先頭を探す
	start := pos
	for start > 0 && line[start-1] != ' ' {
		start--
	}
	prefix := string(line[start:pos])
	tokens := strings.Fields(string(line[:start]))

	// スペースで終わっている場合はprefixなしで次階層候補
	if pos > 0 && line[pos-1] == ' ' {
		tokens = append(tokens, "")
		prefix = ""
		start = pos
	}

	// ? で候補一覧表示
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

	// 補完候補を取得
	completions := getCompletionsStrict(tokens, prefix)
	if len(completions) == 0 {
		return nil, pos
	}

	// 候補をソート
	sort.Strings(completions)

	// 結果を返す
	var result [][]rune
	for _, comp := range completions {
		// すべての補完候補にスペースを追加
		comp = comp + " "
		// 補完語そのものを返す（prefix部分はreadlineが削除する）
		result = append(result, []rune(comp[len(prefix):]))
	}
	return result, start
}

// getCompletionsStrict: prefix以外のトークンでツリーを降り、prefix一致のみ候補
func getCompletionsStrict(tokens []string, prefix string) []string {
	node := rootCmdNode
	for _, t := range tokens {
		if t == "" {
			// 空トークンは階層を降りない（スペース直後）
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

	// 子ノードがない場合は補完しない
	if node.Children == nil {
		return nil
	}

	// 候補を収集
	var res []string
	for k := range node.Children {
		// 部分一致の場合のみ候補に追加
		if strings.HasPrefix(k, prefix) {
			res = append(res, k)
		}
	}

	// 候補が1つの場合は完全一致として扱う
	if len(res) == 1 {
		return res
	}

	// 候補が複数の場合は部分一致のみを返す
	return res
}
