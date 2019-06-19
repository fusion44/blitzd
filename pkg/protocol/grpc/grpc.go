package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	pb "github.com/frennkie/blitzd/pkg/api/v1"
)

func runGracefulGRPCServer(ctx context.Context, server *grpc.Server, lis net.Listener) {

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Printf("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server as goroutine
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

}

func runGRPCServer(server *grpc.Server, lis net.Listener) {
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// RunServer runs gRPC service to publish
func RunServer(ctx context.Context, v1API pb.GreeterServer) error {

	// load peer cert/key, ca cert
	serverCert, err := tls.LoadX509KeyPair(viper.GetString("server.tlscert"), viper.GetString("server.tlskey"))
	if err != nil {
		log.Printf("load server cert/key error:%v", err)
		return err
	}
	clientRootCaCert, err := ioutil.ReadFile(viper.GetString("client.cacert"))
	if err != nil {
		log.Printf("read ca cert file error:%v", err)
		return err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(clientRootCaCert)
	ta := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// register service(s)
	server := grpc.NewServer(grpc.Creds(ta))
	pb.RegisterGreeterServer(server, v1API)

	port := fmt.Sprintf("%d", viper.GetInt("server.rpc.port"))
	if viper.GetBool("server.rpc.localhost_only") {
		log.Printf("Starting gRPC Server (localhost) on port: %s", port)

		lisLocalhostV4, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
		if err != nil {
			log.Printf("server listen on port %s error:%v", port, err)
			return err
		}
		//go runGracefulGRPCServer(ctx, server, lisLocalhostV4)
		go runGRPCServer(server, lisLocalhostV4)

		lisLocalhostV6, err := net.Listen("tcp", fmt.Sprintf("[::1]:%s", port))
		if err != nil {
			log.Printf("server listen on port %s error:%v", port, err)
			return err
		}
		//go runGracefulGRPCServer(ctx, server, lisLocalhostV6)
		go runGRPCServer(server, lisLocalhostV6)

		return nil

	} else {
		log.Printf("Starting gRPC Server (all interfaces) on port: %s", port)

		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			log.Printf("server listen on port %s error:%v", port, err)
			return err
		}
		//go runGracefulGRPCServer(ctx, server, lis)
		go runGRPCServer(server, lis)

		return nil
	}
}