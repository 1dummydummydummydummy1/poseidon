package ldapsearch

import (
  // Standard
  "encoding/json""
  "os"
  "fmt"
  "strings"
  "gopkg.in/ldap.v2"
  "crypto/tls"

	"github.com/MythicAgents/poseidon/Payload_Type/poseidon/agent_code/pkg/utils/structs"
)

type Arguments struct {
	      ServerAddress      string `json:"serveraddress"`
	      BindUser      string `json:"binduser"`
        BindPassword      string `json:"bindpassword"`
        SearchFilter      string `json:"searchfilter"`
        BaseSearchDN      string `json:"basesearchdn"`

}

func Run(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID

	var args Arguments
	err := json.Unmarshal([]byte(task.Params), &args)
	if err != nil {
		msg.UserOutput = err.Error()
		msg.Completed = true
		msg.Status = "error"
		task.Job.SendResponses <- msg
		return
	}

	if LDAPSearch(args.ServerAddress,args.BindUser,args.BindPassword,args.SearchFilter,args.BaseSearchDN)==nil{
		msg.UserOutput = "LDAP query executed successfully"
		msg.Completed = true
		task.Job.SendResponses <- msg
	} else {
		msg.UserOutput = "LDAP query failed"
		msg.Completed = true
		task.Job.SendResponses <- msg
	}

}

func LDAPSearch(serveraddress, binduser, bindpassword, searchfilter, basesearchdn string) error {
serverAddress = serveraddress 
bindUser = binduser
bindPassword = bindpassword
searchFilter = searchfilter
baseSearchDN = basesearchdn


  conn, err := establishConnection()

  if err != nil {
    fmt.Printf("Connection failed. %s", err)
    return
  }

  defer conn.Close()

  if err := listEntries(conn); err != nil {
    fmt.Printf("%v", err)
    return
  }

  if err := authenticateUser(conn); err != nil {
    fmt.Printf("%v", err)
    return
  }
}


func establishConnection() (*ldap.Conn, error) {

    var tlsConfig *tls.Config

    tlsConfig = &tls.Config{InsecureSkipVerify: true}

  conn, err := ldap.DialTLS("tcp", serverAddress, tlsConfig)
  //conn, err := ldap.Dial("tcp", serverAddress)

  if err != nil {
    return nil, fmt.Errorf("Connection failed. %s", err)
  }

  if err := conn.Bind(bindUser, bindPassword); err != nil {
    return nil, fmt.Errorf("Bind failed. %s", err)
  }

  return conn, nil
}

func listEntries(conn *ldap.Conn) error {
  result, err := conn.Search(ldap.NewSearchRequest(
    baseSearchDN,
    ldap.ScopeWholeSubtree,
    ldap.NeverDerefAliases,
    0,
    0,
    false,
    buildFilter("*"),
    []string{"*"},
    nil,
  ))

  if err != nil {
    return fmt.Errorf("Search failed. %s", err)
  }

  // Prints all attributes per entry
  for _, entry := range result.Entries {
    entry.Print()
    fmt.Println()
  }

  return nil
}

func authenticateUser(conn *ldap.Conn) error {
  result, err := conn.Search(ldap.NewSearchRequest(
    baseSearchDN,
    ldap.ScopeWholeSubtree,
    ldap.NeverDerefAliases,
    0,
    0,
    false,
    buildFilter(bindUser),
    []string{"dn"},
    nil,
  ))

  if err != nil {
    return fmt.Errorf("User search failed. %s", err)
  }

  if len(result.Entries) < 1 {
    return fmt.Errorf("User not found")
  }

  if len(result.Entries) > 1 {
    return fmt.Errorf("Multiple users found")
  }

  if err := conn.Bind(bindUser, bindPassword); err != nil {
    fmt.Printf("Authentication failed. %s", err)
  } else {
    fmt.Printf("Authentication successful!")
  }

  return nil
}

func buildFilter(user string) string {
  res := strings.Replace(
    searchFilter,
    "{username}",
    user,
    -1,
  )

  return res
}
