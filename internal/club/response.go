package club

import (
	"club-tracker/internal/shift"
	"strconv"
)

// Response returns to client after processing event
type Response struct {
	Time     shift.Time
	Code     int
	UserName string
	Table    int
	Err      string
}

func newResponse(time shift.Time, code int, userName string, table int, err string) Response {
	return Response{
		Time:     time,
		Code:     code,
		UserName: userName,
		Table:    table,
		Err:      err,
	}
}

func responseErr(err string, time shift.Time) Response {
	return Response{Code: Error, Err: err, Time: time}
}

func responseZero() Response {
	return Response{Code: Zero}
}

func (r Response) String() string {
	switch r.Code {
	case ClientForcedToLeave:
		return r.Time.String() + " " + strconv.Itoa(r.Code) + " " + r.UserName
	case ClientGotToTable:
		return r.Time.String() + " " + strconv.Itoa(r.Code) + " " + r.UserName + " " + strconv.Itoa(r.Table)
	case Error:
		return r.Time.String() + " " + strconv.Itoa(r.Code) + " " + r.Err
	default:
		return ""
	}
}
