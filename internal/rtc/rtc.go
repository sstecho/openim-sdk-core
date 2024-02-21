package rtc

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/protocol/rtc"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/openim-sdk-core/v3/internal/util"
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/common"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/constant"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/db/db_interface"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/syncer"
	"github.com/openimsdk/openim-sdk-core/v3/pkg/utils"
)

type Signaling struct {
	listener            func() open_im_sdk_callback.OnSignalingListener
	loginUserID         string
	db                  db_interface.DataBase
	requestRecvSyncer   *syncer.Syncer[*rtc.InvitationInfo, [2]string]
	requestAcceptSyncer *syncer.Syncer[*rtc.InvitationInfo, [2]string]
	requestRejectSyncer *syncer.Syncer[*rtc.InvitationInfo, [2]string]
	requestCancelSyncer *syncer.Syncer[*rtc.InvitationInfo, [2]string]
	loginTime           int64
	conversationCh      chan common.Cmd2Value
	isCanceled          bool
}

func NewRtc(dataBase db_interface.DataBase, loginUserID string, conversationCh chan common.Cmd2Value) *Signaling {
	s := &Signaling{db: dataBase, loginUserID: loginUserID, conversationCh: conversationCh}
	s.initSyncer()
	return s
}

func (s *Signaling) initSyncer() {
	s.requestRecvSyncer = syncer.New(func(ctx context.Context, value *rtc.InvitationInfo) error {
		log.ZDebug(ctx, "requestRecvSyncer-Insert", "server", value)
		return nil
	}, func(ctx context.Context, value *rtc.InvitationInfo) error {
		return nil
	}, func(ctx context.Context, server *rtc.InvitationInfo, local *rtc.InvitationInfo) error {
		return nil
	}, func(value *rtc.InvitationInfo) [2]string {
		return [...]string{value.InviterUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *rtc.InvitationInfo) error {
		if s.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			log.ZDebug(ctx, "requestRecvSyncer-Notice-Insert", "server", server)
			s.listener().OnReceiveNewInvitation(utils.StructToJsonString(server))
		case syncer.Delete:
		case syncer.Update:
			log.ZDebug(ctx, "requestRecvSyncer-Notice-Update", "server", server)
			s.listener().OnReceiveNewInvitation(utils.StructToJsonString(server))
		}
		return nil
	})

	s.requestAcceptSyncer = syncer.New(func(ctx context.Context, value *rtc.InvitationInfo) error {
		log.ZDebug(ctx, "requestAcceptSyncer-Insert", "server", value)
		return nil
	}, func(ctx context.Context, value *rtc.InvitationInfo) error {
		return nil
	}, func(ctx context.Context, server *rtc.InvitationInfo, local *rtc.InvitationInfo) error {
		return nil
	}, func(value *rtc.InvitationInfo) [2]string {
		return [...]string{value.InviterUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *rtc.InvitationInfo) error {
		if s.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			log.ZDebug(ctx, "requestAcceptSyncer-Notice-Insert", "server", server)
			s.listener().OnInviteeAccepted(utils.StructToJsonString(server))
		case syncer.Delete:
		case syncer.Update:
			log.ZDebug(ctx, "requestAcceptSyncer-Notice-Update", "server", server)
			s.listener().OnInviteeAccepted(utils.StructToJsonString(server))
		}
		return nil
	})

	s.requestRejectSyncer = syncer.New(func(ctx context.Context, value *rtc.InvitationInfo) error {
		log.ZDebug(ctx, "requestRejectSyncer-Insert", "server", value)
		return nil
	}, func(ctx context.Context, value *rtc.InvitationInfo) error {
		return nil
	}, func(ctx context.Context, server *rtc.InvitationInfo, local *rtc.InvitationInfo) error {
		return nil
	}, func(value *rtc.InvitationInfo) [2]string {
		return [...]string{value.InviterUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *rtc.InvitationInfo) error {
		if s.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			log.ZDebug(ctx, "requestRejectSyncer-Notice-Insert", "server", server)
			s.listener().OnInviteeRejected(utils.StructToJsonString(server))
		case syncer.Delete:
		case syncer.Update:
			log.ZDebug(ctx, "requestRejectSyncer-Notice-Update", "server", server)
			s.listener().OnInviteeRejected(utils.StructToJsonString(server))
		}
		return nil
	})

	s.requestCancelSyncer = syncer.New(func(ctx context.Context, value *rtc.InvitationInfo) error {
		log.ZDebug(ctx, "requestCancelSyncer-Insert", "server", value)
		return nil
	}, func(ctx context.Context, value *rtc.InvitationInfo) error {
		return nil
	}, func(ctx context.Context, server *rtc.InvitationInfo, local *rtc.InvitationInfo) error {
		return nil
	}, func(value *rtc.InvitationInfo) [2]string {
		return [...]string{value.InviterUserID, value.ToUserID}
	}, nil, func(ctx context.Context, state int, server, local *rtc.InvitationInfo) error {
		if s.listener == nil {
			return nil
		}
		switch state {
		case syncer.Insert:
			log.ZDebug(ctx, "requestCancelSyncer-Notice-Insert", "server", server)
			s.listener().OnInvitationCancelled(utils.StructToJsonString(server))
		case syncer.Delete:
		case syncer.Update:
			log.ZDebug(ctx, "requestCancelSyncer-Notice-Update", "server", server)
			s.listener().OnInvitationCancelled(utils.StructToJsonString(server))
		}
		return nil
	})
}

// SetListener sets the user's listener.
func (s *Signaling) SetListener(listener func() open_im_sdk_callback.OnSignalingListener) {
	s.listener = listener
}

func (s *Signaling) Invite(ctx context.Context, inviteReq *rtc.InvitationInfo, p *sdkws.OfflinePushInfo) (*rtc.SignalInviteResp, error) {
	inviteReq.InviterUserID = s.loginUserID
	inviteReq.InitiateTime = utils.GetCurrentTimestampBySecond()
	// 获取token信息，并发送消息到被邀请人，使用的是对服务器发送信息，由服务器找到对应的ws客户端并发送通知信息
	req := rtc.SignalReq{
		Payload: &rtc.SignalReq_Invite{
			Invite: &rtc.SignalInviteReq{
				Invitation:      inviteReq,
				OfflinePushInfo: p,
				UserID:          s.loginUserID,
			},
		},
	}

	resp, err := util.CallApi[rtc.SignalMessageAssembleResp](ctx, constant.RtcSignalMessage, req)
	if err != nil {
		return nil, err
	}
	return resp.SignalResp.GetInvite(), nil
}

func (s *Signaling) Accept(ctx context.Context, acceptReq *rtc.InvitationInfo, p *sdkws.OfflinePushInfo) (*rtc.SignalAcceptResp, error) {
	acceptReq.InviterUserID = s.loginUserID
	acceptReq.InitiateTime = utils.GetCurrentTimestampBySecond()
	req := &rtc.SignalReq{
		Payload: &rtc.SignalReq_Accept{
			Accept: &rtc.SignalAcceptReq{
				Invitation:      acceptReq,
				OfflinePushInfo: p,
				UserID:          s.loginUserID,
			},
		},
	}

	resp, err := util.CallApi[rtc.SignalMessageAssembleResp](ctx, constant.RtcSignalMessage, req)
	if err != nil {
		return nil, err
	}
	return resp.SignalResp.GetAccept(), nil
}

func (s *Signaling) Reject(ctx context.Context, rejectReq *rtc.InvitationInfo, p *sdkws.OfflinePushInfo) (*rtc.SignalRejectResp, error) {
	rejectReq.InviterUserID = s.loginUserID
	rejectReq.InitiateTime = utils.GetCurrentTimestampBySecond()
	req := &rtc.SignalReq{
		Payload: &rtc.SignalReq_Reject{
			Reject: &rtc.SignalRejectReq{
				Invitation:      rejectReq,
				OfflinePushInfo: p,
				UserID:          s.loginUserID,
			},
		},
	}

	resp, err := util.CallApi[rtc.SignalMessageAssembleResp](ctx, constant.RtcSignalMessage, req)
	if err != nil {
		return nil, err
	}
	return resp.SignalResp.GetReject(), nil
}

func (s *Signaling) Cancel(ctx context.Context, cancelReq *rtc.InvitationInfo, p *sdkws.OfflinePushInfo) (*rtc.SignalCancelResp, error) {
	cancelReq.InviterUserID = s.loginUserID
	cancelReq.InitiateTime = utils.GetCurrentTimestampBySecond()
	req := &rtc.SignalReq{
		Payload: &rtc.SignalReq_Cancel{
			Cancel: &rtc.SignalCancelReq{
				Invitation:      cancelReq,
				OfflinePushInfo: p,
				UserID:          s.loginUserID,
			},
		},
	}

	resp, err := util.CallApi[rtc.SignalMessageAssembleResp](ctx, constant.RtcSignalMessage, req)
	if err != nil {
		return nil, err
	}
	return resp.SignalResp.GetCancel(), nil
}

func (s *Signaling) DoNotification(ctx context.Context, msg *sdkws.MsgData) {
	go func() {
		log.ZDebug(ctx, "DoNotification-func", "msg", msg)
		if err := s.doNotification(ctx, msg); err != nil {
			log.ZError(ctx, "doNotification error", err, "msg", msg)
		}
	}()
}

func (s *Signaling) doNotification(ctx context.Context, msg *sdkws.MsgData) error {
	if s.listener == nil {
		return errors.New("s.listener == nil")
	}
	log.ZDebug(ctx, "doNotification", "msg", msg)
	switch msg.ContentType {
	case constant.SignalingInviteNotification:
		tips := rtc.InvitationInfo{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return s.SyncBothSignalingRequest(ctx, &tips)
	case constant.SignalingAcceptNotification:
		tips := rtc.InvitationInfo{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return s.SyncAcceptSignalingRequest(ctx, &tips)
	case constant.SignalingRejectNotification:
		tips := rtc.InvitationInfo{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return s.SyncRejectSignalingRequest(ctx, &tips)
	case constant.SignalingCancelNotification:
		tips := rtc.InvitationInfo{}
		if err := utils.UnmarshalNotificationElem(msg.Content, &tips); err != nil {
			return err
		}
		return s.SyncCancelSignalingRequest(ctx, &tips)
	default:
		return fmt.Errorf("type failed %d", msg.ContentType)
	}
}

func (s *Signaling) SyncBothSignalingRequest(ctx context.Context, resp *rtc.InvitationInfo) error {
	if resp.ToUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncBothSignalingRequest", "toUserID", resp.ToUserID, "loginUserID", s.loginUserID)
		return s.requestRecvSyncer.Sync(ctx, util.Batch(ServerChatRequestToLocalChatRequest, []*rtc.InvitationInfo{resp}), nil, nil)
	} else if resp.InviterUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncBothSignalingRequest", "fromUserID", resp.InviterUserID, "loginUserID", s.loginUserID)
		return nil
	}
	return nil
}

func (s *Signaling) SyncAcceptSignalingRequest(ctx context.Context, resp *rtc.InvitationInfo) error {
	if resp.ToUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "toUserID", resp.ToUserID, "loginUserID", s.loginUserID)
		return s.requestAcceptSyncer.Sync(ctx, util.Batch(ServerChatRequestToLocalChatRequest, []*rtc.InvitationInfo{resp}), nil, nil)
	} else if resp.InviterUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "fromUserID", resp.InviterUserID, "loginUserID", s.loginUserID)
		return nil
	}
	return nil
}

func (s *Signaling) SyncRejectSignalingRequest(ctx context.Context, resp *rtc.InvitationInfo) error {
	if resp.ToUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "toUserID", resp.ToUserID, "loginUserID", s.loginUserID)
		return s.requestRejectSyncer.Sync(ctx, util.Batch(ServerChatRequestToLocalChatRequest, []*rtc.InvitationInfo{resp}), nil, nil)
	} else if resp.InviterUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "fromUserID", resp.InviterUserID, "loginUserID", s.loginUserID)
		return nil
	}
	return nil
}

func (s *Signaling) SyncCancelSignalingRequest(ctx context.Context, resp *rtc.InvitationInfo) error {
	if resp.ToUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "toUserID", resp.ToUserID, "loginUserID", s.loginUserID)
		return s.requestCancelSyncer.Sync(ctx, util.Batch(ServerChatRequestToLocalChatRequest, []*rtc.InvitationInfo{resp}), nil, nil)
	} else if resp.InviterUserID == s.loginUserID {
		log.ZDebug(ctx, "SyncAcceptSignalingRequest", "fromUserID", resp.InviterUserID, "loginUserID", s.loginUserID)
		return nil
	}
	return nil
}

func ServerChatRequestToLocalChatRequest(info *rtc.InvitationInfo) *rtc.InvitationInfo {
	return info
}
