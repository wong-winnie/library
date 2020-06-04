package timehelp

import "time"

const defaultFormat = "2006-01-02 15:04:05.000"

//日期格式化
func TimeFormat(needTime, oldFormat, newFormat string) string {
	t, _ := time.Parse(oldFormat, needTime)
	return t.Format(newFormat)
}

//计算两个时间差
//sTime, _ := time.Parse("2006-01-02 15:04:05", thirdData.Data.StartTime)
//eTime, _ := time.Parse("2006-01-02 15:04:05", thirdData.Data.EndTime)
//mTime := eTime.Sub(sTime)
