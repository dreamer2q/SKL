package skl

import "time"

type User struct {
	QQ        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserInfo

	Token  string
	Groups []*Group `gorm:"many2many:user_groups;"`
}

type Group struct {
	GroupID    int64 `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CourseName string
	//CourseID    string
	TeacherName string
	TeacherID   int64

	Users []*User `gorm:"many2many:user_groups;"`
}

type codeError struct {
	statusCode int //400,验证码不正确;401,签到码无效;200,签到成功
}

func (c codeError) Error() string {
	switch c.statusCode {
	case 400:
		return "验证码不正确"
	case 401:
		return "签到码无效"
	}
	return "error not matched"
}

type UserInfo struct {
	UserID   string `json:"id"`
	UserName string `json:"userName"`
	UserType int    `json:"userType"`
	UnitID   string `json:"unitId"`
	UnitCode string `json:"unitCode"`
	UnitName string `json:"unitName"`
	Grade    string `json:"grade"`
	ClassNo  string `json:"classNo"`
	Sex      string `json:"sex"`
	Major    string `json:"major"`
	//RoleList interface{} `json:"roleList"`
}

type SKLCheckData []struct {
	UserID   string `json:"userId"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`

	AbsentCount      int     `json:"absentCount"`
	LateCount        int     `json:"lateCount"`
	AbsentLeaveCount int     `json:"absentLeaveCount"`
	RightCount       int     `json:"rightCount"`
	LeaveCount       int     `json:"leaveCount"`
	AbsentTimeCount  float64 `json:"absentTimeCount"`
}

type SKLCheckListStruct []struct {
	CourseSchemaID      string `json:"courseSchemaId"`
	StudentID           string `json:"studentId"`
	StudentName         string `json:"studentName"`
	UnitName            string `json:"unitName"`
	ClassNo             string `json:"classNo"`
	CreateTime          int64  `json:"createTime"`
	CourseName          string `json:"courseName"`
	CourseCode          string `json:"courseCode"`
	TeacherUnit         string `json:"teacherUnit"`
	TeacherName         string `json:"teacherName"`
	SchoolYear          string `json:"schoolYear"`
	Semester            int    `json:"semester"`
	StartWeek           int    `json:"startWeek"`
	EndWeek             int    `json:"endWeek"`
	StartSection        int    `json:"startSection"`
	EndSection          int    `json:"endSection"`
	Period              string `json:"period"`
	ClassRoom           string `json:"classRoom"`
	WeekDay             int    `json:"weekDay"`
	RecordDate          int64  `json:"recordDate"`
	CheckInStatus       string `json:"checkInStatus"`
	CourseStudentStatus string `json:"courseStudentStatus"`
	UpdateBy            string `json:"updateBy"`
	UpdateMode          string `json:"updateMode"`
}
