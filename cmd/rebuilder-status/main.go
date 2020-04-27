package main

import (
	"os"

	"go.seankhliao.com/rebuilderd-go/status"
	"go.seankhliao.com/usvc"
)

func main() {
	svc := status.NewServer(os.Args)
	svc.Svc.Log.Fatal().Err(usvc.Run(usvc.SignalContext(), svc.Svc)).Msg("exited")

}
