package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/acubed-tm/edge/helpers"
	proto "github.com/acubed-tm/edge/protofiles"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
)

const service = "authentication-service.acubed:50551"

func register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		_, err := c.Register(ctx, &proto.RegisterRequest{Email: req.Email, Password: req.Password})
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not log in: %v", err))
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccess(w, r)
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	type reply struct {
		Token string `json:"token"`
	}

	resp, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		resp, err := c.Login(ctx, &proto.LoginRequest{Email: req.Email, Password: req.Password})
		if err != nil {
			return nil, errors.New(fmt.Sprintf("could not log in: %v", err))
		}
		return reply{Token: resp.Token}, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, resp)
}

func getUserUuidAndInvites(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	type resp struct {
		Uuid    string   `json:"uuid"`
		Invites []string `json:"invites"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	accountUuid, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		resp, err := c.IsEmailRegistered(ctx, &proto.IsEmailRegisteredRequest{Email: req.Email})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		if resp.IsRegistered {
			return resp.AccountUuid, nil
		} else {
			return nil, nil
		}
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	if accountUuid == nil || accountUuid == "" {
		// could return empty response too
		helpers.WriteErrorJson(w, r, errors.New("email not found"))
		return
	}

	inviteUuids, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		resp, err := c.GetInvites(ctx, &proto.GetInvitesRequest{AccountUuid: accountUuid.(string)})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return resp.OrganizationUuids, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, resp{
		Uuid:    accountUuid.(string),
		Invites: inviteUuids.([]string),
	})
}

func verifyEmail(w http.ResponseWriter, r *http.Request) {
	emailVerificationToken := chi.URLParam(r, "token")

	_, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.ActivateEmail(ctx, &proto.ActivateEmailRequest{Token: emailVerificationToken})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	http.Redirect(w, r, "https://portal.acubed.app", 301)
}

func dropCurrentToken(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GetJwtToken(r)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.DropSingleToken(ctx, &proto.DropSingleTokenRequest{Token: token})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccess(w, r)
}

func dropAllTokens(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GetJwtToken(r)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.DropAllTokens(ctx, &proto.DropAllTokensRequest{Token: token})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccess(w, r)
}
