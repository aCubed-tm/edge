package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/acubed-tm/edge/helpers"
	"github.com/acubed-tm/edge/protofiles"
	"google.golang.org/grpc"
	"net/http"
)

const service = "authenticationms.acubed:50551"

func register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.GetJsonFromPostRequest(r, &req)
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

	err := helpers.GetJsonFromPostRequest(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	type reply struct {
		Token    string `json:"token"`
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
		Uuid string `json:"uuid"`
		Invites []string `json:"invites"`
	}

	err := helpers.GetJsonFromPostRequest(r, &req)
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

	if accountUuid == nil {
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
		return resp.InviteUuids, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, resp {
		Uuid:    accountUuid.(string),
		Invites: inviteUuids.([]string),
	})
}

func dropCurrentToken(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GetJwtToken(r)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.DropSingleToken(ctx, &proto.DropSingleTokenRequest{Token: token})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return nil, nil
	})

	helpers.WriteSuccess(w, r)
}

func dropAllTokens(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GetJwtToken(r)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.DropAllTokens(ctx, &proto.DropAllTokensRequest{Token: token})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return nil, nil
	})

	helpers.WriteSuccess(w, r)
}
