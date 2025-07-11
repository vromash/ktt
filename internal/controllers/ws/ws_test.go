package ws

import (
	"encoding/json"
	"financing-aggregator/internal/exchange"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/suite"
	"net/http/httptest"
	"testing"
)

type webSocketTestSuite struct {
	suite.Suite
	server   *httptest.Server
	wsURL    string
	wsServer WebSocketHandler
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(webSocketTestSuite))
}

func (s *webSocketTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	s.wsServer = NewWebSocketHandler(nil)
	r.GET("/ws/applications/:id", s.wsServer.SubscribeToApplicationUpdates)

	ts := httptest.NewServer(r)
	s.server = ts
	s.wsURL = "ws" + ts.URL[4:] + "/ws/applications/test-app-id"
}

func (s *webSocketTestSuite) TearDownSuite() {
	if s.server != nil {
		s.server.Close()
	}
	if s.wsServer != nil {
		s.wsServer.CloseAll()
	}
}

func (s *webSocketTestSuite) Test_WebSocket_ConnectionAndMessage() {
	offerResponse := getTestOfferResponse()

	c, _, err := websocket.DefaultDialer.Dial(s.wsURL, map[string][]string{})
	s.NoError(err)
	defer c.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, actualBytes, err := c.ReadMessage()
			s.NoError(err)

			actual := exchange.OfferResponse{}
			err = json.Unmarshal(actualBytes, &actual)
			s.NoError(err)
			s.Equal(offerResponse, actual)
			return
		}
	}()

	s.wsServer.BroadcastNewOffer("test-app-id", offerResponse)

	<-done
}

func getTestOfferResponse() exchange.OfferResponse {
	return exchange.OfferResponse{
		MonthlyPaymentAmount: 50,
		TotalRepaymentAmount: 150,
		NumberOfPayments:     3,
		AnnualPercentageRate: 10.0,
		FirstRepaymentDate:   "2025-01-01",
	}
}
