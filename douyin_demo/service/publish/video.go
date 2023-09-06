package publish

import (
	"douyin_demo/dao"
	"douyin_demo/log"
	"fmt"
	"sync"
)

type VideoList struct {
	Id            int64         `json:"id"`
	Author        dao.VideoUser `json:"author"`
	PlayUrl       string        `json:"play_url,omitempty"`
	CoverUrl      string        `json:"cover_url,omitempty"`
	FavoriteCount int64         `json:"favorite_count,omitempty"`
	CommentCount  int64         `json:"comment_count,omitempty"`
	IsFavorite    bool          `json:"is_favorite,omitempty"`
	Title         string        `json:"title,omitempty"`
}

// AssembleVideoList 组装完整视频列表信息
func AssembleVideoList(userId int) []VideoList {
	video, err := dao.NewVideDao().GetVideoList()
	if err != nil {
		return nil
	}
	var videos []VideoList
	favoriteList, err := dao.NewFavoriteDao().FindFavoritesByUserId(userId)
	if err != nil {
		return nil
	}
	favorites := make(map[int]bool)
	for _, fav := range favoriteList {
		favorites[fav.VideoId] = true
	}
	follows := make(map[int]bool)
	relations := dao.NewRelationDao().FollowById(userId)
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(relations))
	for _, relation := range relations {
		go func(relation dao.Relation) {
			defer wg.Done()
			var userToFetch int
			if relation.UserIdOne == userId {
				userToFetch = relation.UserIdTwo
			} else {
				userToFetch = relation.UserIdOne
			}
			lock.Lock()
			follows[userToFetch] = true
			lock.Unlock()
		}(relation)
	}
	wg.Wait()
	for _, videoInfo := range video {
		videoUser := dao.NewVideDao().GetVideoUser(videoInfo.Author)
		videoUser.IsFollow = follows[int(videoUser.Id)]
		videos = append(videos, VideoList{
			Id:            int64(videoInfo.ID),
			Author:        *videoUser,
			PlayUrl:       fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, videoInfo.PlayUrl),
			CoverUrl:      fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, videoInfo.CoverUrl),
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    favorites[int(videoInfo.ID)],
			Title:         videoInfo.Title,
		})
	}
	return videos
}

// AssembleVideoListById 获取视频发布者对应信息
func AssembleVideoListById(id int) []VideoList {
	video, err := dao.NewVideDao().GetVideoListById(id)
	if err != nil {
		log.SetLog().Error("[video] AssembleVideoListById")
		return nil
	}
	var videos []VideoList
	for _, videoInfo := range video {
		videoUser := dao.NewVideDao().GetVideoUser(videoInfo.Author)
		videos = append(videos, VideoList{
			Id:            int64(videoInfo.ID),
			Author:        *videoUser,
			PlayUrl:       fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, videoInfo.PlayUrl),
			CoverUrl:      fmt.Sprintf("http://%s:%s/%s", dao.HOST, dao.PORT, videoInfo.CoverUrl),
			FavoriteCount: videoInfo.FavoriteCount,
			CommentCount:  videoInfo.CommentCount,
			IsFavorite:    videoInfo.IsFavorite,
			Title:         videoInfo.Title,
		})
	}
	return videos
}
