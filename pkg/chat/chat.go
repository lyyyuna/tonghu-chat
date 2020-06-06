package chat

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type API struct {
	store *ChatStore
}

func NewChatServer(r *gin.Engine, store *ChatStore) {
	api := &API{
		store: store,
	}
	channels := r.Group("/channels")
	{
		channels.POST("/register", api.register)
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

func (api *API) register(ctx *gin.Context) {
	var req registerReq
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.String(http.StatusBadRequest, "bad register request")
		return
	}

	req.Channel
}

func (api *API) listMembers(ctx *gin.Context) {

}
