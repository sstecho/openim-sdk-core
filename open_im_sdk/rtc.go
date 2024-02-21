package open_im_sdk

import (
	"github.com/openimsdk/openim-sdk-core/v3/open_im_sdk_callback"
)

func SignalingInvite(callback open_im_sdk_callback.Base, operationID string, inviteReq, offlinePushInfo string) {
	call(callback, operationID, UserForSDK.Rtc().Invite, inviteReq, offlinePushInfo)
}

func SignalingAccept(callback open_im_sdk_callback.Base, operationID string, acceptReq, offlinePushInfo string) {
	call(callback, operationID, UserForSDK.Rtc().Accept, acceptReq, offlinePushInfo)
}

func SignalingReject(callback open_im_sdk_callback.Base, operationID string, rejectReq, offlinePushInfo string) {
	call(callback, operationID, UserForSDK.Rtc().Reject, rejectReq, offlinePushInfo)
}

func SignalingCancel(callback open_im_sdk_callback.Base, operationID string, cancelReq, offlinePushInfo string) {
	call(callback, operationID, UserForSDK.Rtc().Cancel, cancelReq, offlinePushInfo)
}
