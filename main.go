package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	
	//steve's custom package to speak for the client to speak with sigma
	pc "github.com/sigma-dev/sigma/plugin/client"
)

type GSMVars struct {
	Name      string `json:"Name"`
	SleepTime string `json:"SleepTime"`
	Output    string `json:"Output"`
	Game      string `json:"Game"`
	NetworkPort   string `json:"NetworkPort"`
}

type GSM struct {
	client *pc.PluginClient
	Vars   GSMVars
	arg    pc.PluginArg
}

// QUAKE 3 PLUGIN and SERVICE

// New() loads the pluginArg and the vars, and returns a GSM struct
// Returns an error if couldn't load vars, or if one of the expected vars is missing
func New() (*GSM, error) {
	h := GSM{}
	h.client = pc.NewPluginClient(http.DefaultClient)

	go h.client.HeartbeatEvery(h.client.HeartbeatTTL / 2)

	h.arg = h.client.GetArg()

	vars, _, err := h.client.LoadVars()
	if err != nil {
		return &h, err
	}
	parsed := GSMVars{}
	err = h.client.UnmarshalVars(vars, "GSM", 0, &parsed)
	if err != nil {
		return &h, err
	}
	h.Vars = parsed

	return &h, nil
}

func (hw *GSM) deploy(op pc.PluginArg) error {
	fmt.Printf("%s: GSM %s\n", op, hw.Vars.Name)
	if err := hw.writeOutputs("deploy"); err != nil {
		return err
	}
	ExecCommand("./linuxgsm install q3server")
	ExecCommand("./q3server auto-install")
	ExecCommand("./q3server start")
	return nil
}

func (hw *GSM) destroy(op pc.PluginArg) error {
	fmt.Printf("%s: GSM %s\n", op, hw.Vars.Name)
	if err := hw.writeOutputs("destroy"); err != nil {
		return err
	}
	ExecCommand("./q3server stop")
	os.Exit(0)
	return nil
}

func (hw *GSM) status(op pc.PluginArg) error {
	fmt.Printf("%s: GSM %s\n", op, hw.Vars.Name)
	if err := hw.writeOutputs("status"); err != nil {
		return err
	}
	return nil
}

func (hw *GSM) update(op pc.PluginArg) error {
	fmt.Printf("%s: GSM %s\n", op, hw.Vars.Name)
	if err := hw.writeOutputs("update"); err != nil {
		return err
	}
	ExecCommand("./q3server update")
	return nil
}

func main() {
	fmt.Println("GSM Plugin")
	hw, err := New()
	if err != nil {
		fail(err.Error(), "Failed to Initialize", 1)
	}

	for {
		if err := hw.Run(); err != nil {
			fail(err.Error(), "", 1)
		}		
	}
}

func (hw *GSM) Run() error {
	switch hw.arg {
	case pc.PluginArgDeploy:
		hw.deploy(hw.arg)
	case pc.PluginArgDestroy:
		hw.destroy(hw.arg)
	case pc.PluginArgStatus:
		hw.status(hw.arg)
	case pc.PluginArgUpdate:
		hw.update(hw.arg)
	case pc.PluginArgUndefined:
		return errors.New("Undefined")
	default:
		return errors.New("Not Impl.")
	}
	return nil
}

// ----------------------- helpers -----------------------------

func (hw *GSM) writeOutputs(op string) error {
	hw.Vars.Output = fmt.Sprintf("%s: GSM %s", op, hw.Vars.Name)
	return hw.client.SendOutput(hw.Vars, "Succeeded") //todo write real status
}

// EXITS with failure.
func fail(err string, progress string, exitCode int) {
	if err != "" {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(exitCode)
}

func ExecCommand(cmd2exec string) {
  	cmd := exec.Command("/bin/sh/", cmd2exec)  	
  	var out bytes.Buffer
  	cmd.Stdout = &out
  	err := cmd.Run()
  	if err != nil {
  		log.Fatal(err)
  	}
  	fmt.Printf("Command output: %q\n", out.String())	
  }

