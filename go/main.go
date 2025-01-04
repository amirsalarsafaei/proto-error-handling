package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	helloworldPb "github.com/amirsalarsafaei/proto-error-handling/autogenerated/go/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/amirsalarsafaei/proto-error-handling/go/internal/helloworld"
)

var defaultMarshaler = protojson.MarshalOptions{
	Indent:    "\t",
	Multiline: true,
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program [command]")
		fmt.Println("Available commands:")
		fmt.Println("  server   Start the gRPC server")
		fmt.Println("  client   Start the gRPC client")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "server":
		serve()
	case "client":
		clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
		username := clientCmd.String("username", "", "Username for the new user")
		email := clientCmd.String("email", "", "Email for the new user")

		err := clientCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println("Error parsing client flags:", err)
			os.Exit(1)
		}

		if *username == "" || *email == "" {
			fmt.Println("Both username and email flags are required")
			clientCmd.PrintDefaults()
			os.Exit(1)
		}

		client(*username, *email)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func serve() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}))
	slog.SetDefault(logger)

	userRepo := helloworld.NewInMemoryUserRepository()
	userService := helloworld.NewUserService(userRepo)

	server := grpc.NewServer()
	helloworldPb.RegisterUserServiceServer(server, userService)

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 8000})
	if err != nil {
		slog.Error("could not listen", slog.Any("error", err))
		return
	}

	go func() {
		err := server.Serve(lis)
		if err != nil {
			slog.Error("could not serve grpc", slog.Any("error", err))
			errChan <- err
		}
	}()

	select {
	case <-sigChan:
		server.GracefulStop()
		slog.Info("gracefully shutdown server")
	case <-errChan:
		slog.Error("server shutdown unexpectedly")
	}

}

func client(username, email string) {
	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	client, err := grpc.NewClient("localhost:8000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error("could not create client", slog.Any("error", err))
		os.Exit(-1)
	}

	userClient := helloworldPb.NewUserServiceClient(
		client,
	)

	resp, err := userClient.CreateUser(ctx, &helloworldPb.CreateUserRequest{
		Username: username,
		Email:    email,
	})
	if spb, ok := status.FromError(err); ok {
		errorJson, err := defaultMarshaler.Marshal(spb.Proto())
		if err != nil {
			slog.Error("could not marshal error", slog.Any("error", err))
			os.Exit(-1)
		}

		fmt.Println("error: ", string(errorJson))
		return
	}

	respJson, err := defaultMarshaler.Marshal(resp)
	if err != nil {
		slog.Error("could not marshal resp", slog.Any("error", err))
		os.Exit(-1)
	}

	fmt.Println("response: ", string(respJson))

}
