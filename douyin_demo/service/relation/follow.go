package relation

import "douyin_demo/dao"

// FollowList 实现查询关注列表
func FollowList(id int) ([]dao.VideoUser, error) {
	videoUserList, err := dao.NewRelationDao().FollowList(id)
	if err != nil {
		return nil, err
	}
	return videoUserList, nil
}
