package PushHttps

import (
	"TKGoBase/IO/Log"
	"TKGoBase/TKBaseDic"
	"TKGoBase/Tool"
	"TKMSGDevInterfaceService/PushMsgMgr"
	"TKMSGDevInterfaceService/TKDic"
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

//var M_DCQueMgr map[uint32]TKDic.DCQueueMgr
type DCQueueMgr struct {
	PDataCacheQueue *list.List
	M_QueLock       sync.Mutex
	FHuaWeiPush     uint32 //华为发送频率
	FMeiZuPush      uint32 //魅族发送频率
	FHuaWeiPushCopy uint32 //华为发送频率比较变量
	FMeiZuPushCopy  uint32 //魅族发送频率比较变量
}

type UserAllInfo struct {
	Userinfo  *TKDic.PWDInfo
	Queueinfo *DCQueueMgr
}

func InitDCList() *DCQueueMgr {
	var tt DCQueueMgr
	tt.PDataCacheQueue = list.New()
	return &tt
}

func CheckDataCacheOnTimer() {
	for _, v := range M_mapCertMgr {
		go v.Queueinfo.CheckDataCacheQueue()
	}
}

func (m *DCQueueMgr) CheckDataCacheQueue() {
	defer Tool.PanicRecover()
	ncount := m.PDataCacheQueue.Len()
	if ncount > 0 {
		m.M_QueLock.Lock()
		templist := m.PDataCacheQueue
		m.PDataCacheQueue = list.New()
		m.M_QueLock.Unlock()

		tklog.WriteInfolog("PDataCacheQueue Elements Count:%d.FhuaWei:%d.FMeiZu:%d.", templist.Len(), m.FHuaWeiPushCopy, m.FMeiZuPushCopy)
		for e := templist.Front(); e != nil; {
			stData, ok := e.Value.(*TKDic.PushDataCacheQueueSt)
			if !ok {
				tklog.WriteInfolog("Assert Type *TKDic.PushDataCacheQueueSt Fail!Please Check it!")
				next := e.Next()
				templist.Remove(e)
				e = next
				continue
			}

			switch stData.DevType {
			case TKBaseDic.AndroidCompany_HuaWei:
				for {
					if m.FHuaWeiPushCopy < TKDic.MaxHuaWeiPushF {
						go PushMsg2HuaWei(&stData.Header, stData.Data, stData.TryTimes)
						next := e.Next()
						templist.Remove(e)
						e = next
						break
					}
					time.Sleep(time.Second * 1)
				}

			case TKBaseDic.AndroidCompany_MeiZu:
				for {
					if m.FMeiZuPushCopy < TKDic.MaxMeiZuPushF {
						go PushMsg2MeiZu(&stData.Header, stData.Data, stData.TryTimes)
						next := e.Next()
						templist.Remove(e)
						e = next
						break
					}
					time.Sleep(time.Second * 1)
				}

			default:
				e = e.Next()
				tklog.WriteErrorlog("DevType is err!Maybe Prase Interface err!Please Check it Carefully!DevType is %d.", stData.DevType)
			}
		}
	}
}

func (m *DCQueueMgr) Dealwithfrequency(appid uint32) {
	m.FHuaWeiPushCopy, m.FMeiZuPushCopy = m.FHuaWeiPush, m.FMeiZuPush
	PushMsgMgr.PushFrequency(uint32(m.PDataCacheQueue.Len()), appid, m.FHuaWeiPushCopy, m.FMeiZuPushCopy)

	//if m.FHuaWeiPushCopy > TKDic.MaxHuaWeiPushF || m.FMeiZuPushCopy > TKDic.MaxMeiZuPushF {
	//	var h uint32 = m.FHuaWeiPushCopy / TKDic.MaxHuaWeiPushF
	//	var m uint32 = m.FMeiZuPushCopy / TKDic.MaxMeiZuPushF
	//	if h > m {
	//		time.Sleep(time.Second * time.Duration(h))
	//	} else {
	//		time.Sleep(time.Second * time.Duration(m))
	//	}
	//}
	atomic.SwapUint32(&m.FHuaWeiPush, 0) //3000以下正常
	atomic.SwapUint32(&m.FMeiZuPush, 0)  //500以下正常
}

func test() {
	//time.Sleep(time.Second * 5)
	//stData := TKDic.PushDataCacheQueueSt{}
	//stData.Data = []byte(`nihao`)
	//stData.DevType = TKDic.AndroidCompany_HuaWei
	//M_DCQueMgr.PDataCacheQueue.PushBack(&stData)
	//M_DCQueMgr.CheckDataCacheQueue()
}
