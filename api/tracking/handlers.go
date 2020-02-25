package tracking

import (
	"context"
	"github.com/acubed-tm/edge/helpers"
	proto "github.com/acubed-tm/edge/protofiles"
	"google.golang.org/grpc"
	"net/http"
)

const service = "tracking-service.acubed:50551"

func addCapture(w http.ResponseWriter, r *http.Request) {
	// this struct may change
	var req struct {
		CaptureX   float32 `json:"captureX"`
		CaptureY   float32 `json:"captureY"`
		Time       string  `json:"time"`
		ObjectUuid string  `json:"objectUuid"`
		CameraUuid string  `json:"cameraUuid"`
	}

	err := helpers.GetJsonFromRequestBody(r, &req)
	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	_, err = helpers.RunGrpc(service, func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error) {
		c := proto.NewTrackingServiceClient(conn)
		return c.AddCapture(ctx, &proto.AddCaptureRequest{
			CaptureX:   req.CaptureX,
			CaptureY:   req.CaptureY,
			Time:       req.Time,
			ObjectUuid: req.ObjectUuid,
			CameraUuid: req.CameraUuid,
		})
	})

	if err != nil {
		helpers.WriteErrorJson(w, r, err)
		return
	}

	helpers.WriteSuccess(w, r)
}
