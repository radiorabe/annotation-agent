package archive

import (
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/radiorabe/annotation-agent/annotator/archive/raar/client"
	"github.com/radiorabe/annotation-agent/annotator/archive/raar/client/audio_file"
	"github.com/radiorabe/annotation-agent/annotator/archive/raar/client/broadcast"
	"github.com/radiorabe/annotation-agent/annotator/archive/raar/client/user"
	"github.com/radiorabe/annotation-agent/annotator/archive/raar/models"
)

// ClientInterface for clients
type ClientInterface interface {
	Login(username string, password string) LoggedInClientInterface
}

// LoggedInClientInterface ...
type LoggedInClientInterface interface {
	GetBroadcasts(date time.Time) []*models.Broadcast
	GetRecord(uri string) *models.Broadcast
	GetFiles(record *models.Broadcast) []*models.AudioFile
}

// BroadcastServiceInterface for RAAR API
type BroadcastServiceInterface interface {
	GetBroadcastsYearMonthDay(
		params *broadcast.GetBroadcastsYearMonthDayParams,
		authInfo runtime.ClientAuthInfoWriter,
	) (*broadcast.GetBroadcastsYearMonthDayOK, error)
	GetBroadcastsID(
		params *broadcast.GetBroadcastsIDParams,
		authInfo runtime.ClientAuthInfoWriter,
	) (*broadcast.GetBroadcastsIDOK, error)
}

// UserServiceInterface ...
type UserServiceInterface interface {
	PostLogin(params *user.PostLoginParams) (*user.PostLoginOK, error)
}

// AudioFileServiceInterface ...
type AudioFileServiceInterface interface {
	GetBroadcastsBroadcastIDAudioFiles(
		params *audio_file.GetBroadcastsBroadcastIDAudioFilesParams,
		authInfo runtime.ClientAuthInfoWriter,
	) (*audio_file.GetBroadcastsBroadcastIDAudioFilesOK, error)
}

// Client concrete type
type Client struct {
	broadcastClient BroadcastServiceInterface
	userClient      UserServiceInterface
	audioFileClient AudioFileServiceInterface

	auth runtime.ClientAuthInfoWriter
}

var raarClient ClientInterface

// GetClient returns the acrClient singleton
func GetClient() ClientInterface {
	if raarClient == nil {
		t := httptransport.New(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes)
		t.Consumers["application/vnd.api+json"] = runtime.JSONConsumer()
		t.Producers["application/vnd.api+json"] = runtime.JSONProducer()
		c := client.New(t, strfmt.Default)

		raarClient = &Client{
			broadcastClient: c.Broadcast,
			userClient:      c.User,
			audioFileClient: c.AudioFile,
		}
	}
	return raarClient
}

// Login ...
func (c *Client) Login(username string, password string) LoggedInClientInterface {
	login, err := c.userClient.PostLogin(
		user.NewPostLoginParams().WithUsername(username).WithPassword(password),
	)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to log in to RAAR")
	}
	c.auth = httptransport.BearerToken(login.GetPayload().Data.Attributes.APIToken)
	logrus.WithField("id", login.GetPayload().Data.ID).Info("Login succeeded")
	return c
}

// GetRecord ...
func (c *Client) GetRecord(uri string) *models.Broadcast {
	u, err := url.Parse(uri)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	id, err := strconv.Atoi(path.Base(u.Path))
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	result, err := c.broadcastClient.GetBroadcastsID(broadcast.NewGetBroadcastsIDParams().WithID(int64(id)), c.auth)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return result.GetPayload().Data
}

// GetFiles ...
func (c *Client) GetFiles(record *models.Broadcast) []*models.AudioFile {
	idFloat, _ := strconv.ParseFloat(record.ID, 64)
	files, err := c.audioFileClient.GetBroadcastsBroadcastIDAudioFiles(
		audio_file.NewGetBroadcastsBroadcastIDAudioFilesParams().WithBroadcastID(int64(idFloat)),
		c.auth,
	)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return files.GetPayload().Data
}

// GetBroadcasts ...
func (c *Client) GetBroadcasts(date time.Time) []*models.Broadcast {
	broadcasts, err := c.broadcastClient.GetBroadcastsYearMonthDay(
		broadcast.NewGetBroadcastsYearMonthDayParams().
			WithYear(int64(date.Year())).
			WithMonth(int64(date.Month())).
			WithDay(int64(date.Day())),
		c.auth,
	)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return broadcasts.GetPayload().Data
}
