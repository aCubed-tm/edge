package profile

import (
	"context"
	"errors"
	"github.com/acubed-tm/edge/helpers"
	proto "github.com/acubed-tm/edge/protofiles"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"net/http"
)

const service = "profile-service.acubed:50551"

func getProfileUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// TODO(authorization): ensure admin or sharing organization

	type resp struct {
		FirstName   string `json:"firstName"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	response, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewProfileServiceClient(conn)
		profile, err := c.GetProfile(ctx, &proto.GetProfileRequest{Uuid: uuid})
		if err != nil {
			return nil, errors.New(err.Error())
		}
		return resp{
			FirstName:   profile.FirstName,
			Name:        profile.LastName,
			Description: profile.Description,
		}, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, response)
}

func updateProfileUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// TODO(authorization): ensure admin or self

	var req struct {
		FirstName   string `json:"firstName"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewProfileServiceClient(conn)
		_, err := c.UpdateProfile(ctx, &proto.UpdateProfileRequest{
			Uuid:        uuid,
			FirstName:   req.FirstName,
			LastName:    req.Name,
			Description: req.Description,
		})
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

func createProfileUser(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	// TODO(authorization): ensure admin or self
	// TODO(validation): check if already exists

	var req struct {
		FirstName   string `json:"firstName"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewProfileServiceClient(conn)
		_, err := c.CreateProfile(ctx, &proto.CreateProfileRequest{
			Uuid:        uuid,
			FirstName:   req.FirstName,
			LastName:    req.Name,
			Description: req.Description,
		})
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
