package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	store ChatStore
}

func NewChatServer(r *gin.Engine, store ChatStore) {
	api := &API{
		store: store,
	}
	channels := r.Group("/channel")
	{
		channels.POST("/register", api.register)
	}

	adminChannels := r.Group("/admin")
	{
		adminChannels.POST("/channel", api.createChannel)
	}

}

type registerReq struct {
	UID           string `json:"uid"`
	DisplayName   string `json:"display_name"`
	Email         string `json:"email"`
	Secret        string `json:"secret"`
	Channel       string `json:"channel"`
	ChannelSecret string `json:"channel_secret"`
}

type registerRes struct {
}

func (api *API) register(ctx *gin.Context) {
	var req registerReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "bad register request")
		zap.S().Infof("Bad register request, err: %v", err)
		return
	}

	// find the channel in store
	ch, err := api.store.GetChannel(req.Channel)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "cannot find the channel")
		zap.S().Infof("Cannot find the channel, err: %v", err)
		return
	}

	// register the uid in channel
	err = ch.Register(&User{
		UID:         req.UID,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Secret:      req.Secret,
	})

	if err != nil {
		ctx.String(http.StatusInternalServerError, "cannot register to this channel")
		zap.S().Infof("Cannot register to this channel, err: %v", err)
		return
	}

	// save this new store
	err = api.store.SaveChannel(ch)
	if err != nil {
		ch.Leave(req.UID)
		ctx.String(http.StatusInternalServerError, "cannot update the channel")
		zap.S().Infof("Cannot update the channel, err: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, &registerRes{})
}

type createReq struct {
	ChannelName string `json:"channelname"`
}

type createRes struct {
}

func (api *API) createChannel(ctx *gin.Context) {
	var req createReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "bad create channel request")
		zap.S().Infof("Bad create channel request, err: %v", err)
		return
	}

	// create new channel
	ch := NewChannel(req.ChannelName)
	if err := api.store.SaveChannel(ch); err != nil {
		ctx.String(http.StatusInternalServerError, "cannot save the new channel")
		zap.S().Infof("Cannot save the new channel, err: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, &createRes{})
}

func (api *API) listMembers(ctx *gin.Context) {

}
