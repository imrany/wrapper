package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/joho/godotenv"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

    pb "github.com/imrany/wrapper/proto/gen/api/v1"
    apiv1 "github.com/imrany/wrapper/router/api/v1"
)

var rootCmd = &cobra.Command{
    Use:   "wrapper",
    Short: "Wrapper is a AI gRPC + REST service",
    Run:   runServer,
}

func runServer(_ *cobra.Command, _ []string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle graceful shutdown
    go func() {
        sigCh := make(chan os.Signal, 1)
        signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
        <-sigCh
        log.Println("🛑 Shutting down server...")
        cancel()
    }()

    port := viper.GetInt("port")
    if port == 0 {
        port = 8080
    }

    apiKey := viper.GetString("api-key")
    if apiKey == "" {
        log.Println("⚠️  No API key provided, e.g Gemini api")
    }

    addr := fmt.Sprintf("0.0.0.0:%d", port)
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Failed to listen on %s: %v", addr, err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterAiServiceServer(grpcServer, &apiv1.APIV1Service{
        APIKey: apiKey,
    })
    reflection.Register(grpcServer)

    // Start gRPC server in a goroutine
    go func() {
        log.Printf("🚀 gRPC server listening on %s", addr)
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("gRPC server failed: %v", err)
        }
    }()

    // Setup REST gateway
    mux := http.NewServeMux()
    gw := runtime.NewServeMux()

    dialOpts := []grpc.DialOption{grpc.WithInsecure()}
    if err := pb.RegisterAiServiceHandlerFromEndpoint(ctx, gw, addr, dialOpts); err != nil {
        log.Fatalf("Failed to register gateway: %v", err)
    }

    mux.Handle("/api/", gw)
    mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("proto/gen/api/v1"))))

    log.Println("🌐 REST gateway + Swagger UI listening on :8090")
    if err := http.ListenAndServe(":8090", mux); err != nil {
        log.Fatalf("HTTP server failed: %v", err)
    }

    <-ctx.Done()
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found or failed to load")
    }

    viper.AutomaticEnv()

    envBindings := map[string]string{
        "port":     "PORT",
        "api-key":  "API_KEY",
    }

    rootCmd.PersistentFlags().Int("port", 8080, "Port to run the gRPC server on")
    rootCmd.PersistentFlags().String("api-key", "", "API key, e.g Gemini API Key")

    for key, env := range envBindings {
        if err := viper.BindEnv(key, env); err != nil {
            panic(fmt.Errorf("failed to bind env var '%s': %w", key, err))
        }
        if err := viper.BindPFlag(key, rootCmd.PersistentFlags().Lookup(key)); err != nil {
            panic(fmt.Errorf("failed to bind flag '%s': %w", key, err))
        }
    }
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("Command failed: %v", err)
        os.Exit(1)
    }
}
