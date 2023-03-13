package handler

import (
	"go-red-envelope-rain/response"
	"net/http"

	"go-red-envelope-rain/rain/api/internal/logic"
	"go-red-envelope-rain/rain/api/internal/svc"
)

func openEnvelopeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewOpenEnvelopeLogic(r.Context(), svcCtx)
		resp, err := l.OpenEnvelope()
		response.Response(w, resp, err)
	}
}
