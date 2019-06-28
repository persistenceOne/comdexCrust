package keys

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	
	"github.com/comdex-blockchain/crypto/keys"
	"github.com/spf13/cobra"
)

// CheckPassWordCommand :
func CheckPassWordCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check <name>",
		Short: "check password for the given key",
		RunE:  runCheckCmd,
		Args:  cobra.ExactArgs(2),
	}
	return cmd
}
func runCheckCmd(cmd *cobra.Command, args []string) error {
	name := args[0]
	password := args[1]
	kb, err := GetKeyBase()
	if err != nil {
		return err
	}
	
	_, err = kb.ExportPrivateKeyObject(name, password)
	if err != nil {
		return err
	}
	fmt.Println("message: Password is correct")
	return nil
}

// Body :
type Body struct {
	Name     string `json:"name"`
	Password string `json:"passWord"`
}

// CheckPasswordHandler : It will check the password for a key
func CheckPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var m Body
	var kb keys.Keybase
	body, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(body, &m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	kb, err = GetKeyBase()
	_, err = kb.ExportPrivateKeyObject(m.Name, m.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"error":%t ,"message":"%v"}`, true, err.Error())))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"error":%t,"message":"correct password"}`, false)))
		return
	}
}
