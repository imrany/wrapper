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
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

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

	port := viper.GetInt("port")
	if port == 0 {
		port = 8080
	}

	apiKey := viper.GetString("api-key")
	if apiKey == "" {
		log.Println("⚠️  No API key provided, e.g Gemini api")
	}

	model := viper.GetString("model")
	if model == "" {
		log.Fatal("⚠️  No model provided, e.g gemini-2.0-flash")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterAiServiceServer(grpcServer, &apiv1.APIV1Service{
		APIKey: apiKey,
		Model:  model,
	})
	reflection.Register(grpcServer)

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("🚀 gRPC server listening on %s", addr)
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	// Setup REST gateway
	mux := http.NewServeMux()
	gw := runtime.NewServeMux()
	dialOpts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterAiServiceHandlerFromEndpoint(ctx, gw, addr, dialOpts); err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	mux.Handle("/", gw)
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("proto/gen/api/v1"))))

	// Create HTTP server with proper shutdown support
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8090",
		Handler: withLogging(withCORS(mux)),
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Println("🌐 REST gateway + Swagger UI listening on :8090")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("🛑 Shutting down servers gracefully...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	} else {
		log.Println("✅ HTTP server stopped gracefully")
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()
	log.Println("✅ gRPC server stopped gracefully")
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found or failed to load")
	}

	func withLogging(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("--> %s %s", r.Method, r.RequestURI)
			h.ServeHTTP(w, r)
		})
	}

	viper.AutomaticEnv()

	envBindings := map[string]string{
		"port":    "PORT",
		"api-key": "API_KEY",
		"model":   "MODEL",
	}

	rootCmd.PersistentFlags().Int("port", 8080, "Port to run the gRPC server on")
	rootCmd.PersistentFlags().String("api-key", "", "API key, e.g Gemini API Key")
	rootCmd.PersistentFlags().String("model", "", "Model, e.g Gemini API Model")

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
	}
}
