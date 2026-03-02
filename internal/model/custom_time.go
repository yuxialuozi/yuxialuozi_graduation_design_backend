package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// CustomTime 自定义时间类型，JSON序列化为 "2006-01-02 15:04:05" 格式
type CustomTime struct {
	time.Time
}

// MarshalJSON 实现JSON序列化接口
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf(`"%s"`, ct.Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON 实现JSON反序列化接口
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// 移除引号
	str := string(data)
	if str == "null" || str == `""` {
		return nil
	}
	str = str[1 : len(str)-1]

	// 尝试多种格式解析
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	var err error
	for _, format := range formats {
		ct.Time, err = time.Parse(format, str)
		if err == nil {
			return nil
		}
	}

	return fmt.Errorf("无法解析时间: %s", str)
}

// Value 实现 driver.Valuer 接口
func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

// Scan 实现 sql.Scanner 接口
func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(time.Time); ok {
		ct.Time = v
		return nil
	}
	return fmt.Errorf("无法扫描时间: %v", value)
}

// NewCustomTime 创建自定义时间
func NewCustomTime(t time.Time) CustomTime {
	return CustomTime{Time: t}
}

// Now 返回当前时间
func NowCustomTime() CustomTime {
	return CustomTime{Time: time.Now()}
}
