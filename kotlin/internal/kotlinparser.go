package internal

import (
	"log"
	"os/exec"

	"github.com/bazelbuild/rules_go/go/runfiles"
)

type KotlinParser struct {
	parserServerBinaryPath string
	parserServerBinaryEnv  []string
	parserServerPid        int
}

func getServerBinaryPathAndEnv() (string, []string) {
	r, err := runfiles.New()
	if err != nil {
		log.Fatalf("Error initializing runfiles: %s", err)
	}
	path, err := r.Rlocation("gazelle-kotlin/kotlin/src/com/github/srmocher/gazelle_kotlin/kotlinparser/parser_server")
	if err != nil {
		log.Fatalf("Could not find runfiles for server binary: %s", err)
	}
	return path, r.Env()
}

func (kp *KotlinParser) startServer() {
	cmd := exec.Command(kp.parserServerBinaryPath)
	cmd.Env = kp.parserServerBinaryEnv
	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not start parser server: %s", err)
	}

	kp.parserServerPid = cmd.Process.Pid
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("Server command failed: %s", err)
		}
	}()
}

func NewKotlinParser() *KotlinParser {
	serverPath, env := getServerBinaryPathAndEnv()
	return &KotlinParser{
		parserServerBinaryPath: serverPath,
		parserServerBinaryEnv:  env,
	}
}
