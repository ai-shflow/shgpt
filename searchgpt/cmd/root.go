package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"searchgpt/review"
)

const (
	columnNum = 4

	defaultCount = 100
	defaultStart = 1

	timeLayout = "2006-01-02"
)

var (
	querySearch string
	queryUser   string

	queryStart int
	queryCount int

	querySince string
	queryUntil string
)

var rootCmd = &cobra.Command{
	Use: fmt.Sprintf("/search <query>"),
	Args: func(cmd *cobra.Command, args []string) error {
		var err error
		if len(args) < 1 {
			now := time.Now()
			querySince = now.Format(timeLayout)
			yyyy, mm, dd := now.Date()
			queryUntil = time.Date(yyyy, mm, dd+1, 0, 0, 0, 0, now.Location()).Format(timeLayout)
			querySearch = fmt.Sprintf("owner:%s since:%s until:%s", queryUser, querySince, queryUntil)
			return nil
		}
		if len(args) == 1 && args[0] == "help" {
			return errors.New("invalid argument\n")
		}
		if querySearch, err = parse(args[0]); err != nil {
			return errors.Wrap(err, "failed to parse\n")
		}
		return nil
	},
	Example: "\n" +
		"  /search\n" +
		fmt.Sprintf("  /search -s %d -c %d\n", defaultStart, defaultCount) +
		"  /search \"project:name branch:master since:2024-01-01 until:2024-01-02\"\n",
	Run: func(cmd *cobra.Command, args []string) {
		if err := execute(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// nolint: gochecknoinits
func init() {
	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().IntVarP(&queryStart, "start", "s", defaultStart, "query start")
	rootCmd.PersistentFlags().IntVarP(&queryCount, "count", "c", defaultCount, "query count")
	rootCmd.PersistentFlags().StringVarP(&queryUser, "user", "u", "", "query user")

	_ = rootCmd.MarkPersistentFlagRequired("user")
}

func parse(query string) (string, error) {
	if strings.Contains(query, "owner:self") {
		if queryUser == "" {
			return "", errors.New("invalid query user")
		}
		return strings.Replace(query, "owner:self", "owner:"+queryUser, -1), nil
	}

	return query, nil
}

func execute() error {
	if queryStart < 1 {
		return errors.New("invalid query start")
	}

	if queryCount < 1 {
		return errors.New("invalid query count")
	}

	if err := queryReview(queryStart, queryCount); err != nil {
		return errors.Wrap(err, "failed to query review")
	}

	return nil
}

func queryReview(start, count int) error {
	r := review.New()
	if r == nil {
		return errors.New("failed to new review")
	}

	changes, err := r.Query(querySearch, start-1, count)
	if err != nil {
		return errors.Wrap(err, "failed to query changes")
	}

	cols := len(changes) / columnNum
	remain := len(changes) % columnNum

	if querySince != "" && queryUntil != "" {
		fmt.Printf("%s: changes: since:%s until:%s\n\n", queryUser, querySince, queryUntil)
	} else {
		if querySince != "" && queryUntil == "" {
			fmt.Printf("%s: changes: since:%s\n\n", queryUser, querySince)
		} else if querySince == "" && queryUntil != "" {
			fmt.Printf("%s: changes: until:%s\n\n", queryUser, queryUntil)
		} else {
			fmt.Printf("%s: changes\n\n", queryUser)
		}
	}

	for i := range cols {
		for index := i * columnNum; index < (i*columnNum + columnNum); index++ {
			number := int(changes[index].(map[string]interface{})["_number"].(float64))
			status := changes[index].(map[string]interface{})["status"].(string)
			fmt.Printf("%d (%s) ", number, strings.ToLower(status))
		}
		fmt.Println()
	}

	for i := range remain {
		index := i + (cols * columnNum)
		number := int(changes[index].(map[string]interface{})["_number"].(float64))
		status := changes[index].(map[string]interface{})["status"].(string)
		fmt.Printf("%d (%s) ", number, strings.ToLower(status))
	}

	return nil
}
