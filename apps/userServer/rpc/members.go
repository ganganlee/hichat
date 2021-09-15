package rpc

import (
	"context"
	"fmt"
	"hichat.zozoo.net/rpc/userGroupMembers"
)

type (
	Members struct {
	}
)

func NewMembers() *Members {
	return &Members{}
}

//添加群成员
func (u *Members) AddMember(ctx context.Context, res *userGroupMembers.AddMemberRequest, rsp *userGroupMembers.AddMemberResponse) error {

	fmt.Println("进去添加群成员rpc方法")
	rsp.Msg = "添加成功"
	return nil
}

//删除群成员
func (u *Members) DelByMemberId(ctx context.Context, res *userGroupMembers.DelByMemberIdRequest, rsp *userGroupMembers.DelByMemberIdResponse) error {
	return nil
}

//解散群
func (u *Members) DelMembers(ctx context.Context, res *userGroupMembers.DelMembersRequest, rsp *userGroupMembers.DelMembersResponse) error {
	return nil
}
