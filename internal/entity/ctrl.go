// Package entity defines main entities for business logic (services), data base mapping and
// HTTP response objects if suitable. Each logic group entities in own file.
package entity

import "time"

// Cluster -.
type Cluster struct {
	ID            string `json:"id"       rac:"cluster"                        example:"UUID like"`
	Host          string `json:"host"     rac:"host"                           example:"localhost"`
	Port          string `json:"port"     rac:"port"                           example:"1541"`
	Name          string `json:"name"     rac:"name"                           example:"name as text"`
	Exp           int    `json:"exp"      rac:"expiration-timeout"             example:"int"`
	LT            int    `json:"lt"       rac:"lifetime-limit"                 example:"int"`
	MaxMemSize    int    `json:"mms"      rac:"max-memory-size"                example:"int"`
	MaxMemTimeLim int    `json:"mmts"     rac:"max-memory-time-limit"          example:"int"`
	SecLevel      int    `json:"sl"       rac:"security-level"                 example:"int"`
	SesFTLevel    int    `json:"sftl"     rac:"session-fault-tolerance-level"  example:"int"`
	LBMode        string `json:"lb"       rac:"load-balancing-mode"            example:"name as text"`
	ErrCountTh    int    `json:"errth"    rac:"errors-count-threshold"         example:"int"`
	KillPP        int    `json:"kpp"      rac:"kill-problem-process"           example:"int"`
}

// Infobase -.
type Infobase struct {
	ID   string `json:"id"    rac:"infobase"   example:"UUID like"`
	Name string `json:"name"  rac:"name"       example:"name as text"`
	Desc string `json:"desc"  rac:"descr"      example:"some comments"`
}

// Session -.
type Session struct {
	ID             string    `json:"id"              rac:"session"     example:"UUID like"`
	SID            int       `json:"sid"             rac:"session-id"  example:"Int like"`
	InfobaseID     string    `json:"ib"              rac:"infobase"    example:"UUID of infobase"`
	ConnectionID   string    `json:"conn"            rac:"connection"  example:"UUID of connection"`
	ProcessID      string    `json:"proc"            rac:"process"     example:"UUID of process"`
	UserName       string    `json:"uname"           rac:"user-name"   example:"Name of the user"`
	Host           string    `json:"host"            rac:"host"        example:"Host of the user"`
	AppID          string    `json:"appid"           rac:"app-id"      example:"Application identifier"`
	Loc            string    `json:"loc"             rac:"locale"      example:"Lang of session string like"`
	Started        time.Time `json:"started"         rac:"started-at"      example:"Time of start"`
	LastActive     time.Time `json:"active"          rac:"last-active-at"  example:"Time of last activity"`
	Hibernate      string    `json:"hib"             rac:"hibernate"       example:"Hibernation yes/no"`
	HiberTime      int       `json:"hibtm"             rac:"passive-session-hibernate-time"  example:"Passive session hibernation time"`
	HiberTermTime  int       `json:"hibterm"             rac:"hibernate-session-terminate-time" example:"Termination session hibernation time"`
	BlockedDB      int       `json:"blockdb"             rac:"blocked-by-dbms"  example:"Blocked by dbms"`
	BlockedLS      int       `json:"blockls"             rac:"blocked-by-ls"  example:"Blocked by ls"`
	Bytes          int       `json:"bytes"             rac:"bytes-all"  example:"Bytes all"`
	Bytes5m        int       `json:"bytes5m"             rac:"bytes-last-5min"  example:"Bytes last 5 min"`
	Calls          int       `json:"calls"             rac:"calls-all"  example:"Calls all"`
	Calls5m        int       `json:"calls5m"             rac:"calls-last-5min"  example:"Calls last 5 min"`
	BytesDB        int       `json:"bytesdb"             rac:"dbms-bytes-all"  example:"Bytes dbms all"`
	BytesDB5m      int       `json:"bytesdb5m"             rac:"dbms-bytes-last-5min"  example:"Bytes dbms last 5 min"`
	DBProcInfo     string    `json:"dbproci"             rac:"db-proc-info"  example:"DB proc info"`
	DBProc         int       `json:"dbproc"             rac:"db-proc-took"  example:"DB proc took"`
	DBProcAt       string    `json:"dbprocat"             rac:"db-proc-took-at"  example:"DB proc took at time"`
	Duration       int       `json:"dur"             rac:"duration-all"  example:"Duration all"`
	DurationDB     int       `json:"durdb"             rac:"duration-all-dbms"  example:"Duration DB all"`
	DurationCur    int       `json:"durcur"             rac:"duration-current"  example:"Duration current"`
	DurationCurDB  int       `json:"durcurdb"             rac:"duration-current-dbms"  example:"Duration db current"`
	Duration5m     int       `json:"dur5m"             rac:"duration-last-5min"  example:"Duration last 5 min"`
	DurationDB5m   int       `json:"durdb5m"             rac:"duration-last-5min-dbms"  example:"Duration db last 5 min"`
	MemoryCur      int       `json:"memcur"             rac:"memory-current"  example:"Memory current"`
	Memory5m       int       `json:"mem5m"             rac:"memory-last-5min"  example:"Memory last 5 min"`
	Memory         int       `json:"mem"             rac:"memory-total"  example:"Memory all"`
	ReadCur        int       `json:"readcur"             rac:"read-current"  example:"Read current"`
	Read5m         int       `json:"read5m"             rac:"read-last-5min"  example:"Read last 5 min"`
	Read           int       `json:"read"             rac:"read-total"  example:"Read all"`
	WriteCur       int       `json:"writecur"             rac:"write-current"  example:"Write current"`
	Write5m        int       `json:"write5m"             rac:"write-last-5min"  example:"Write last 5 min"`
	Write          int       `json:"write"             rac:"write-total"  example:"Write all"`
	DurationSvcCur int       `json:"dursvccur"             rac:"duration-current-service"  example:"Duration service current"`
	DurationSvc5m  int       `json:"dursvc5m"             rac:"duration-last-5min-service"  example:"Duration service last 5 min"`
	DurationSvc    int       `json:"dursvc"             rac:"duration-all-service"  example:"Duration service all"`
	Svc            string    `json:"svc"             rac:"current-service-name"  example:"Current servie name"`
	CPUCur         int       `json:"cpucur"             rac:"cpu-time-current"  example:"CPU current time"`
	CPU5m          int       `json:"cpu5m"             rac:"cpu-time-last-5min"  example:"CPU last 5 min"`
	CPU            int       `json:"cpu"             rac:"cpu-time-total"  example:"CPU all"`
	Sep            string    `json:"sep"             rac:"data-separation"  example:"Data separation"`
}

// Connection -.
type Connection struct {
	ID         string    `json:"id"          rac:"connection" example:"UUID like"`
	CID        int       `json:"cid"         rac:"connection-id"  example:"Int like"`
	InfobaseID string    `json:"ib"          rac:"infobase" example:"UUID of infobase"`
	ProcessID  string    `json:"proc"        rac:"process" example:"UUID of process"`
	Host       string    `json:"host"        rac:"host" example:"Host of the user"`
	AppID      string    `json:"appid"       rac:"application" example:"Application identifier"`
	Connected  time.Time `json:"connected"   rac:"connected-at"      example:"Time of start"`
	SID        int       `json:"sid"         rac:"session-number"  example:"Int like session number"`
	Blocked    int       `json:"blocked"     rac:"blocked-by-ls"  example:"Int like blocked by ls"`
}
