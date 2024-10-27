package tasks

import (
	"fmt"
	"os"

	"github.com/cligpt/shflow/model"
	"github.com/cligpt/shflow/utils"
)

func Flow() model.Flow {
	// 引入config中的agents和task
	pathagentsmd := utils.GetPath("agents.md")
	pathtasksmd := utils.GetPath("tasks.md")
	agentsMd := pathagentsmd
	tasksMd := pathtasksmd

	// 读取markdown内容
	contentAgent, err := os.ReadFile(agentsMd)
	if err != nil {
		fmt.Println(err)
	}
	contentTask, errTask := os.ReadFile(tasksMd)
	if errTask != nil {
		fmt.Println(errTask)
	}
	flow := utils.ParseMarkdown(string(contentAgent), string(contentTask))
	return flow
}
