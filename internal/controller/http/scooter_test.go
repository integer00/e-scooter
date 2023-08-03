package http

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/integer00/e-scooter/config"
	"github.com/integer00/e-scooter/internal/entity"
	mock_usecase "github.com/integer00/e-scooter/internal/mocks/usecase"
	"github.com/integer00/e-scooter/pkg/httpserver"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestParseRequest(t *testing.T) {
	require := require.New(t)

	controller := ScoController{}
	s := strings.NewReader(`{
		"scooterid": "kappa_ride",
		"userid": "alice"
	}`)

	msg, err := controller.parseRequest(s)

	require.IsType(&entity.Message{}, msg)
	require.NotNil(msg)
	require.NoError(err)

}

func TestParseRequestError(t *testing.T) {
	cases := []struct {
		name string
		in   io.Reader
		err  error
		// err  error
		// s    string
	}{
		{
			name: "parsing_bad_input_1",
			in:   strings.NewReader(`{"scooter":"asd"}`),
			err:  ErrMalformedJsonPayload,
		},
		{
			name: "parsing_bad_input_2",
			in:   strings.NewReader(`{"scooterid":"kappa_ride"}`),
			err:  ErrMalformedJsonPayload,
		},
		{
			name: "parsing_bad_input_3",
			in:   strings.NewReader(`asd`),
			err:  ErrMalformedJsonPayload,
		},
	}
	controller := ScoController{}
	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			msg, err := controller.parseRequest(tCase.in)
			logrus.Info(msg == nil)
			logrus.Info(err)
			// assert.IsType(t, entity.Message{}, err)
			require.EqualError(t, tCase.err, err.Error())
			require.Nil(t, msg)
		})
	}

}

func TestController(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	usecase := mock_usecase.NewMockUseCase(ctl)

	controller := NewScooterController(usecase)

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")

	config := config.NewConfig()

	mux := controller.NewMux()

	handler := cors.Default().Handler(mux)

	httpServer := httpserver.New(handler, config.Host+":"+config.Port)

	logrus.Info(httpServer)

}
