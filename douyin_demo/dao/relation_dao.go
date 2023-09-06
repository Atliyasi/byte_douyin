package dao

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type Relation struct {
	gorm.Model
	UserIdOne int  `gorm:"column:user_id_one"`
	UserIdTwo int  `gorm:"column:user_id_two"`
	Forward   bool `gorm:"column:forward"` // 正向关系 UserOne->UserTwo
	Reverse   bool `gorm:"column:reverse"` // 反向关系 UserOne<-UserTwo
}

func (*Relation) TableName() string {
	return "relation"
}

type RelationDao struct{}

var relationDao *RelationDao
var RelationOnce sync.Once

func NewRelationDao() *RelationDao {
	RelationOnce.Do(func() {
		relationDao = &RelationDao{}
	})
	return relationDao
}

// CreateRelation 关注底层
func (*RelationDao) CreateRelation(userIdOne int, userIdTwo int) error {
	db := GetDB()
	tx := db.Begin()
	var relation *Relation
	if err := tx.Where("user_id_one=?", userIdOne).Where("user_id_two=?", userIdTwo).Or("user_id_one=?", userIdTwo).Where("user_id_two=?", userIdOne).First(&relation).Error; err != nil {
		fmt.Println("进入error")
		relation = &Relation{
			Model:     gorm.Model{},
			UserIdOne: userIdOne,
			UserIdTwo: userIdTwo,
			Forward:   true,
			Reverse:   false,
		}
		if err := tx.Create(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := addFollowCount(userIdOne, tx, 1); err != nil {
			return err
		}
		if err := addFollowerCount(userIdTwo, tx, 1); err != nil {
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}
		return nil
	}
	fmt.Println("relation.UserIdOne: ", relation.UserIdOne, ", userIdOne: ", userIdOne)
	if relation.UserIdOne == userIdOne {
		relation.Forward = true
		if err := tx.Save(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := addFollowCount(userIdOne, tx, 1); err != nil {
			return err
		}
		if err := addFollowerCount(userIdTwo, tx, 1); err != nil {
			return err
		}
	} else {
		relation.Reverse = true
		if err := tx.Save(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := addFollowCount(userIdOne, tx, 1); err != nil {
			return err
		}
		if err := addFollowerCount(userIdTwo, tx, 1); err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// CancelRelation 取消关注
func (*RelationDao) CancelRelation(userIdOne int, userIdTwo int) error {
	var relation Relation
	db := GetDB()
	tx := db.Begin()
	if err := tx.Where("user_id_one=?", userIdOne).Where("user_id_two=?", userIdTwo).Or("user_id_one=?", userIdTwo).Where("user_id_two=?", userIdOne).First(&relation).Error; err != nil {
		tx.Rollback()
		return err
	}
	if relation.UserIdOne == userIdOne {
		relation.Forward = false
		if err := tx.Save(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := addFollowCount(userIdOne, tx, -1); err != nil {
			return err
		}
		if err := addFollowerCount(userIdTwo, tx, -1); err != nil {
			return err
		}
	} else {
		relation.Reverse = false
		if err := tx.Save(relation).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := addFollowCount(userIdOne, tx, -1); err != nil {
			return err
		}
		if err := addFollowerCount(userIdTwo, tx, -1); err != nil {
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// addFollowCount 实现关注后关注着关注数变化
func addFollowCount(id int, tx *gorm.DB, num int64) error {
	var userOne VideoUser
	if err := tx.Where("id=?", id).First(&userOne).Error; err != nil {
		tx.Rollback()
		return err
	}
	userOne.FollowCount += num
	if err := tx.Save(&userOne).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// addFollowerCount 实现关注后被关注者粉丝数变化
func addFollowerCount(id int, tx *gorm.DB, num int64) error {
	var userTwo VideoUser
	if err := tx.Where("id=?", id).First(&userTwo).Error; err != nil {
		tx.Rollback()
		return err
	}
	userTwo.FollowerCount += num
	if err := tx.Save(&userTwo).Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// FollowList 实现关注列表查询
func (*RelationDao) FollowList(id int) ([]VideoUser, error) {
	var videoUserList []VideoUser
	var relations []Relation
	db := GetDB()
	tx := db.Begin()
	if err := tx.Where("user_id_one=?", id).Where("forward=?", true).Or("user_id_two=?", id).Where("reverse=?", true).Find(&relations).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, relation := range relations {
		if relation.UserIdOne == id {
			var videoUser VideoUser
			if err := tx.Where("id=?", relation.UserIdTwo).First(&videoUser).Error; err != nil {
				continue
			}
			videoUser.IsFollow = true
			videoUserList = append(videoUserList, videoUser)
		} else {
			var videoUser VideoUser
			if err := tx.Where("id=?", relation.UserIdOne).First(&videoUser).Error; err != nil {
				continue
			}
			videoUser.IsFollow = true
			videoUserList = append(videoUserList, videoUser)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return videoUserList, nil
}

// FollowerList 实现粉丝列表查询
func (*RelationDao) FollowerList(id int) ([]VideoUser, error) {
	var videoUserList []VideoUser
	var relations []Relation
	db := GetDB()
	tx := db.Begin()
	if err := tx.Where("user_id_one=?", id).Where("reverse=?", true).Or("user_id_two=?", id).Where("forward=?", true).Find(&relations).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	for _, relation := range relations {
		if relation.UserIdOne == id {
			var videoUser VideoUser
			if err := tx.Where("id=?", relation.UserIdTwo).First(&videoUser).Error; err != nil {
				continue
			}
			if relation.Forward == true {
				videoUser.IsFollow = true
			}
			videoUserList = append(videoUserList, videoUser)
		} else {
			var videoUser VideoUser
			if err := tx.Where("id=?", relation.UserIdOne).First(&videoUser).Error; err != nil {
				continue
			}
			if relation.Reverse == true {
				videoUser.IsFollow = true
			}
			videoUserList = append(videoUserList, videoUser)
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return videoUserList, nil
}

func (*RelationDao) FollowById(userId int) []Relation {
	var relation []Relation
	if err := GetDB().Where("user_id_one=? AND forward=? OR user_id_two=? AND reverse=?", userId, true, userId, true).Find(&relation).Error; err != nil {
		return nil
	}
	return relation
}

// FriendList 互关为朋友
func (*RelationDao) FriendList(userId int) ([]VideoUser, error) {
	db := GetDB()
	tx := db.Begin()
	var relations []Relation
	if err := tx.Where("forward = ? AND reverse = ? AND (user_id_one = ? OR user_id_two = ?)", true, true, userId, userId).Find(&relations).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	var videoUserList []VideoUser
	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(relations))
	for _, relation := range relations {
		go func(relation Relation) {
			defer wg.Done()
			var userToFetch int
			if relation.UserIdOne == userId {
				userToFetch = relation.UserIdTwo
			} else {
				userToFetch = relation.UserIdOne
			}
			var videoUser VideoUser
			if err := tx.Where("id = ?", userToFetch).First(&videoUser).Error; err != nil {
				tx.Rollback()
				return
			}
			mutex.Lock()
			videoUserList = append(videoUserList, videoUser)
			mutex.Unlock()
		}(relation)
	}
	wg.Wait()
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	return videoUserList, nil
}
