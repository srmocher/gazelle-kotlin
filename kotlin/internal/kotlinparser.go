package internal

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/bazelbuild/rules_go/go/runfiles"
	pb "github.com/srmocher/gazelle-kotlin/kotlin/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	kotlinParserAddr = "[::1]:50051"
)

type KotlinParser struct {
	parserServerBinaryPath string
	parserServerBinaryEnv  []string
	parserServerPid        int

	parserClient pb.KotlinParserClient
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
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatalf("Could not start parser server: %s", err)
	}

	kp.parserServerPid = cmd.Process.Pid
	log.Printf("Starting server with pid %d", cmd.Process.Pid)
	// run server in background
	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("Server command failed: %s", err)
		}
	}()
	time.Sleep(2)
}

func NewKotlinParser() (*KotlinParser, error) {
	serverPath, env := getServerBinaryPathAndEnv()
	kp := KotlinParser{
		parserServerBinaryPath: serverPath,
		parserServerBinaryEnv:  env,
	}
	kp.startServer()
	conn, err := grpc.Dial(kotlinParserAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c := pb.NewKotlinParserClient(conn)
	kp.parserClient = c
	return &kp, nil
}

func (kp *KotlinParser) ParseKotlinFiles(files []string) ([]*pb.SourceFileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := kp.parserClient.ParseKotlinFiles(ctx, &pb.KotlinParserRequest{
		KotlinSourceFile: files,
	})

	if err != nil {
		return nil, err
	}

	return r.GetSourceFileInfos(), nil
}
