package ldapsearch

import (
  // Standard
  "encoding/json"
  "fmt"
  "strings"
  "strconv"
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
	SizeLimit      string `json:"sizelimit"`

}

var(
   serverAddress = ""
   bindUser = ""
   bindPassword = ""
   searchFilter = ""
   baseSearchDN = ""
   sizeLimit = ""
)

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

	outpt,err := LDAPSearch(args.ServerAddress,args.BindUser,args.BindPassword,args.SearchFilter,args.BaseSearchDN,args.SizeLimit)
        if err != nil {
                msg.UserOutput = err.Error()
                msg.Completed = true
                msg.Status = "error"
                task.Job.SendResponses <- msg
                return
        } else {
		msg.UserOutput = outpt
		msg.Completed = true
		task.Job.SendResponses <- msg
	}
}

func LDAPSearch(serveraddress, binduser, bindpassword, searchfilter, basesearchdn, sizelimit string) (string,error) {

  serverAddress = serveraddress
  bindUser = binduser
  bindPassword = bindpassword
  searchFilter = searchfilter
  baseSearchDN = basesearchdn
  sizeLimit = sizelimit

  conn, err := establishConnection()

  if err != nil {
    return "nil",err
  }

  defer conn.Close()


    sizelimitint, err := strconv.Atoi(sizeLimit)
    if err != nil {
        return "nil",err
    }


   rslts,err := listEntries(conn,sizelimitint)
    if err != nil {
    return "nil",err
  }

  if err := authenticateUser(conn); err != nil {
    return "nil",err
  }
return rslts,nil
}


func establishConnection() (*ldap.Conn, error) {

    var tlsConfig *tls.Config

    tlsConfig = &tls.Config{InsecureSkipVerify: true}

  conn, err := ldap.DialTLS("tcp", serverAddress, tlsConfig)

  if err != nil {
    return nil, fmt.Errorf("Connection failed. %s", err)
  }

  if err := conn.Bind(bindUser, bindPassword); err != nil {
    return nil, fmt.Errorf("Bind failed. %s", err)
  }

  return conn, nil
}


func listEntries(conn *ldap.Conn,sizelimit int) (string,error) {
  result, err := conn.Search(ldap.NewSearchRequest(
    baseSearchDN,
    ldap.ScopeWholeSubtree,
    ldap.NeverDerefAliases,
    sizelimit,
    0,
    false,
    buildFilter("*"),
    []string{"*"},
    nil,
  ))


  if err != nil {
    return "error",fmt.Errorf("Search failed. %s", err)
  }

    var entries []string

    for _, entry := range result.Entries {
        var entryDetails []string
        entryDetails = append(entryDetails, fmt.Sprintf("DN: %s", entry.DN))

        for _, attr := range entry.Attributes {

            attrValues := strings.Join(attr.Values, ", ")
            entryDetails = append(entryDetails, fmt.Sprintf("%s: %s", attr.Name, attrValues))
        }

        entryString := strings.Join(entryDetails, "\n")
        entries = append(entries, entryString)
    }

        entriesString := strings.Join(entries, "\n\n")

        return entriesString, nil

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

  if err := conn.Bind(bindUser, bindPassword); err != nil {
    return fmt.Errorf("Authentication failed. %s", err)
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
