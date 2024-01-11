package internal

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	pb "github.com/srmocher/gazelle-kotlin/kotlin/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type KotlinParser struct {
	parserServerBinaryPath string
	parserServerBinaryEnv  []string
	parserServerCmd        *exec.Cmd
	parserServerPid        int

	parserClient pb.KotlinParserClient
}

// getServerBinaryPathAndEnv returns the path to the Kotlin parser server binary and the environment variables
// using Bazel runfiles
func getServerBinaryPathAndEnv() (string, []string, error) {
	r, err := runfiles.New()
	if err != nil {
		return "", nil, fmt.Errorf("Error initializing runfiles: %s", err)
	}
	path, err := r.Rlocation("gazelle-kotlin/kotlin/src/com/github/srmocher/gazelle_kotlin/kotlinparser/parser_server")
	if err != nil {
		return "", nil, fmt.Errorf("Could not find runfiles for server binary: %s", err)
	}
	return path, r.Env(), nil
}

// getServerPort waits for the server to be ready and returns the port it's running on
// it makes upto 5 attempts to retrieve it and waits 2 seconds between each attempt
func (kp *KotlinParser) getServerPort() (string, error) {
	maxAttempts := 5

	portFile := filepath.Join("/tmp", "gazelle-kotlin", fmt.Sprintf("kotlinparser.%d.port", kp.parserServerPid))
	for i := 0; i < maxAttempts; i++ {
		if _, err := os.Stat(portFile); err == nil {
			f, err := os.Open(portFile)
			if err != nil {
				return "", err
			}
			b, err := io.ReadAll(f)
			if err != nil {
				return "", err
			}
			return strings.TrimSuffix(string(b), "\n"), nil
		} else {
			// server not ready yet
			log.Println("Server not ready yet, waiting 2 seconds")
			time.Sleep(2 * time.Second)
		}
	}
	return "", fmt.Errorf("could not find server port file after %d attempts", maxAttempts)
}

// startServer starts the Kotlin parser server without waiting for it to complete
func (kp *KotlinParser) startServer() error {
	cmd := exec.Command(kp.parserServerBinaryPath)
	cmd.Env = kp.parserServerBinaryEnv
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Could not start parser server: %s", err)
	}

	kp.parserServerCmd = cmd
	kp.parserServerPid = cmd.Process.Pid
	log.Printf("Starting server with pid %d", cmd.Process.Pid)
	return nil
}

// Stop stops the Kotlin parser server by sending a sigterm to the process
func (kp *KotlinParser) Stop() error {
	// run server in background
	process, err := os.FindProcess(kp.parserServerPid)
	if err != nil {
		return fmt.Errorf("Could not find server process with pid %d: %s", kp.parserServerPid, err)
	}

	if err = process.Signal(os.Interrupt); err != nil {
		return fmt.Errorf("could not send interrupt signal to server process with pid %d: %s", kp.parserServerPid, err)
	}
	return nil
}

func NewKotlinParser() (*KotlinParser, error) {
	serverPath, env, err := getServerBinaryPathAndEnv()
	if err != nil {
		return nil, err
	}

	kp := KotlinParser{
		parserServerBinaryPath: serverPath,
		parserServerBinaryEnv:  env,
	}
	if err := kp.startServer(); err != nil {
		return nil, err
	}
	// Wait for server to be ready and collect the port it's running on
	port, err := kp.getServerPort()
	if err != nil {
		return nil, fmt.Errorf("failed to get server port: %v", err)
	}
	kotlinParserAddr := fmt.Sprintf("localhost:%s", port)
	log.Printf("Connecting to server at %s\n", kotlinParserAddr)
	conn, err := grpc.Dial(kotlinParserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewKotlinParserClient(conn)
	kp.parserClient = c
	return &kp, nil
}

// ParseKotlinFiles parses the given Kotlin files and returns a list of SourceFileInfo objects
func (kp *KotlinParser) ParseKotlinFiles(files []string) ([]*pb.SourceFileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := kp.parserClient.ParseKotlinFiles(ctx, &pb.KotlinParserRequest{
		KotlinSourceFile: files,
	})

	if err != nil {
		return nil, err
	}

	if r.GetError() != nil {
		return nil, fmt.Errorf("Failed to parse Kotlin files: %s", r.GetError().GetMessage())
	}

	return r.GetSourceFileInfos(), nil
}
