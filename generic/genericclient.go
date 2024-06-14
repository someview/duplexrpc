package generic

import (
	"context"
	"fmt"
	netpoll "gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gonetpoll.git"
	"gitlab.dev.wiqun.com/tl/goserver/chat/l2/tl.gorpc.git/protocol"
)

func (c *xClient) GenericCall(ctx context.Context, method string, args any, callopts ...CallOpt) (res any, err error) {
	if c.isShutdown.Load() {
		return nil, ErrXClientShutdown
	}

	callOpt := ApplyCallOption(callopts...)
	defer callOpt.Recycle()
	if c.info.OneWay() {
		callOpt.oneway = true
	}
	end, err := c.selectClient(ctx, callOpt)
	if err != nil {
		return nil, err
	}
	methodInfo := c.GenericInfo.GetMethod(method)
	if methodInfo == nil {
		return nil, fmt.Errorf("method %s not registered", method)
	}
	req := protocol.NewMessageWithCtx(ctx)
	req.SetReq(args)
	req.SetMsgType(byte(c.info.Type()))
	seqID := genSeqID()
	req.SetSeqID(seqID)

	var sendCallBack netpoll.CallBack
	sendCallBack = func(err error) {

		if callOpt.oneway && err == nil {
			// onewayCall
			cb(nil, nil)
		} else if !callOpt.oneway && err == nil {
			// set callback, to wait server response
			c.cbMap.Set(seqID, cb)
		} else {
			// has error
			if callOpt.failedCount > c.opt.failureLimit {
				cb(nil, fmt.Errorf("failure limit :%w", err))
			} else {
				// record failedCount
				callOpt.failedCount++
				switch callOpt.failMode {
				// you must return after any retry!!!! if you don't, req.Recycle will cause panic
				case Failover:
					end, err = c.selectFailoverClient(end.Address())
					if err != nil {
						cb(nil, err)
					} else {
						end.AsyncSend(req, sendCallBack)
						return
					}
				case Failfast:
					cb(nil, err)
				case Failtry:
					end.AsyncSend(req, sendCallBack)
					return
				case Failbackup:
					cb(nil, fmt.Errorf("failBackup is not implemented"))
				}
			}

		}
		callOpt.Recycle()
		req.Recycle()
	}

	end.AsyncSend(req, sendCallBack)

}
