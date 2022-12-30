package db

import "time"

type ZCMUSER struct {
	ITG_USER_ID             string    `gorm:"size:20"`
	X_webfarm_code          string    `gorm:"size:4"`
	X_user_alias            string    `gorm:"size:20"`
	Itg_user_nm             string    `gorm:"size:100"`
	Itg_org_cm              string    `gorm:"size:10"`
	Itg_org_nm              string    `gorm:"size:100"`
	User_rbof_cd            string    `gorm:"size:4"`
	User_rbof_nm            string    `gorm:"size:100"`
	Itg_eml_adr             string    `gorm:"size:200"`
	User_co_cd              string    `gorm:"size:4"`
	User_g_scn_cd           string    `gorm:"type:char(1)"`
	User_dtl_scn_cd         string    `gorm:"type:char(1)"`
	User_co_nm              string    `gorm:"size:100"`
	UserTnSbc               string    `gorm:"size:100"`
	UserFaxTnSbc            string    `gorm:"size:100"`
	UserHpTnSbc             string    `gorm:"size:100"`
	UserBeepNoSbc           string    `gorm:"size:50"`
	UserDtpcCd              string    `gorm:"size:4"`
	UserDtpcNm              string    `gorm:"size:100"`
	UserPoaCd               string    `gorm:"size:4"`
	UserPoaNm               string    `gorm:"size:50"`
	PoaLCd                  string    `gorm:"size:4"`
	X_AddDeptInfo           string    `gorm:"type:char(1)"`
	UserChrgAffrSbc         string    `gorm:"size:100"`
	DtlChrgAffrSbc          string    `gorm:"size:100"`
	Pw                      string    `gorm:"size:100"`
	UserEnNm                string    `gorm:"size:100"`
	ItgYn                   string    `gorm:"type:char(1)"`
	SortCd                  string    `gorm:"size:100"`
	DelYn                   string    `gorm:"size:100"`
	X_UserClass             string    `gorm:"type:char(1)"`
	CreDtm                  time.Time `gorm:"size:7"`
	MdfyDtm                 time.Time `gorm:"size:7"`
	PwAltrDtm               time.Time `gorm:"size:7"`
	LgiSucsDtm              time.Time `gorm:"size:7"`
	LgiFailOft              string    `gorm:"type:decimal(6)"`
	X_Personalviewcheck     string    `gorm:"type:char(1)"`
	OtpUseYn                string    `gorm:"type:char(1)"`
	LgiYn                   string    `gorm:"type:char(1)"`
	X_PolicyIgnorelogonflag string    `gorm:"type:char(1)"`
	LangCd                  string    `gorm:"size:3"`
	X_IdcCode               string    `gorm:"size:4"`
	ItgOrgEnNm              string    `gorm:"size:100"`
	UserRofEnNm             string    `gorm:"size:100"`
	UserPoaEnNm             string    `gorm:"size:100"`
	UserDtpcEnNm            string    `gorm:"size:100"`
	UserCoENNm              string    `gorm:"size:100"`
	X_exOption              string    `gorm:"type:char(1)"`
	IfDtm                   time.Time `gorm:"size:7"`
	UserScnCd               string    `gorm:"type:char(1)"`
	SrcSysCd                string    `gorm:"size:2"`
	UseYn                   string    `gorm:"type:char(1)"`
	UserOtpUseYn            string    `gorm:"type:char(1)"`
	VbgCreOrgCd             string    `gorm:"size:10"`
	VbgCreUserId            string    `gorm:"size:20"`
	VbgCreDtm               time.Time `gorm:"size:7"`
	FinMdfyOrgCd            string    `gorm:"size:10"`
	FinMdfyUserId           string    `gorm:"size:20"`
	FinMdfyDtm              time.Time `gorm:"size:7"`
	MobileNat               string    `gorm:"size:10"`
	PartnerYn               string    `gorm:"size:10"`
	AuthCls                 string    `gorm:"size:5"`
	V_DrmChkYn              string    `gorm:"size:1"`
	ItgOrgCdPre             string    `gorm:"size:10"`
	IsNondrm                string    `gorm:"type:char(1)"`
}

type SecurityPhrase struct { // 보안문구
	Index   int    `gorm:"primaryKey;autoIncrement:true"`
	Status  string // 문구 상황
	Country string // 문구 언어
	Phrases string // 문구
}

type WebLogin struct { // 암호화로그인
	SsoId         string `gorm:"primaryKey"`
	KeyValue      string
	IP            string
	LastLoginTime time.Time
	KeyStoreTime  time.Time
	CreateAt      time.Time `gorm:"autoCreateTime:nano"`
	IsOnline      bool
}
