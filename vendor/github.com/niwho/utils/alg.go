package utils

func inner2find(startpos int, val int64, ids []int64) int {
	if len(ids)/2 == 0 {
		if len(ids) > 0 && val < ids[0] {
			startpos += 1
		}
		return startpos
	}
	if ids[len(ids)/2] < val {
		return inner2find(startpos, val, ids[0:len(ids)/2])
	} else {
		return inner2find(startpos+(len(ids))/2, val, ids[len(ids)/2:])
	}
}

// 在有序列表中二分查找位置
func GetPosBefore(mid int64, groupMsgIds []int64) int {
	return inner2find(0, mid, groupMsgIds)

}
