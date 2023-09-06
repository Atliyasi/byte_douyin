package router

import (
	"douyin_demo/controller"
	"douyin_demo/log"
	"douyin_demo/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	//r.Static("/public/mov/", "./public/mov")
	//r.Static("/public/pic/", "./public/pic")
	douYin := r.Group("/douyin")
	{
		douYin.POST("/user/register/", controller.Register)
		douYin.POST("/user/login/", controller.Login)
		douYin.GET("/user/", middleware.AuthMiddleWare(), controller.UserInfo)
	}
	{
		douYin.GET("/feed/", controller.Feed)
	}
	{
		douYin.POST("/publish/action/", controller.PublishAction)
		douYin.GET("/publish/list/", controller.PublishList)
	}
	{
		douYin.POST("/favorite/action/", middleware.AuthMiddleWare(), controller.Upvote)
		douYin.GET("/favorite/list/", controller.FavoriteList)
	}
	{
		douYin.POST("/comment/action/", middleware.AuthMiddleWare(), controller.HandleComment)
		douYin.GET("/comment/list/", controller.GetCommentList)
	}
	{
		douYin.POST("/relation/action/", middleware.AuthMiddleWare(), controller.Relation)
		douYin.GET("/relation/follow/list/", middleware.AuthMiddleWare(), controller.Follow)
		douYin.GET("/relation/follower/list/", middleware.AuthMiddleWare(), controller.Follower)
		douYin.GET("/relation/friend/list/", middleware.AuthMiddleWare(), controller.Friend)
	}
	{
		douYin.POST("/message/action/", middleware.AuthMiddleWare(), controller.SendMessage)
		douYin.GET("/message/chat/", middleware.AuthMiddleWare(), controller.GetMessage)
	}
	PORT := ":8080"
	// 启动服务
	err := r.Run(PORT)
	if err != nil {
		log.SetLog().Error(err.Error())
	}
}
