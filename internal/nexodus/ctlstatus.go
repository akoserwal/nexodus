package nexodus

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"runtime"

	"github.com/urfave/cli/v2"
)

func callNexd(method string) (string, error) {
	conn, err := net.Dial("unix", "/run/nexd.sock")
	if err != nil {
		return "", fmt.Errorf("Failed to connect to nexd: %w\n", err)
	}
	defer conn.Close()

	client := jsonrpc.NewClient(conn)

	var result string
	err = client.Call("NexdCtl."+method, nil, &result)
	if err != nil {
		return "", fmt.Errorf("Failed to execute method (%s): %w\n", method, err)
	}

	return result, nil
}

// CtlStatus attempt to retrieve the status of the nexd service
func CtlStatus(cCtx *cli.Context) (string, error) {
	if runtime.GOOS != Linux.String() {
		return "", fmt.Errorf("nexd ctl interface not yet supported on %s", runtime.GOOS)
	}

	result, err := callNexd("Status")
	if err != nil {
		return "", err
	}

	return result, err
}
