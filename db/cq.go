package db

// PostForm 用于接收go-cqhttp消息
type PostForm struct {
	ID                        uint   `gorm:"primaryKey;autoIncrement"`
	GroupId                   int    `json:"group_id"`
	Interval                  int    `json:"interval"`
	MetaEventType             string `json:"meta_event_type"`
	Font                      int    `json:"font"`
	Message                   string `json:"message"`
	MessageId                 int    `json:"message_id"`
	MessageSeq                int    `json:"message_seq"`
	MessageType               string `json:"message_type"`
	PostType                  string `json:"post_type"`
	RawMessage                string `json:"raw_message"`
	SelfId                    int64  `json:"self_id"`
	SenderAge                 int    `json:"sender_age"`
	SenderArea                string `json:"sender_area"`
	SenderCard                string `json:"sender_card"`
	SenderLevel               string `json:"sender_level"`
	SenderNickname            string `json:"sender_nickname"`
	SenderRole                string `json:"sender_role"`
	SenderSex                 string `json:"sender_sex"`
	SenderTitle               string `json:"sender_title"`
	SenderUserId              int    `json:"sender_user_id"`
	StatusAppEnabled          bool   `json:"status_app_enabled"`
	StatusAppGood             bool   `json:"status_app_good"`
	StatusAppInitialized      bool   `json:"status_app_initialized"`
	StatusGood                bool   `json:"status_good"`
	StatusOnline              bool   `json:"status_online"`
	StatusStatPacketReceived  int    `json:"status_stat_packet_received"`
	StatusStatPacketSent      int    `json:"status_stat_packet_sent"`
	StatusStatPacketLost      int    `json:"status_stat_packet_lost"`
	StatusStatMessageReceived int    `json:"status_stat_message_received"`
	StatusStatMessageSent     int    `json:"status_stat_message_sent"`
	StatusStatLastMessageTime int    `json:"status_stat_last_message_time"`
	StatusStatDisconnectTimes int    `json:"status_stat_disconnect_times"`
	StatusStatLostTimes       int    `json:"status_stat_lost_times"`
	SubType                   string `json:"sub_type"`
	TargetId                  int64  `json:"target_id"`
	Time                      int    `json:"time"`
	DataTime                  string `json:"data_time"`
	UserId                    int    `json:"user_id"`
}

func (gb *GormDb) InsertCqPostFrom(cpf PostForm) (int64, error) {
	res := gb.DB.Create(&cpf)
	return res.RowsAffected, res.Error
}
