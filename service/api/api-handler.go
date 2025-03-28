package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// User operations
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.PUT("/users/:uID/username", rt.wrap(rt.setMyUserName))
	rt.router.PUT("/users/:uID/photo", rt.wrap(rt.setMyPhoto))
	rt.router.GET("/users/:uID/search", rt.wrap(rt.searchUsers))

	// Membership operations
	rt.router.DELETE("/memberships/:convID/members/:uID", rt.wrap(rt.leaveGroup))
	rt.router.POST("/memberships/:convID", rt.wrap(rt.addToGroup))
	rt.router.GET("/group/:convID", rt.wrap(rt.getGroupInfo))

	// Message operations
	rt.router.POST("/conversations/:convID/messages/:messID", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversations/:convID/messages/:messID", rt.wrap(rt.deleteMessage))
	rt.router.POST("/conversations/:convID/messages", rt.wrap(rt.sendMessage))

	// Comment (reaction) operations
	rt.router.GET("/conversations/:convID/messages/:messID/reactions", rt.wrap(rt.getComments))
	rt.router.PUT("/conversations/:convID/messages/:messID/reactions", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversations/:convID/messages/:messID/reactions", rt.wrap(rt.uncommentMessage))

	// Conversation operations
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversations", rt.wrap(rt.startConversation))
	rt.router.POST("/group", rt.wrap(rt.createGroup))
	rt.router.PUT("/conversations/:convID/manage/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversations/:convID/manage/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/conversations/:convID/open", rt.wrap(rt.getConversation))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
