package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func ReadToString(reader io.Reader, dst *string) error {
	buf := new(strings.Builder)

	_, err := io.Copy(buf, reader)
	if err != nil {
		return err
	}

	*dst = buf.String()

	return nil
}

func TestDevSetup(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	instance := NewDockerCompose()

	type CommandResponse struct {
		msg  string
		code int
	}

	type Expected struct {
		userMsg string
		nsMsg   string
	}

	cases := []struct {
		description string
		test        func(compose *DockerCompose) (*Expected, error)
		expected    Expected
	}{
		{
			description: "succeeds",
			test: func(dockerCompose *DockerCompose) (*Expected, error) {
				cli := dockerCompose.GetServiceCLI()
				networks, err := cli.Networks(context.TODO())
				if err != nil {
					return nil, err
				}
				logrus.Info(networks)

				// Try to create a new user
				_, reader, err := cli.Exec(ctx, strings.Split("./cli user create john_doe secret john.doe@test.com", " "))
				if err != nil {
					return nil, err
				}

				var userRes string
				if err := ReadToString(reader, &userRes); err != nil {
					return nil, err
				}
				logrus.Info(userRes)

				// Try to create a new namespace
				_, reader, err = cli.Exec(ctx, strings.Split("./cli namespace create dev john_doe 00000000-0000-4000-0000-000000000000", " "))
				if err != nil {
					return nil, err
				}

				var nsRes string
				if err := ReadToString(reader, &nsRes); err != nil {
					return nil, err
				}
				logrus.Info(nsRes)

				userAuth := new(models.UserAuthResponse)
				res1, err := resty.
					New().
					R().
					SetBody(map[string]string{
						"username": "john_doe",
						"password": "secret",
					}).
					SetResult(userAuth).
					Post(fmt.Sprintf("http://localhost:%s/api/login", dockerCompose.GetEnv("SHELLHUB_HTTP_PORT")))
				if err != nil {
					return nil, err
				}
				logrus.Info(res1.Status())
				logrus.Info(string(res1.Body()))

				time.Sleep(60 * time.Second)

				devicesList := make([]models.Device, 1)
				res2, err := resty.
					New().
					R().
					SetHeader("authorization", fmt.Sprintf("Bearer %s", userAuth.Token)).
					SetResult(&devicesList).
					Get(fmt.Sprintf("http://localhost:%s/api/devices", dockerCompose.GetEnv("SHELLHUB_HTTP_PORT")))
				if err != nil {
					return nil, err
				}
				for _, device := range devicesList {
					logrus.Infof("%+v", device)
				}
				logrus.Info(res2.Status())
				logrus.Info(string(res2.Body()))

				_, err = resty.
					New().
					R().
					SetHeader("authorization", fmt.Sprintf("Bearer %s", userAuth.Token)).
					Patch(fmt.Sprintf("http://localhost:%s/api/devices/%s/accept", dockerCompose.GetEnv("SHELLHUB_HTTP_PORT"), devicesList[0].UID))
				if err != nil {
					return nil, err
				}

				time.Sleep(60 * time.Second)

				devicesList = make([]models.Device, 1)
				_, err = resty.
					New().
					R().
					SetHeader("authorization", fmt.Sprintf("Bearer %s", userAuth.Token)).
					SetResult(&devicesList).
					Get(fmt.Sprintf("http://localhost:%s/api/devices", dockerCompose.GetEnv("SHELLHUB_HTTP_PORT")))
				if err != nil {
					return nil, err
				}
				for _, device := range devicesList {
					logrus.Infof("%+v", device)
				}

				return &Expected{
					userMsg: userRes,
					nsMsg:   nsRes,
				}, nil
			},
			expected: Expected{
				// TODO: how can we assert the exit's code?
				userMsg: "\nUsername: john_doe\nEmail: john.doe@test.com\n",
				nsMsg:   "Namespace created successfully\nNamespace: dev\nTenant: 00000000-0000-4000-0000-000000000000\nOwner:", // TODO: how can we assert the Owner ID?
			},
		},
	}

	for i, tt := range cases {
		tc := tt

		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			dockerCompose := instance.Clone().WithEnv("SHELLHUB_NETWORK", fmt.Sprintf("shellhub_network_%d", i+1))

			teardown, err := dockerCompose.Start()
			if !assert.NoError(t, err) {
				t.Fatal(err)
			}
			defer teardown() // nolint: errcheck

			actual, err := tc.test(dockerCompose)
			if assert.NoError(t, err) {
				assert.Contains(t, actual.userMsg, tc.expected.userMsg)
				assert.Contains(t, actual.nsMsg, tc.expected.nsMsg)
			}
		})
	}
}
