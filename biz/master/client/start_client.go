package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/EquaApps/frp/dao"
	"github.com/EquaApps/frp/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/sirupsen/logrus"
)

func StartFRPCHandler(ctx context.Context, req *pb.StartFRPCRequest) (*pb.StartFRPCResponse, error) {
	logrus.Infof("master get a start client request, origin is: [%+v]", req)

	userInfo := common.GetUserInfo(ctx)
	clientID := req.GetClientId()

	if !userInfo.Valid() {
		return &pb.StartFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	if len(clientID) == 0 {
		return &pb.StartFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid client id"},
		}, nil
	}

	client, err := dao.GetClientByClientID(userInfo, clientID)
	if err != nil {
		return nil, err
	}

	client.Stopped = false

	if err = dao.UpdateClient(userInfo, client); err != nil {
		return nil, err
	}

	go func() {
		resp, err := rpc.CallClient(context.Background(), req.GetClientId(), pb.Event_EVENT_START_FRPC, req)
		if err != nil {
			logrus.WithError(err).Errorf("start client event send to client error, client id: [%s]", req.GetClientId())
		}

		if resp == nil {
			logrus.Errorf("cannot get response, client id: [%s]", req.GetClientId())
		}
	}()

	return &pb.StartFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
