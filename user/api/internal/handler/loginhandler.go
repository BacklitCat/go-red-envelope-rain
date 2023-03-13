package handler

import (
	"go-red-envelope-rain/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-red-envelope-rain/user/api/internal/logic"
	"go-red-envelope-rain/user/api/internal/svc"
	"go-red-envelope-rain/user/api/internal/types"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		response.Response(w, resp, err)
	}
}
