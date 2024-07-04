package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// App route

	rt.router.GET("/", rt.wrap(rt.homeHandle))

	// Login routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User routes
	rt.router.GET("/users", rt.wrap(rt.searchUser))
	rt.router.GET("/users/:userId/profile", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/users/:userId/profile/edit", rt.wrap(rt.setMyUserName))
	rt.router.GET("/users/:userId/stream", rt.wrap(rt.getMyStream))

	// Photos routes
	rt.router.POST("/users/:userId/photos", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/users/:userId/photos/:photosId", rt.wrap(rt.deletePhoto))
	rt.router.GET("/users/:userId/photos/:photosId", rt.wrap(rt.getPhoto))

	// Likes routes
	rt.router.POST("/users/:userId/photos/:photosId/likes", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:userId/photos/:photosId/likes/:likesId", rt.wrap(rt.unlikePhoto))

	// Comments routes
	rt.router.POST("/users/:userId/photos/:photosId/comments", rt.wrap(rt.commentPhoto))
	rt.router.GET("/users/:userId/photos/:photosId/comments", rt.wrap(rt.getPhotoComments))
	rt.router.DELETE("/users/:userId/photos/:photosId/comments/:commentsId", rt.wrap(rt.uncommentPhoto))

	// Follows routes
	rt.router.GET("/users/:userId/follows", rt.wrap(rt.getUserFollows))
	rt.router.POST("/users/:userId/follows/:followedId", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:userId/follows/:followedId", rt.wrap(rt.unfollowUser))
	rt.router.GET("/users/:userId/followers", rt.wrap(rt.getUserFollowers))
	rt.router.GET("/users/:userId/follows/:followedId", rt.wrap(rt.getIsFollowed))

	// Ban routes
	rt.router.GET("/users/:userId/bans", rt.wrap(rt.getUserBans))
	rt.router.POST("/users/:userId/bans/:bannedId", rt.wrap(rt.banUser))
	rt.router.DELETE("/users/:userId/bans/:bannedId", rt.wrap(rt.unbanUser))
	rt.router.GET("/users/:userId/bans/:bannedId", rt.wrap(rt.getIsBanned))

	return rt.router
}
