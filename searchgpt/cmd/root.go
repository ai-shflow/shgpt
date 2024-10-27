package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"searchgpt/review"
)

const (
	columnNum = 6

	defaultCount = 100
	defaultStart = 1
)

var (
	queryCount  int
	querySearch string
	queryStart  int
)

var rootCmd = &cobra.Command{
	Use: fmt.Sprintf("/search <query>"),
	Args: func(cmd *cobra.Command, args []string) error {
		var err error
		if len(args) < 1 {
			querySearch = "owner:self"
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
}

func parse(data string) (string, error) {
	return data, nil
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

	for i := range cols {
		for index := i * columnNum; index < columnNum; index++ {
			number := int(changes[index].(map[string]interface{})["_number"].(float64))
			fmt.Printf("%d ", number)
		}
		fmt.Println()
	}

	for i := range remain {
		index := i + (cols * columnNum)
		number := int(changes[index].(map[string]interface{})["_number"].(float64))
		fmt.Printf("%d ", number)
	}

	return nil
}
