package proto

import (
	"errors"
	"strconv"
)

type GroupName string
type TaskStatus int
type TaskError string
type TaskId string

const (
	MIN_FACE_SIZE = 50
	MODE_SINGLE   = "SINGLE"
	MODE_LARGEST  = "LARGEST"
)
const (
	_ TaskStatus = iota
	CREATED
	PENDING
	STARTED
	STOPPING
	STOPPED
	COMPLETED
)

type Task struct {
	TaskId       TaskId     `bson:"task_id"`
	Uid          uint32     `bson:"uid"`
	GroupName    GroupName  `bson:"group_name"`
	Config       TaskConfig `bson:"task_config"`
	TotalCount   int        `bson:"total_count"`
	HandledCount int        `bson:"handled_count"`
	Status       TaskStatus `bson:"status"`
	FileName     string     `bson:"file_name"`
	FileExt      string     `bson:"file_ext"`
	LastError    TaskError  `bson:"last_error"`
}

type TaskConfig struct {
	FilterPose  bool   `json:"filter_pose" bson:"filter_pose"`
	FilterBlur  bool   `json:"filter_blur" bson:"filter_blur"`
	FilterCover bool   `json:"filter_cover" bson:"filter_cover"`
	MinWidth    int    `json:"min_width" bson:"min_width"`
	MinHeight   int    `json:"min_height" bson:"min_height"`
	Mode        string `json:"mode" bson:"mode"`
}

func GetValidFilter(s string) (bool, error) {
	if s == "" {
		return false, nil
	}
	filter, err := strconv.ParseBool(s)
	if err != nil {
		return false, errors.New("not a bool")
	}
	return filter, nil
}

func GetValidSizeLimit(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("not a number")
	}
	if num < MIN_FACE_SIZE {
		num = MIN_FACE_SIZE
	}
	return num, nil
}

func GetValidMode(s string) (string, error) {
	switch s {
	case "", MODE_SINGLE:
		return MODE_SINGLE, nil
	case MODE_LARGEST:
		return MODE_LARGEST, nil
	default:
		return "", errors.New("invalid mode")
	}
}
