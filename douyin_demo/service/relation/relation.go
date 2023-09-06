package relation

import "douyin_demo/dao"

type RelationData struct {
	userId   int
	toUserId int
}

func NewRelation(userId int, toUserId int) *RelationData {
	return &RelationData{
		userId:   userId,
		toUserId: toUserId,
	}
}

// Relation 关注逻辑
func (r *RelationData) Relation() error {
	if err := dao.NewRelationDao().CreateRelation(r.userId, r.toUserId); err != nil {
		return err
	}
	return nil
}

// CancelRelation 取消关注逻辑
func (r *RelationData) CancelRelation() error {
	if err := dao.NewRelationDao().CancelRelation(r.userId, r.toUserId); err != nil {
		return err
	}
	return nil
}
