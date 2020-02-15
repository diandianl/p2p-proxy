/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/diandianl/p2p-proxy/endpoint"
	"github.com/diandianl/p2p-proxy/p2p"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newEndpointCmd(ctx context.Context) *cobra.Command {

	// endpointCmd represents the endpoint command
	var endpointCmd = &cobra.Command{
		Use:   "endpoint",
		Short: "endpoint command run at local for proxy agent",
		Long: "endpoint command run at local for proxy agent",
		RunE: func(cmd *cobra.Command, args []string) error {

			checkCfgs := []string{
				"Identity.PrivKey",
				"Endpoint.Proxy",
			}

			for _, k := range checkCfgs {
				if !viper.IsSet(k) {
					return fmt.Errorf("Config '%s' is required", k)
				}
			}

			viper.SetDefault("P2P.Addr", []string{"/ip4/127.0.0.1/tcp/8888"})
			viper.SetDefault("Endpoint.HTTP", "127.0.0.1:8010")

			priv := viper.GetString("Identity.PrivKey")
			privKey, err := base64.StdEncoding.DecodeString(priv)
			if err != nil {
				return err
			}

			opts := []endpoint.Option {
				endpoint.AddP2POption(p2p.Identity(privKey)),
				endpoint.AddP2POption(p2p.Addrs(viper.GetStringSlice("P2P.Addr")...)),
				endpoint.Listen(viper.GetString("Endpoint.HTTP")),
				endpoint.Proxy(viper.GetString("Endpoint.Proxy")),
			}

			if viper.GetBool("P2P.BandwidthReporter.Enable") {
				period := 10 * time.Second
				if viper.IsSet("P2P.BandwidthReporter.Period") {
					period = viper.GetDuration("P2P.BandwidthReporter.Period")
				}
				opts = append(opts, endpoint.AddP2POption(p2p.BandwidthReporter(ctx, period)))
			}

			ep, err := endpoint.New(opts...)

			if err != nil {
				return err
			}

			return ep.Start(ctx)
		},
	}

	endpointCmd.Flags().StringP("proxy", "p", "", "proxy server address")
	viper.BindPFlag("Endpoint.Proxy", endpointCmd.Flags().Lookup("proxy"))

	endpointCmd.Flags().String("http", "", "local http(s) proxy agent listen address")
	viper.BindPFlag("Endpoint.HTTP", endpointCmd.Flags().Lookup("http"))

	return endpointCmd
}

