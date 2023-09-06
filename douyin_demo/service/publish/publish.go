package publish

import (
	"douyin_demo/dao"
	"douyin_demo/util"
	"gorm.io/gorm"
)

// SavePublishVideo 保存视频信息
func SavePublishVideo(saveFile string, finalName string, title string, id int64) error {
	snapshotName := util.GetSnapshot(saveFile, finalName)
	video := &dao.Video{
		Model:         gorm.Model{},
		Author:        id,
		PlayUrl:       finalName,
		CoverUrl:      snapshotName,
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
		Title:         title,
	}
	if err := dao.NewVideDao().SetVideo(video); err != nil {
		return err
	}
	if err := dao.NewVideDao().UpdateVideoUserByWork(id); err != nil {
		return err
	}
	return nil
}
