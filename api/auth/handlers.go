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

	if accountUuid == nil {
		accountUuid = ""
	}

	inviteUuids, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		resp, err := c.GetInvitesByEmail(ctx, &proto.GetInvitesByEmailRequest{Email: req.Email})
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

func putEmail(w http.ResponseWriter, r *http.Request) {
	// TODO(authorization): if admin or self
	emailUuid := chi.URLParam(r, "uuid")

	var req struct {
		IsPrimary bool `json:"isPrimary"` // should always be true
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	if !req.IsPrimary {
		helpers.WriteErrorJson(w, r, errors.New("cannot make email non-primary, make another email primary instead"))
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.MakeEmailPrimary(ctx, &proto.MakeEmailPrimaryRequest{EmailUuid: emailUuid})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}
}

func addEmail(w http.ResponseWriter, r *http.Request) {
	// TODO(authorization): if admin or self
	// TODO: send verification email
	var req struct {
		UserUuid string `json:"userUuid"`
		Email    string `json:"email"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.AddEmail(ctx, &proto.AddEmailRequest{
			AccountUuid: req.UserUuid,
			Email:       req.Email,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}
}

func deleteEmail(w http.ResponseWriter, r *http.Request) {
	// TODO(authorization): if admin or self
	emailUuid := chi.URLParam(r, "uuid")

	_, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewAuthServiceClient(conn)
		_, err := c.DeleteEmail(ctx, &proto.DeleteEmailRequest{
			Uuid: emailUuid,
		})
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}
}

// TODO: should move this to helper function
func _(w http.ResponseWriter, r *http.Request) (string, error) { // getCurrentUserUuid
	token, err := helpers.GetJwtToken(r)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return "", nil
	}

	accountUuid, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		// Contact the server and print out its response.
		c := proto.NewAuthServiceClient(conn)
		resp, err := c.GetUuidFromToken(ctx, &proto.GetUuidFromTokenRequest{Token: token})
		if err != nil {
			return "", errors.New(err.Error())
		}
		return resp.Uuid, nil
	})

	if err != nil {
		return "", err
	}
	return accountUuid.(string), nil
}
