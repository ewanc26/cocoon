package server

import (
	"github.com/haileyok/cocoon/internal/helpers"
	"github.com/labstack/echo/v4"
)

type AccountRevokeRequest struct {
	Token string `form:"token"`
}

func (s *Server) handleAccountRevoke(e echo.Context) error {
	var req AccountRevokeRequest
	if err := e.Bind(&req); err != nil {
		s.logger.Error("could not bind account revoke request", "error", err)
		return helpers.ServerError(e, nil)
	}

	repo, sess, err := s.getSessionRepoOrErr(e)
	if err != nil {
		return e.Redirect(303, "/account/signin")
	}

	if err := s.db.Exec("DELETE FROM oauth_tokens WHERE sub = ? AND token = ?", nil, repo.Repo.Did, req.Token).Error; err != nil {
		s.logger.Error("couldnt delete oauth session for account", "did", repo.Repo.Did, "token", req.Token, "error", err)
		sess.AddFlash("Unable to revoke session. See server logs for more details.", "error")
		sess.Save(e.Request(), e.Response())
		return e.Redirect(303, "/account")
	}

	sess.AddFlash("Session successfully revoked!", "success")
	sess.Save(e.Request(), e.Response())
	return e.Redirect(303, "/account")
}
