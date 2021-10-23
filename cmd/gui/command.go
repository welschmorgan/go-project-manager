package gui

import (
	"context"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/welschmorgan/go-release-manager/api"
	"github.com/welschmorgan/go-release-manager/config"
	"github.com/welschmorgan/go-release-manager/log"
)

var Command = &cobra.Command{
	Use:   "gui",
	Short: "Interface to show workspace",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		listenAddr := config.Get().API.ListenAddr
		v := NewWebAppView(WebAppViewConfig{
			width:    1024,
			height:   768,
			debug:    true,
			title:    "GRLM:UI",
			homePage: "home",
			hostPort: listenAddr,
		})
		s := api.NewAPIServer(listenAddr)
		viewClosed := make(chan bool)
		apiStopped := make(chan bool)
		allTasksDone := make(chan bool)

		go func() {
			// Run api server
			wg := sync.WaitGroup{}
			wg.Add(1)
			go func() {
				defer wg.Done()
				s.Serve()
				apiStopped <- true
			}()

			// Run web app view
			wg.Add(1)
			go func() {
				defer wg.Done()
				v.Start()
				viewClosed <- true
			}()

			wg.Wait()
			close(allTasksDone)
		}()

		done := false
		for !done {
			select {
			case <-apiStopped:
				log.Infof("Api server stopped")
				done = true
			case <-viewClosed:
				log.Info("Shutting down api server gracefully")
				if err := s.Shutdown(context.TODO()); err != nil {
					panic(err) // failure/timeout shutting down the server gracefully
				}
				done = true
			case <-allTasksDone:
				log.Infof("All tasks done")
				done = true
			default:
				time.Sleep(500 * time.Millisecond)
			}
		}
		log.Infof("Bye")

		return nil
	},
}
