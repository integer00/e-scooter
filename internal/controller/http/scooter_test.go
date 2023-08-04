package http

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/integer00/e-scooter/internal/entity"
	mock_usecase "github.com/integer00/e-scooter/internal/mocks/usecase"
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

	require := require.New(t)

	usecase := mock_usecase.NewMockUseCase(ctl)

	controller := NewScooterController(usecase)

	usecase.EXPECT().UserLogin("alice").Times(1)
	usecase.EXPECT().GetEndpoints().Return([]byte{}).AnyTimes()

	mux := controller.NewMux()

	// 	cases := []struct {
	// 		name   string
	// 		req    *http.Request
	// 		rr     *httptest.ResponseRecorder
	// 		cookie http.Cookie
	// 		err    error
	// 	}{
	// 		{
	// 			name: "getting login",
	// 			req:  httptest.NewRequest(http.MethodGet, "/login", nil),
	// 			rr:   httptest.NewRecorder(),
	// 		},
	// 		{
	// 			name: "getting scooter with cookie",
	// 			req:  httptest.NewRequest(http.MethodGet, "/scooters", nil),
	// 			rr:   httptest.NewRecorder(),
	// 		},
	// 	}

	// 	for _, tCase := range cases {
	// 		t.Run(tCase.name, func(t *testing.T) {
	// 			mux.ServeHTTP(tCase.rr, tCase.req)

	// 		})
	// 	}
	// }

	//should get cookie
	req1 := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req1)
	require.Equal(http.StatusMovedPermanently, rr.Result().StatusCode)
	require.IsType([]*http.Cookie{}, rr.Result().Cookies())

	//save token for future
	cookies := rr.Result().Cookies()

	//no cookie
	req2 := httptest.NewRequest(http.MethodGet, "/scooters", nil)
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req2)
	require.Equal(http.StatusUnauthorized, rr.Result().StatusCode)

	//bad cookie
	req2.AddCookie(&http.Cookie{Name: "token", Value: "asdasdasd"})
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req2)
	// body, _ := io.ReadAll(rr.Result().Body)
	// require.Equal(t, string(body), ErrInvalidCookie.Error())
	require.Equal(http.StatusBadRequest, rr.Result().StatusCode)

	//good cookie
	req2 = httptest.NewRequest(http.MethodGet, "/scooters", nil)
	req2.AddCookie(cookies[0])
	rr = httptest.NewRecorder()
	mux.ServeHTTP(rr, req2)
	require.Equal(http.StatusOK, rr.Result().StatusCode)

}
