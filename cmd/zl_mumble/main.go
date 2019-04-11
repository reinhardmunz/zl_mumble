package main

import (
	"crypto/rand"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/dchote/gumble/gumble"
	_ "github.com/dchote/gumble/opus"
	"github.com/reinhardmunz/zl_mumble"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Command line flags
	server := flag.String("server", "localhost:64738", "the server to connect to")
	username := flag.String("username", "control", "the username of the client")
	password := flag.String("password", "", "the password of the server")
	insecure := flag.Bool("insecure", true, "skip server certificate verification")
	certificate := flag.String("certificate", os.Getenv("GOPATH") + "/cert/2019-04-11-mumble-cert-control-nopw.pem", "PEM encoded certificate and private key")
	channel := flag.String("channel", "Radio", "mumble channel to join by default")
	pttStart := flag.String("pttStart", os.Getenv("GOPATH") + "/audio/pttStart.wav", "ptt start WAV sound file")
	pttStop := flag.String("pttStop", os.Getenv("GOPATH") + "/audio/pttStop.wav", "ptt stop WAV sound file")

	flag.Parse()

	// Initialize
	b := zl_mumble.Talkiepi{
		Config:       gumble.NewConfig(),
		Address:      *server,
		ChannelName:  *channel,
		PttStartFile: *pttStart,
		PttStopFile:  *pttStop,
	}

	// if no username specified, lets just auto-generate a random one
	if len(*username) == 0 {
		buf := make([]byte, 6)
		_, err := rand.Read(buf)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		buf[0] |= 2
		b.Config.Username = fmt.Sprintf("talkiepi-%02x%02x%02x%02x%02x%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	} else {
		b.Config.Username = *username
	}

	b.Config.Password = *password

	if *insecure {
		b.TLSConfig.InsecureSkipVerify = true
	}
	if *certificate != "" {
		cert, err := tls.LoadX509KeyPair(*certificate, *certificate)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		b.TLSConfig.Certificates = append(b.TLSConfig.Certificates, cert)
	}

	b.Init()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	exitStatus := 0

	<-signals
	b.CleanUp()

	os.Exit(exitStatus)
}
