package service

import (
	"time"
)

// UpdateDashboardIncomeChart 更新Dashboard的收入图表生成逻辑
func (s *ReportService) UpdateDashboardIncomeChart() []IncomeChartData {
	incomeChart := make([]IncomeChartData, 0)

	// 生成近12个月的收入图表（按月）
	now := time.Now()
	for i := 11; i >= 0; i-- {
		month := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, time.Local)
		endOfMonth := month.AddDate(0, 1, -1)

		// 使用GetIncomeReport获取按月收入
		incomeReport, err := s.GetIncomeReport(month, endOfMonth, "month")
		if err != nil {
			// 如果出错，使用0值
			incomeChart = append(incomeChart, IncomeChartData{
				Date:   month.Format("2006-01"),
				Amount: 0,
			})
		} else {
			incomeChart = append(incomeChart, IncomeChartData{
				Date:   month.Format("2006-01"),
				Amount: incomeReport.Total,
			})
		}
	}

	return incomeChart
}