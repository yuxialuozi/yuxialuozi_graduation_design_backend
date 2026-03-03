package service

import (
	"time"
)

// UpdateDashboardIncomeChart 更新Dashboard的收入图表生成逻辑
func (s *ReportService) UpdateDashboardIncomeChart() []IncomeChartData {
	incomeChart := make([]IncomeChartData, 0)

	// 生成2025年全年的收入图表（按月）
	for month := 1; month <= 12; month++ {
		startOfMonth := time.Date(2025, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		// 使用GetIncomeReport获取按月收入
		incomeReport, err := s.GetIncomeReport(startOfMonth, endOfMonth, "month")
		if err != nil {
			// 如果出错，使用0值
			incomeChart = append(incomeChart, IncomeChartData{
				Date:   startOfMonth.Format("2006-01"),
				Amount: 0,
			})
		} else {
			incomeChart = append(incomeChart, IncomeChartData{
				Date:   startOfMonth.Format("2006-01"),
				Amount: incomeReport.Total,
			})
		}
	}

	return incomeChart
}
