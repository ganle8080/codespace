package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func main() {
	var xlsxFilePath = flag.String("f", "", "指定文件名")
	flag.Parse()

	file, _ := xlsx.OpenFile(*xlsxFilePath)

	sheet := file.Sheets[0]
	fmt.Printf("正在处理工作表: %s\n", sheet.Name)

	// 从第5行开始（索引4，因为从0开始计数）
	rowIndex := 4 // 第5行
	iCol := 8     // I列（索引8，A=0, B=1,..., I=8）
	jCol := 9     // J列（索引9）
	kCol := 10
	var total float64

	for {
		row, err := sheet.Row(rowIndex)
		if err != nil {
			break // 到达文件末尾或出错时退出循环
		}

		// 获取I列和J列的单元格
		iCell := row.GetCell(iCol)
		jCell := row.GetCell(jCol)
		kCell := row.GetCell(kCol)

		// 获取单元格的值
		iValue, _ := iCell.FormattedValue()
		jValue, _ := jCell.FormattedValue()
		kValue, _ := kCell.FormattedValue()

		if iValue != "未打卡" && jValue != "未打卡" && kValue == "2" {
			d, _ := calculateExtraTime(iValue, jValue)
			total += d
		}

		// 如果两个单元格都为空，则停止读取
		if iValue == "" && jValue == "" {
			break
		}

		// 将成对数据添加到列表中
		rowIndex++

	}

	fmt.Printf("本月额外工时为: %v\n", total)
}

// 将字符串时间("HH:MM")转换为当天的time.Time
func parseTime(timeStr string) (time.Time, error) {
	// 使用当前日期作为基准，因为我们只关心时间部分
	now := time.Now()
	return time.ParseInLocation("15:04", timeStr, now.Location())
}

// 计算两个时间点之间的额外工作时间（不在08:30-18:30之间）
// 返回浮点数小时数，负值表示早于工作时间
func calculateExtraTime(startStr, endStr string) (float64, error) {
	// 解析时间字符串
	startTime, err := parseTime(startStr)
	if err != nil {
		return 0, fmt.Errorf("无效的开始时间: %v", err)
	}

	endTime, err := parseTime(endStr)
	if err != nil {
		return 0, fmt.Errorf("无效的结束时间: %v", err)
	}

	// 定义工作时间段
	workStart, _ := parseTime("08:30")
	// workEnd, _ := parseTime("18:00")

	eTime, _ := parseTime("18:30")

	var extraHours float64

	// 1. 计算开始时间之前的额外时间+
	if startTime.Before(workStart) {
		diff := workStart.Sub(startTime)
		extraHours += diff.Hours()
	}

	if startTime.After(workStart) {
		diff := workStart.Sub(startTime)
		extraHours -= diff.Hours()
	}

	// 2. 计算结束时间之后的额外时间
	if endTime.After(eTime) {
		diff := endTime.Sub(eTime)
		extraHours += diff.Hours()
	}

	return roundToTwoDecimal(extraHours), nil
}

// 四舍五入保留两位小数
func roundToTwoDecimal(num float64) float64 {
	return math.Round(num*100) / 100
}
