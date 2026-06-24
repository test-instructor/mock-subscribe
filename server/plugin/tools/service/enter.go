package service

var Service = new(service)

type service struct {
	Environment  environment
	UserRelation userRelation
	FanFollow    fanFollow
	SendChat     sendChat
}
