package tracking

import (
	"context"
	"github.com/acubed-tm/edge/helpers"
	proto "github.com/acubed-tm/edge/protofiles"
	"github.com/go-chi/chi"
	"google.golang.org/grpc"
	"net/http"
)

const service = "tracking-service.acubed:50551"

func addCapture(w http.ResponseWriter, r *http.Request) {
	// this struct may change
	var req []struct {
		CaptureX   float32 `json:"x"`
		CaptureY   float32 `json:"y"`
		Time       int64   `json:"time"`
		ObjectUuid string  `json:"code"`
		CameraUuid string  `json:"camera"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	for _, e := range req {
		// ensure ms epochs
		if e.Time < 1500000000000 {
			e.Time *= 1000
		}
		_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
			c := proto.NewTrackingServiceClient(conn)
			return c.AddCapture(ctx, &proto.AddCaptureRequest{
				CaptureX:   e.CaptureX,
				CaptureY:   e.CaptureY,
				Time:       e.Time,
				ObjectUuid: e.ObjectUuid,
				CameraUuid: e.CameraUuid,
			})
		})
	}

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccess(w, r)
}

type objectLocation struct {
	X    float32 `json:"x"`
	Y    float32 `json:"y"`
	Z    float32 `json:"z"`
	Time int64   `json:"time"`
}

func getAllObjects(w http.ResponseWriter, r *http.Request) {
	type objectInfo struct {
		Uuid     string         `json:"uuid"`
		Name     string         `json:"name"`
		Note     string         `json:"note"`
		Location objectLocation `json:"lastLocation"`
	}

	_, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewTrackingServiceClient(conn)
		return c.UpdatePositions(ctx, &proto.UpdatePositionsRequest{Uuid:""})
	})
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	objects, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewTrackingServiceClient(conn)
		resp, err := c.GetAllObjects(ctx, &proto.GetAllObjectsRequest{})
		if err != nil {
			return nil, err
		}
		ret := make([]objectInfo, len(resp.Objects))
		for i, e := range resp.Objects {
			ret[i] = objectInfo{
				Uuid: e.Uuid,
				Name: e.Name,
				Note: e.Note,
				Location: objectLocation{
					X:    e.LastLocation.X,
					Y:    e.LastLocation.Y,
					Z:    e.LastLocation.Z,
					Time: e.LastLocation.Time,
				},
			}
		}
		return ret, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, objects)
}

func getObject(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")

	_, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewTrackingServiceClient(conn)
		return c.UpdatePositions(ctx, &proto.UpdatePositionsRequest{Uuid:uuid})
	})
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	objects, err := helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewTrackingServiceClient(conn)
		resp, err := c.GetObject(ctx, &proto.GetObjectRequest{Uuid: uuid})
		if err != nil {
			return nil, err
		}
		ret := make([]objectLocation, len(resp.Locations))
		for i, e := range resp.Locations {
			ret[i] = objectLocation{
				X:    e.X,
				Y:    e.Y,
				Z:    e.Z,
				Time: e.Time,
			}
		}
		return ret, nil
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccessJson(w, r, objects)
}
