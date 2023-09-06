package favorite

import (
	"douyin_demo/dao"
	"douyin_demo/service/publish"
	"fmt"
)

type FavoriteFlow struct {
	videoId    int
	actionType int
}

func NewFavoriteFlow(videoId int, actionType int) *FavoriteFlow {
	return &FavoriteFlow{
		videoId:    videoId,
		actionType: actionType,
	}
}

// Upvote 实现点赞或取消点赞逻辑
func (f *FavoriteFlow) Upvote(userId int) error {
	if f.actionType == 1 {
		err := dao.NewFavoriteDao().Like(f.videoId, userId)
		if err != nil {
			return err
		}
	}
	if f.actionType == 2 {
		err := dao.NewFavoriteDao().Unlike(f.videoId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

// FavoriteList 实现喜欢列表的显示
func FavoriteList(userId int) ([]publish.VideoList, error) {
	var favoriteList []dao.FavoriteList
	favoriteList, err := dao.NewFavoriteDao().FindFavoriteList(userId)
	if err != nil {
		return nil, err
	}
	var videoList []publish.VideoList
	for _, favorite := range favoriteList {
		videoUser := dao.NewVideDao().GetVideoUser(int64(favorite.UserId))
		video := dao.NewVideDao().GetVideoById(int64(favorite.VideoId))
		videoList = append(videoList, publish.VideoList{
			Id:            int64(video.ID),
			Author:        *videoUser,
			PlayUrl:       fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, video.PlayUrl),
			CoverUrl:      fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, video.CoverUrl),
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    true,
			Title:         video.Title,
		})
	}
	return videoList, nil
}
