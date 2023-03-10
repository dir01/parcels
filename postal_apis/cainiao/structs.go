package cainiao

type response struct {
	Module  []module `json:"module"`
	Success bool     `json:"success"`
}

type module struct {
	MailNo        string `json:"mailNo"`
	OriginCountry string `json:"originCountry"`
	DestCountry   string `json:"destCountry"`
	Status        string `json:"status"`
	StatusDesc    string `json:"statusDesc"`
	MailNoSource  string `json:"mailNoSource"`
	DaysNumber    string `json:"daysNumber"`

	ProcessInfo struct {
		ProgressStatus    string  `json:"progressStatus"`
		ProgressRate      float64 `json:"progressRate"`
		Type              string  `json:"type"`
		ProgressPointList []struct {
			PointName string `json:"pointName"`
			Light     bool   `json:"light,omitempty"`
			Reload    bool   `json:"reload,omitempty"`
		} `json:"progressPointList"`
	} `json:"processInfo"`

	GlobalEtaInfo struct {
		EtaDesc         string `json:"etaDesc"`
		DeliveryMinTime int64  `json:"deliveryMinTime"`
		DeliveryMaxTime int64  `json:"deliveryMaxTime"`
	} `json:"globalEtaInfo"`

	LatestTrace struct {
		Time         int64  `json:"time"`
		TimeStr      string `json:"timeStr"`
		Desc         string `json:"desc"`
		StanderdDesc string `json:"standerdDesc"`
		DescTitle    string `json:"descTitle"`
		TimeZone     string `json:"timeZone"`
		ActionCode   string `json:"actionCode"`
		Group        group  `json:"group"`
	} `json:"latestTrace"`

	DetailList []detail `json:"detailList"`
}

type detail struct {
	Time         int64  `json:"time"`
	TimeStr      string `json:"timeStr"`
	Desc         string `json:"desc"`
	StanderdDesc string `json:"standerdDesc"`
	DescTitle    string `json:"descTitle"`
	TimeZone     string `json:"timeZone"`
	ActionCode   string `json:"actionCode"`
	Group        group  `json:"group,omitempty"`
}

type group struct {
	NodeCode       string `json:"nodeCode"`
	NodeDesc       string `json:"nodeDesc"`
	CurrentIconUrl string `json:"currentIconUrl"`
	HistoryIconUrl string `json:"historyIconUrl"`
}
