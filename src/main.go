package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"os"
)

func main() {

	cert, err := tls.LoadX509KeyPair("../certs/server.pem", "../certs/server.key")
	if err != nil {
		log.Fatalf("server: Loadkeys: %s", err)
		os.Exit(1)
	}
	certPool := x509.NewCertPool()
	pem, err := os.ReadFile("../certs/ca.pem")
	if err != nil {
		log.Fatalf("Failed to read client certificate authority: %v", err)
	}
	if !certPool.AppendCertsFromPEM(pem) {
		log.Fatalf("Can't parse client certificate authority")
	}

	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", "0.0.0.0:2083", &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			break
		}
		defer conn.Close()
		tlsConnection, ok := conn.(*tls.Conn)
		if ok {
			log.Print("server: conn: type assert to TLS succeeded")
			err := tlsConnection.Handshake()
			if err != nil {
				log.Fatalf("server: handshake failed: %s", err)
			} else {
				log.Print("server: conn: Handshake completed")
			}
			state := tlsConnection.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(v.Subject.Organization)
				log.Print(v.Subject.OrganizationalUnit)
				log.Print(v.Subject.CommonName)
				// log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
			go handleClient(conn)
		}
		log.Println("server: conn: closed")
	}
}

func handleClient(conn net.Conn) {
	buf := make([]byte, 2048)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break

		}
		log.Printf("server: conn: echo %v\n", buf[:n])
		// n, err = conn.Write(buf[:n])
		// log.Printf("server: conn: wrote %d bytes", n)
		// if err != nil {
		// 	log.Printf("server: write: %s", err)
		// 	break
		// }
	}
}
