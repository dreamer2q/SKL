package skl

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
