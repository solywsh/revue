package cq

// GetCqCodeFace qq表情
func GetCqCodeFace(faceId string) string {
	return "[CQ:face,qq=" + faceId + "]"
}

// GetCqCodeAt @某人
func GetCqCodeAt(qq, name string) string {
	// qq 为空或者为all时,@全体成员
	if name == "" {
		// 在本群时
		return "[CQ:at,qq=" + qq + "]"
	} else {
		// 不在本群时
		return "[CQ:at,qq=123,name=" + name + "]"
	}
}

// GetCqCodePoke 戳一戳
func GetCqCodePoke(qq string) string {
	return "[CQ:poke,qq=" + qq + "]"
}

// GetCqCodeMusic 分享音乐(标准)
func GetCqCodeMusic(musicType, musicId string) string {
	// musicType : qq 163 xm
	return "[CQ:music,type=" + musicType + ",id=" + musicId + "]"
}

// GetCqCodeImg 分享图片
func GetCqCodeImg(url string) string {
	return "[CQ:image,file=" + url + "]"
}
