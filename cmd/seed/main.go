package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"yuxialuozi_graduation_design_backend/internal/model"
)

func main() {
	cfg := loadConfig()

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal("连接数据库失败: ", err)
	}
	log.Println("数据库连接成功")

	// 清空所有表数据
	if err := cleanAllTables(db); err != nil {
		log.Fatal("清空数据失败: ", err)
	}
	log.Println("已清空所有表数据")

	// 重新创建表结构
	if err := recreateTables(db); err != nil {
		log.Fatal("重建表失败: ", err)
	}
	log.Println("已重建表结构")

	// 生成测试数据
	seedData(db)
	log.Println("测试数据生成完成！")
}

type dbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func loadConfig() *dbConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("未找到配置文件，使用默认值")
	}

	return &dbConfig{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		DBName:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.sslmode"),
	}
}

func connectDB(cfg *dbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                               logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
}

func cleanAllTables(db *gorm.DB) error {
	tables := []string{"maintenances", "fees", "contracts", "rooms", "users", "tenants"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("删除表 %s 失败: %w", table, err)
		}
		log.Printf("已删除表: %s", table)
	}
	return nil
}

func recreateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Tenant{},
		&model.Contract{},
		&model.Room{},
		&model.Fee{},
		&model.Maintenance{},
	)
}

func seedData(db *gorm.DB) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// ===== 1. 创建管理员用户 =====
	adminPwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	admin := model.User{
		Username:    "admin",
		Password:    string(adminPwd),
		Nickname:    "系统管理员",
		Phone:       "13800000001",
		Email:       "admin@tenant.com",
		Role:        "admin",
		Status:      "active",
		Permissions: []string{"read", "write", "delete", "admin"},
		TenantID:    0,
	}
	if err := db.Create(&admin).Error; err != nil {
		log.Fatal("创建管理员失败: ", err)
	}
	log.Println("已创建管理员: admin/123456")

	// ===== 2. 创建房间 =====
	buildings := []string{"A栋", "B栋", "C栋"}
	rooms := make([]model.Room, 0)
	roomIndex := 0
	for _, building := range buildings {
		for floor := 1; floor <= 6; floor++ {
			for roomNum := 1; roomNum <= 4; roomNum++ {
				roomIndex++
				area := 25.0 + float64(r.Intn(30)) + float64(r.Intn(100))/100.0
				rent := 1500.0 + float64(r.Intn(3000)) + float64(r.Intn(100))/100.0
				status := "vacant"
				if roomIndex <= 18 {
					status = "occupied"
				}
				room := model.Room{
					RoomNo:      fmt.Sprintf("%s%d%02d", building[:1], floor, roomNum),
					Building:    building,
					Floor:       floor,
					Area:        area,
					MonthlyRent: rent,
					Status:      status,
				}
				rooms = append(rooms, room)
			}
		}
	}
	if err := db.Create(&rooms).Error; err != nil {
		log.Fatal("创建房间失败: ", err)
	}
	log.Printf("已创建 %d 个房间", len(rooms))

	// ===== 3. 创建租户 =====
	tenantData := []struct {
		name  string
		phone string
		email string
	}{
		{"张伟", "13800001001", "zhangwei@email.com"},
		{"李娜", "13800001002", "lina@email.com"},
		{"王强", "13800001003", "wangqiang@email.com"},
		{"赵敏", "13800001004", "zhaomin@email.com"},
		{"陈磊", "13800001005", "chenlei@email.com"},
		{"刘洋", "13800001006", "liuyang@email.com"},
		{"孙婷", "13800001007", "sunting@email.com"},
		{"周杰", "13800001008", "zhoujie@email.com"},
		{"吴芳", "13800001009", "wufang@email.com"},
		{"郑浩", "13800001010", "zhenghao@email.com"},
		{"黄丽", "13800001011", "huangli@email.com"},
		{"林峰", "13800001012", "linfeng@email.com"},
		{"何雪", "13800001013", "hexue@email.com"},
		{"马超", "13800001014", "machao@email.com"},
		{"罗静", "13800001015", "luojing@email.com"},
		{"谢军", "13800001016", "xiejun@email.com"},
		{"韩梅", "13800001017", "hanmei@email.com"},
		{"唐明", "13800001018", "tangming@email.com"},
	}

	tenants := make([]model.Tenant, 0)
	for _, td := range tenantData {
		tenant := model.Tenant{
			Name:          td.name,
			ContactPerson: td.name,
			Phone:         td.phone,
			Email:         td.email,
			Status:        "active",
		}
		tenants = append(tenants, tenant)
	}
	if err := db.Create(&tenants).Error; err != nil {
		log.Fatal("创建租户失败: ", err)
	}
	log.Printf("已创建 %d 个租户", len(tenants))

	// ===== 4. 为前18个租户分配房间 =====
	occupiedRooms := make([]model.Room, 0)
	for i := range rooms {
		if rooms[i].Status == "occupied" && i < len(tenants) {
			rooms[i].TenantID = &tenants[i].ID
			occupiedRooms = append(occupiedRooms, rooms[i])
		}
	}
	if err := db.Save(&occupiedRooms).Error; err != nil {
		log.Fatal("分配房间失败: ", err)
	}
	log.Printf("已为 %d 个租户分配房间", len(occupiedRooms))

	// ===== 5. 创建用户账号 =====
	for i, tenant := range tenants {
		pwd, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		user := model.User{
			Username:    fmt.Sprintf("user%d", i+1),
			Password:    string(pwd),
			Nickname:    tenant.Name,
			Phone:       tenant.Phone,
			Email:       tenant.Email,
			Role:        "user",
			Status:      "active",
			Permissions: []string{"read", "write"},
			TenantID:    tenant.ID,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("创建用户 user%d 失败: %v", i+1, err)
		}
	}
	log.Printf("已创建 %d 个租户用户（用户名: user1~user%d，密码: 123456）", len(tenants), len(tenants))

	// ===== 6. 创建合同 =====
	contractStatuses := []string{"active", "active", "active", "active", "active", "active", "active", "active",
		"active", "active", "active", "active", "expired", "expired", "active", "active", "active", "active"}
	contracts := make([]model.Contract, 0)
	for i, tenant := range tenants {
		startDate := parseTime("2025-01-01")
		if i%3 == 1 {
			startDate = parseTime("2025-03-01")
		} else if i%3 == 2 {
			startDate = parseTime("2025-06-01")
		}
		endDate := parseTime("2026-01-01")
		if i%3 == 1 {
			endDate = parseTime("2026-03-01")
		} else if i%3 == 2 {
			endDate = parseTime("2026-06-01")
		}

		amount := 18000.0 + float64(r.Intn(24000))
		contract := model.Contract{
			TenantID:   tenant.ID,
			ContractNo: fmt.Sprintf("HT-2025-%04d", i+1),
			StartDate:  model.CustomTime{Time: startDate},
			EndDate:    model.CustomTime{Time: endDate},
			Amount:     amount,
			Status:     contractStatuses[i],
		}
		contracts = append(contracts, contract)
	}
	if err := db.Create(&contracts).Error; err != nil {
		log.Fatal("创建合同失败: ", err)
	}
	log.Printf("已创建 %d 份合同", len(contracts))

	// ===== 7. 创建费用记录 =====
	feeTypes := []string{"rent", "water", "electricity", "property", "parking"}
	feeStatuses := []string{"paid", "unpaid", "overdue"}

	fees := make([]model.Fee, 0)
	for i, tenant := range tenants {
		roomNo := ""
		if i < len(occupiedRooms) {
			roomNo = occupiedRooms[i].RoomNo
		}

		for month := 1; month <= 4; month++ {
			for _, feeType := range feeTypes {
				var amount float64
				switch feeType {
				case "rent":
					amount = 2000.0 + float64(r.Intn(3000))
				case "water":
					amount = 30.0 + float64(r.Intn(80))
				case "electricity":
					amount = 80.0 + float64(r.Intn(200))
				case "property":
					amount = 150.0 + float64(r.Intn(200))
				case "parking":
					amount = 200.0 + float64(r.Intn(300))
				}

				status := feeStatuses[r.Intn(3)]
				if month <= 2 {
					status = "paid"
				}

				dueDate := parseTime(fmt.Sprintf("2025-%02d-15", month+2))
				fee := model.Fee{
					TenantID: tenant.ID,
					RoomNo:   roomNo,
					FeeType:  feeType,
					Amount:   amount,
					Period:   fmt.Sprintf("2025-%02d", month+2),
					DueDate:  model.CustomTime{Time: dueDate},
					Status:   status,
				}

				if status == "paid" {
					paidDate := parseTime(fmt.Sprintf("2025-%02d-%02d", month+2, 10+r.Intn(5)))
					fee.PaidDate = &model.CustomTime{Time: paidDate}
				}

				fees = append(fees, fee)
			}
		}
	}
	if err := db.Create(&fees).Error; err != nil {
		log.Fatal("创建费用失败: ", err)
	}
	log.Printf("已创建 %d 条费用记录", len(fees))

	// ===== 8. 创建维修工单 =====
	maintenanceTypes := []string{"plumbing", "electrical", "door", "appliance", "wall", "network"}
	priorities := []string{"low", "medium", "high"}
	mStatuses := []string{"pending", "processing", "completed"}
	assignees := []string{"李师傅", "王师傅", "张师傅", "赵师傅"}
	descriptions := map[string][]string{
		"plumbing":   {"卫生间水龙头漏水，需要更换", "厨房下水管道堵塞", "卫生间马桶冲水不畅", "阳台排水管漏水"},
		"electrical": {"客厅灯具不亮，需要检修", "卧室插座没电", "厨房开关失灵", "走廊灯频繁闪烁"},
		"door":       {"入户门锁损坏，无法正常开关", "卧室门合页松动", "窗户玻璃破裂", "阳台推拉门卡顿"},
		"appliance":  {"空调不制冷，需要维修", "洗衣机排水异常", "热水器无法加热", "冰箱温度异常"},
		"wall":       {"卧室墙面渗水发霉", "客厅墙皮脱落", "卫生间瓷砖松动", "天花板有裂缝"},
		"network":    {"网络信号弱，经常断线", "网线接口松动", "路由器无法连接", "WiFi信号覆盖差"},
	}

	maintenances := make([]model.Maintenance, 0)
	for i := 0; i < 25; i++ {
		tenantIdx := r.Intn(len(tenants))
		tenant := tenants[tenantIdx]
		roomNo := ""
		if tenantIdx < len(occupiedRooms) {
			roomNo = occupiedRooms[tenantIdx].RoomNo
		}

		mType := maintenanceTypes[r.Intn(len(maintenanceTypes))]
		descs := descriptions[mType]
		desc := descs[r.Intn(len(descs))]

		priority := priorities[r.Intn(len(priorities))]
		status := mStatuses[r.Intn(len(mStatuses))]
		assignee := ""
		if status != "pending" {
			assignee = assignees[r.Intn(len(assignees))]
		}

		month := 1 + r.Intn(5)
		day := 1 + r.Intn(28)
		createdAt := parseTime(fmt.Sprintf("2025-%02d-%02d", month, day))

		m := model.Maintenance{
			TicketNo:    fmt.Sprintf("WX-2025-%04d", i+1),
			TenantID:    tenant.ID,
			RoomNo:      roomNo,
			Type:        mType,
			Description: desc,
			Priority:    priority,
			Status:      status,
			Assignee:    assignee,
			CreatedAt:   model.CustomTime{Time: createdAt},
		}

		if status == "completed" {
			completedAt := createdAt.Add(time.Duration(1+r.Intn(7)) * 24 * time.Hour)
			m.CompletedAt = &model.CustomTime{Time: completedAt}
		}

		maintenances = append(maintenances, m)
	}
	if err := db.Create(&maintenances).Error; err != nil {
		log.Fatal("创建维修工单失败: ", err)
	}
	log.Printf("已创建 %d 条维修工单", len(maintenances))

	// 打印汇总
	fmt.Println("\n========================================")
	fmt.Println("测试数据生成完成！")
	fmt.Println("========================================")
	fmt.Println("管理员: admin / 123456")
	fmt.Printf("租户用户: user1~user%d / 123456\n", len(tenants))
	fmt.Printf("房间数量: %d\n", len(rooms))
	fmt.Printf("租户数量: %d\n", len(tenants))
	fmt.Printf("合同数量: %d\n", len(contracts))
	fmt.Printf("费用记录: %d\n", len(fees))
	fmt.Printf("维修工单: %d\n", len(maintenances))
	fmt.Println("========================================")
}

func parseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		log.Fatalf("解析时间失败 %s: %v", s, err)
	}
	return t
}
