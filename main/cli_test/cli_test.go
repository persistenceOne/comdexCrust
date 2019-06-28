package clitest

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	
	"github.com/comdex-blockchain/server"
	"github.com/comdex-blockchain/tests"
	"github.com/stretchr/testify/require"
)

var (
	maindHome = ""
)

func init() {
	maindHome = getTestingHomeDir()
}

func TestInitStartSequence(t *testing.T) {
	os.RemoveAll(maindHome)
	servAddr, port, err := server.FreeTCPAddr()
	require.NoError(t, err)
	executeInit(t)
	executeStart(t, servAddr, port)
}

func executeInit(t *testing.T) {
	var (
		chainID string
		initRes map[string]json.RawMessage
	)
	out := tests.ExecuteT(t, fmt.Sprintf("maind --home=%s init", maindHome), "")
	err := json.Unmarshal([]byte(out), &initRes)
	require.NoError(t, err)
	err = json.Unmarshal(initRes["chain_id"], &chainID)
	require.NoError(t, err)
}

func executeStart(t *testing.T, servAddr, port string) {
	proc := tests.GoExecuteTWithStdout(t, fmt.Sprintf("maind start --home=%s --rpc.laddr=%v", maindHome, servAddr))
	defer proc.Stop(false)
	tests.WaitForTMStart(port)
}

func getTestingHomeDir() string {
	tmpDir := os.TempDir()
	maindHome := fmt.Sprintf("%s%s.test_maind", tmpDir, string(os.PathSeparator))
	return maindHome
}
