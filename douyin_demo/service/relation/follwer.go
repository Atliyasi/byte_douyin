package relation

import "douyin_demo/dao"

// FollowerList 实现查询粉丝列表
func FollowerList(id int) ([]dao.VideoUser, error) {
	videoUserList, err := dao.NewRelationDao().FollowerList(id)
	if err != nil {
		return nil, err
	}
	return videoUserList, nil
}
