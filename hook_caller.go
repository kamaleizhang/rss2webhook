package r2w

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type HookCaller interface {
	CallHooks(config Config, request HookRequest) error
}

type HookCallerImpl struct {
	httpClient *http.Client
}

func NewHookCallerImpl(httpClient *http.Client) HookCaller {
	return &HookCallerImpl{
		httpClient: httpClient,
	}
}

func (h *HookCallerImpl) CallHooks(config Config, request HookRequest) error {
	log.Printf("Calling hook: %v", config.TargetWebhook)
	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}
	post, err := h.httpClient.Post(config.TargetWebhook, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	log.Printf("Post response: %v", post)
	return nil
}
