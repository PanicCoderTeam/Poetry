package basic

import "errors"

const (
	// Instance State List
	PendingState      = "PENDING"
	LaunchFailedState = "LAUNCH_FAILED"
	RunningState      = "RUNNING"
	StoppedState      = "STOPPED"
	StartingState     = "STARTING"
	StoppingState     = "STOPPING"
	RebootingState    = "REBOOTING"
	ShutdownState     = "SHUTDOWN"
	TerminatingState  = "TERMINATING"
	TerminatedState   = "TERMINATED"

	// Node State List
	NodeNormal  = "NORMAL"
	NodeLimited = "LIMITED"
	NodeOffline = "OFFLINE"
	NodeSellout = "SELLOUT"

	// ISP List
	CTCC     = "电信"
	CUCC     = "联通"
	CMCC     = "移动"
	ThreeNet = "三网"

	// ISP ID List
	CTCCId = 0
	CUCCId = 1
	CMCCId = 2

	// ISP Name List
	CTCCName    = "CTCC"
	CUCCName    = "CUCC"
	CMCCName    = "CMCC"
	ThreeNetISP = "CTCC;CUCC;CMCC"

	// source 请求来源
	SourceAPI = "API"
	SourceMC  = "MC"

	ISPSplitter = ";"
	IPSplitter  = ";"

	UUIDLen           = 36
	InstanceIdLen     = 12
	AuthFilterName    = "auth_filter"
	ClaimsKeyVal      = ClaimsKey("user_claims")
	CvmApiVersion     = `2017-03-12`
)

type ClaimsKey string

var (
	AllIsp     = []string{CTCC, CUCC, CMCC}
	AllIspId   = []int64{CTCCId, CUCCId, CMCCId}
	AllIspName = []string{CTCCName, CUCCName, CMCCName}

	IrregularStates = []string{PendingState, LaunchFailedState, TerminatedState}
	InvalidStates   = []string{PendingState, LaunchFailedState}
)

func IsNodeState(s string) bool {
	if s != NodeNormal && s != NodeLimited && s != NodeSellout && s != NodeOffline {
		return false
	}
	return true
}

var (
	// ErrReturnEmpty 表示应该返回空结果的错误
	ErrReturnEmpty = errors.New("should return empty")
)

// 10745 | 西安电信咸新区OC1-240G-MU
// 10911 | 重庆移动水土OC7-160G-MY
// 2707  | 深圳电信沙河OC7-30G-V
// 2905  | 太原电信数码西路OC2-30G-V
// 2941  | 长沙联通荷花园OC2-30G-V
// 3172  | 天津联通空港OC7-30G-V
// 3208  | 上海移动松江机房OC5-30G-V
// 3703  | 太原移动开发区OC6-30G-V
// 3928  | 深圳移动宝观OC2-30G-V
// 4828  | 深圳移动宝观OC3-30G-V
// 5326  | 宁波电信镇海OC2-30G-V
// 5527  | 北京电信酒仙桥路OC13-30G-V
// 5673  | 成都电信西区OC15-30G-V
// 6048  | 北京移动大白楼OC3-30G-V
// 6051  | 上海联通金桥OC5-30G-V
// 6108  | 青岛电信高新区OC5-160G-MU
// 6612  | 南昌电信红谷滩OC9-30G-V
// 6726  | 南昌联通高新区OC5-30G-V
// 7043  | 武汉电信火凤凰OC1-160G-MY
// 7405  | 合肥联通滨湖新区OC2-30G-V
// 7417  | 成都联通物理二路OC4-30G-V
// 7605  | 东莞联通华南数据中心OC3-30G-V
// 7645  | 广州联通龙荣路OC2-30G-V
// 7875  | 北京联通亦庄OC9-30G-V
// 7879  | 北京联通亦庄OC11-30G-V
// 8023  | 广州电信东涌OC1-30G-V
// 8025  | 广州电信沙溪OC2-30G-V
// 8299  | 南京电信河西OC4-160G-MU
// 8703  | 北京移动信息港OC4-30G-V
// 8727  | 成都移动西云OC8-30G-V
// 8729  | 成都移动西云OC9-30G-V
// 9125  | 太原联通二长OC13-30G-V
// 9407  | 济南移动孙村OC9-160G-M
// 9425  | 成都电信西区OC25-30G-V
// 9435  | 郑州联通高新区OC27-160G-MY
// 9893  | 石家庄联通双烽OC1-160G-MU
var OldArchitectureIdcIdList = []int{
	10745, 10911, 2707, 2905, 2941, 3172, 3208, 3703, 3928, 4828, 5326, 5527, 5673, 6048, 6051, 6108, 6612, 6726, 7043,
	7405, 7417, 7605, 7645, 7875, 7879, 8023, 8025, 8299, 8703, 8727, 8729, 9125, 9407, 9425, 9435, 9893,
}

// IsOldArchitectureIdc 判断是否是老架构的节点，老架构机房不再新增，因此这里直接记下 idcId
func IsOldArchitectureIdc(idcId int) bool {
	for _, v := range OldArchitectureIdcIdList {
		if v == idcId {
			return true
		}
	}
	return false
}
