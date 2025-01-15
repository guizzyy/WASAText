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

	// Conversation operations
	rt.router.GET("/conversation", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversation", rt.wrap(rt.startConversation))
	rt.router.POST("/conversation/group", rt.wrap(rt.createGroup))
	rt.router.GET("/conversation/:convID/", rt.wrap(rt.getConversation))
	rt.router.PUT("/conversation/:convID/name", rt.wrap(rt.setGroupName))
	rt.router.PUT("/conversation/:convID/photo", rt.wrap(rt.setGroupPhoto))
	rt.router.GET("/conversation/:convID/members", rt.wrap(rt.getMembers))
	rt.router.POST("/conversation/:convID/members", rt.wrap(rt.addToGroup))
	rt.router.DELETE("/conversation/:convID/members/:uID", rt.wrap(rt.leaveGroup))

	// Message operations
	rt.router.POST("/conversation/:convID/messages", rt.wrap(rt.sendMessage))
	rt.router.POST("/conversation/:convID/messages/:messID", rt.wrap(rt.forwardMessage))
	rt.router.DELETE("/conversation/:convID/messages/:messID", rt.wrap(rt.deleteMessage))

	// Comment (reaction) operations
	rt.router.GET("/conversation/:convID/messages/:messID/reactions", rt.wrap(rt.getComments))
	rt.router.PUT("/conversation/:convID/messages/:messID/reactions", rt.wrap(rt.commentMessage))
	rt.router.DELETE("/conversation/:convID/messages/:messID/reactions", rt.wrap(rt.uncommentMessage))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
