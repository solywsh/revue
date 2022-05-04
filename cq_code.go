package main

// qq表情
func cqCodeFace(faceId string) string {
	return "[CQ:face,qq=" + faceId + "]"
}

// @某人
func cqCodeAt(qq, name string) string {
	// qq 为空或者为all时,@全体成员
	if name == "" {
		// 在本群时
		return "[CQ:at,qq=" + qq + "]"
	} else {
		// 不在本群时
		return "[CQ:at,qq=123,name=" + name + "]"
	}
}

// 戳一戳
func cqCodePoke(qq string) string {
	return "[CQ:poke,qq=" + qq + "]"
}

// 分享音乐(标准)
func cqCodeMusic(musicType, musicId string) string {
	// musicType : qq 163 xm
	return "[CQ:music,type=" + musicType + ",id=" + musicId + "]"
}
