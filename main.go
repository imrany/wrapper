package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"

    "github.com/joho/godotenv"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    pb "github.com/imrany/wrapper/proto/gen/api/v1"
    apiv1 "github.com/imrany/wrapper/router/api/v1"
)

var rootCmd = &cobra.Command{
    Use:   "wrapper",
    Short: "Wrapper is a Gemini gRPC service",
    Run:   runServer,
}

func runServer(_ *cobra.Command, _ []string) {
    ctx := context.Background()

    port := viper.GetInt("port")
    if port == 0 {
        port = 8080
    }

    apiKey := viper.GetString("gemini-api-key")
    if apiKey == "" {
        log.Println("⚠️  No Gemini API key provided")
    }

    addr := fmt.Sprintf("0.0.0.0:%d", port)
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Failed to listen on %s: %v", addr, err)
    }

    grpcServer := grpc.NewServer()

    // Register your service implementation
    pb.RegisterAiServiceServer(grpcServer, &apiv1.GeminiService{APIKey: apiKey})

    // Optional: Enable reflection for easier debugging
    reflection.Register(grpcServer)

    log.Printf("🚀 gRPC server listening on %s", addr)
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("gRPC server failed: %v", err)
    }

    <-ctx.Done()
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  No .env file found or failed to load")
    }

    viper.AutomaticEnv()

    envBindings := map[string]string{
        "port":           "PORT",
        "gemini-api-key": "GEMINI_API_KEY",
    }

    rootCmd.PersistentFlags().Int("port", 8080, "Port to run the gRPC server on")
    rootCmd.PersistentFlags().String("gemini-api-key", "", "Gemini API key")

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
