package knowledge

import (
	"sort"
	"strings"
)

// Category represents a knowledge category
type Category string

const (
	CategoryFee         Category = "费用管理"
	CategoryContract    Category = "合同管理"
	CategoryMaintenance Category = "维修管理"
	CategoryRoom        Category = "房间管理"
	CategoryPolicy      Category = "政策制度"
	CategorySystem      Category = "系统使用"
)

// KnowledgeItem represents a single knowledge article
type KnowledgeItem struct {
	ID           string     // Unique identifier
	Category     Category   // Category
	Title        string     // Title/question
	Content      string     // Detailed answer
	Keywords     []string   // Search keywords
	RelatedAPIs  []string   // Related API paths (e.g., "/api/fees", "/api/contracts")
	Priority     int        // Search priority (higher = more important)
	Tags         []string   // Tags for classification
}

// KnowledgeBase is the global knowledge base
var KnowledgeBase = []KnowledgeItem{
	// ==================== 费用管理 ====================
	{
		ID:       "fee-001",
		Category: CategoryFee,
		Title:    "租金是如何计算的？",
		Content: `租金计算方式如下：

1. **月租金计算**：月租金 = 房间单价 × 房间面积（或固定月租金额）
2. **计费周期**：租金按月计费，通常在每月1日生成账单
3. **入住不足整月**：按实际入住天数计算，日租金 = 月租金 ÷ 当月天数
4. **提前退房**：按实际居住天数结算，多付部分退还

租金金额由以下因素决定：
- 房间面积（平方米）
- 房间位置（楼层、朝向）
- 市场行情
- 合同约定

如需查询具体租金，请提供房间号，我可以帮您查询。`,
		Keywords:     []string{"租金", "租金计算", "月租", "房租", "费用", "租金标准", "房租计算"},
		RelatedAPIs:  []string{"/api/rooms", "/api/fees"},
		Priority:     10,
		Tags:         []string{"租金", "费用", "计算"},
	},
	{
		ID:       "fee-002",
		Category: CategoryFee,
		Title:    "水电费如何计算和缴纳？",
		Content: `水电费计算和缴纳说明：

**水费计算**：
- 按实际用水量计算（立方米）
- 水价通常为 2-5 元/立方米（根据地区定价）
- 每月根据水表读数核算

**电费计算**：
- 按实际用电量计算（千瓦时，kWh）
- 电价按阶梯电价计算：
  - 第一档（0-200度）：0.5-0.6元/度
  - 第二档（201-400度）：0.6-0.7元/度
  - 第三档（400度以上）：0.8元/度以上

**缴纳方式**：
1. 线上支付：通过本系统在线支付
2. 线下支付：到物业管理处缴纳
3. 代扣代缴：绑定银行卡自动扣款

**缴纳时间**：每月10日前缴纳上月光电费，逾期将产生滞纳金。`,
		Keywords:     []string{"水电费", "水费", "电费", "水电", "能源费", "阶梯电价", "用水", "用电", "缴费", "支付"},
		RelatedAPIs:  []string{"/api/fees"},
		Priority:     10,
		Tags:         []string{"水电费", "费用", "缴纳"},
	},
	{
		ID:       "fee-003",
		Category: CategoryFee,
		Title:    "物业费包含哪些内容？",
		Content: `物业费包含以下服务内容：

**公共设施维护**：
- 电梯维护和保养
- 公共区域清洁（楼道、走廊、花园等）
- 安保服务
- 消防设施维护

**绿化养护**：
- 小区绿化维护
- 园艺修剪
- 病虫害防治

**垃圾清运**：
- 生活垃圾清运
- 建筑垃圾处理

**公共照明**：
- 公共区域照明
- 景观照明

**其他服务**：
- 社区活动组织
- 停车管理
- 设施设备维修

物业费标准通常为 2-10 元/平方米/月，具体金额见合同约定。`,
		Keywords:     []string{"物业费", "物业服务", "物业", "管理费", "服务费"},
		RelatedAPIs:  []string{"/api/fees"},
		Priority:     8,
		Tags:         []string{"物业费", "服务"},
	},
	{
		ID:       "fee-004",
		Category: CategoryFee,
		Title:    "逾期未缴费会怎样？",
		Content: `逾期未缴费的处理方式：

**短期逾期（1-7天）**：
- 系统发送催缴通知
- 暂无滞纳金

**中期逾期（8-30天）**：
- 发送正式催缴函
- 滞纳金按欠费金额的 0.05%/天 计算
- 限制部分系统功能

**长期逾期（30天以上）**：
- 发送律师函
- 滞纳金累计计算
- 可能影响个人信用记录
- 严重情况下启动法律程序

**注意事项**：
- 如有困难，请及时与物业沟通
- 经济困难可申请分期付款
- 建议设置自动扣款避免逾期

请尽快缴纳欠费，如有疑问请联系物业。`,
		Keywords:     []string{"逾期", "滞纳金", "欠费", "逾期缴费", "欠费处理", "催缴", "罚款"},
		RelatedAPIs:  []string{"/api/fees"},
		Priority:     7,
		Tags:         []string{"逾期", "滞纳金", "欠费"},
	},
	{
		ID:       "fee-005",
		Category: CategoryFee,
		Title:    "如何查询我的账单？",
		Content: `查询账单的方法：

**线上查询**：
1. 登录租户端系统
2. 进入"我的账单"页面
3. 查看所有费用明细，包括：
   - 租金
   - 水电费
   - 物业费
   - 其他费用

**账单内容包括**：
- 费用类型（租金/水费/电费/物业费）
- 费用金额
- 所属周期（年月）
- 应缴日期
- 缴费状态（已缴/未缴/逾期）

**导出账单**：
支持导出 PDF 格式的账单明细，可用于报销或留存。

如有账单疑问，请联系物业核实。`,
		Keywords:     []string{"账单", "查询账单", "费用明细", "账单查询", "看账单", "查看费用"},
		RelatedAPIs:  []string{"/api/tenant/fees"},
		Priority:     9,
		Tags:         []string{"账单", "查询"},
	},

	// ==================== 合同管理 ====================
	{
		ID:       "contract-001",
		Category: CategoryContract,
		Title:    "如何签订租房合同？",
		Content: `签订租房合同的流程：

**签订前准备**：
1. 选择合适的房间
2. 确认租金和付款方式
3. 准备身份证件

**签订流程**：
1. 双方协商合同条款
2. 物业提供合同文本
3. 租户仔细阅读合同内容
4. 确认无误后签字/盖章
5. 缴纳首月租金和押金
6. 交付钥匙，完成入住

**合同内容**：
- 房屋地址和房间号
- 租期（起始和终止日期）
- 租金金额和付款方式
- 押金金额和退还条件
- 双方权利和义务
- 违约责任
- 续租和退租条款

**注意事项**：
- 仔细阅读所有条款
- 保留合同副本
- 核对房间信息是否正确`,
		Keywords:     []string{"签订合同", "合同签订", "签合同", "租赁合同", "合同", "签约"},
		RelatedAPIs:  []string{"/api/contracts", "/api/tenant/contracts"},
		Priority:     9,
		Tags:         []string{"合同", "签订"},
	},
	{
		ID:       "contract-002",
		Category: CategoryContract,
		Title:    "租期一般是多长时间？",
		Content: `租期说明：

**常见租期**：
- 短期租赁：3-6个月
- 标准租期：1年（最常见）
- 长期租赁：2年及以上

**租期选择建议**：
- 首次租房建议签1年，了解居住体验后再决定
- 长期居住可签2-3年，通常有租金优惠
- 短期过渡可选择短租或月租

**租期计算**：
- 租期从入住当天开始计算
- 租期最后一天为合同终止日
- 退房应在租期结束前办理

**提前终止**：
- 需要提前30天书面通知
- 可能需要支付违约金
- 押金退还按合同约定执行`,
		Keywords:     []string{"租期", "合同期限", "租赁期限", "一年", "长期", "短期"},
		RelatedAPIs:  []string{"/api/contracts", "/api/tenant/contracts"},
		Priority:     8,
		Tags:         []string{"租期", "合同"},
	},
	{
		ID:       "contract-003",
		Category: CategoryContract,
		Title:    "租金调整是怎么规定的？",
		Content: `租金调整规定：

**调整条件**：
- 合同到期续签时可调整
- 市场行情大幅变化时
- 房屋装修或升级后

**调整规则**：
- 涨幅通常不超过上一年度的10%
- 提前30天书面通知租户
- 新租金在下一租期生效

**不调整情况**：
- 合同期内租金固定
- 除非双方协商同意

**续租租金**：
- 续租时根据市场行情和房屋状况确定
- 通常会有一定的涨幅
- 老租户通常可享受优惠

**协商建议**：
- 可与物业协商租金
- 长期租户可能有优惠
- 按时缴费有助于谈判`,
		Keywords:     []string{"租金调整", "涨租", "涨房租", "加租", "租金变化", "续租租金"},
		RelatedAPIs:  []string{"/api/contracts"},
		Priority:     7,
		Tags:         []string{"租金", "调整"},
	},
	{
		ID:       "contract-004",
		Category: CategoryContract,
		Title:    "如何办理退租手续？",
		Content: `退租办理流程：

**提前通知**：
- 提前30天书面通知物业
- 说明退租日期
- 确认最后缴费日期

**退租前准备**：
1. 缴清所有费用（租金、水电费、物业费等）
2. 检查房间设施设备
3. 清理个人物品
4. 恢复房间原状（如有改动）

**退租当日流程**：
1. 物业检查房间
2. 核对设施设备完好情况
3. 结算费用（如有损坏需赔偿）
4. 退还押金
5. 交回钥匙和门禁卡
6. 签署退房确认单

**押金退还**：
- 正常退租：全额退还押金
- 损坏赔偿：从押金中扣除
- 欠费：从押金中扣除

**注意事项**：
- 提前通知是义务
- 押金收据要保存好
- 结清费用是退款前提`,
		Keywords:     []string{"退租", "退房", "退租手续", "退房手续", "退房流程", "终止合同"},
		RelatedAPIs:  []string{"/api/contracts"},
		Priority:     8,
		Tags:         []string{"退租", "合同"},
	},
	{
		ID:       "contract-005",
		Category: CategoryContract,
		Title:    "如何办理续租？",
		Content: `续租办理流程：

**提前申请**：
- 租期到期前30-60天申请
- 提前申请可优先保留房间

**续租流程**：
1. 向物业提出续租申请
2. 物业审核续租资格
3. 确认新租金（可能有调整）
4. 签署续租合同
5. 缴纳新租期费用

**续租优惠**：
- 老租户通常可享受折扣
- 长期续租有更多优惠
- 按时缴费记录良好有优惠

**注意事项**：
- 提前申请，避免房间被出租
- 准备好新租期费用
- 仔细阅读新合同条款

**不续租处理**：
- 如不续租，提前通知物业
- 按退租流程办理`,
		Keywords:     []string{"续租", "续签", "合同续签", "续约", "继续租"},
		RelatedAPIs:  []string{"/api/contracts"},
		Priority:     8,
		Tags:         []string{"续租", "合同"},
	},

	// ==================== 维修管理 ====================
	{
		ID:       "maint-001",
		Category: CategoryMaintenance,
		Title:    "如何提交维修申请？",
		Content: `维修申请流程：

**申请方式**：
1. 登录租户端系统
2. 进入"维修工单"页面
3. 点击"提交维修申请"
4. 填写维修信息
5. 提交申请

**申请内容**：
- 维修类型（电气、水管、家电、家具、其他）
- 问题描述（详细描述故障情况）
- 联系方式
- 方便上门的时间

**维修类型说明**：
- **电气类**：灯具、插座、开关、断路器等
- **水管类**：水龙头、马桶、管道漏水等
- **家电类**：空调、冰箱、洗衣机等
- **家具类**：门锁、窗把手、柜门等
- **其他**：其他需要维修的问题

**提交后**：
- 系统生成工单编号
- 物业将在24小时内响应
- 维修人员联系您预约时间

**紧急维修**：
- 水管爆裂、严重漏水
- 电路故障、停电
- 门锁故障（影响安全）
紧急情况请同时电话联系物业。`,
		Keywords:     []string{"维修", "报修", "申请维修", "提交维修", "维修申请", "维修请求", "坏", "故障", "损坏"},
		RelatedAPIs:  []string{"/api/tenant/maintenance"},
		Priority:     10,
		Tags:         []string{"维修", "报修"},
	},
	{
		ID:       "maint-002",
		Category: CategoryMaintenance,
		Title:    "维修响应时间是多长？",
		Content: `维修响应时间：

**按维修类型**：

| 维修类型 | 响应时间 | 处理时间 |
|---------|---------|---------|
| 紧急维修 | 2小时内 | 24小时内完成 |
| 电气故障 | 24小时内 | 1-3天 |
| 水管维修 | 24小时内 | 1-2天 |
| 家电维修 | 1-3天 | 3-7天 |
| 家具维修 | 3-7天 | 5-10天 |
| 其他维修 | 3-7天 | 7-14天 |

**紧急维修范围**：
- 水管爆裂、严重漏水
- 电路故障、跳闸
- 门锁故障无法进入
- 空调/暖气故障（极端天气）
- 燃气泄漏

**影响处理时间的因素**：
- 配件采购时间
- 维修复杂度
- 维修人员排期

**工单查询**：
可在系统中查看工单状态：待处理 → 处理中 → 已完成`,
		Keywords:     []string{"维修时间", "响应时间", "维修周期", "多久", "等待", "处理时间"},
		RelatedAPIs:  []string{"/api/tenant/maintenance"},
		Priority:     8,
		Tags:         []string{"维修", "时间"},
	},
	{
		ID:       "maint-003",
		Category: CategoryMaintenance,
		Title:    "哪些情况需要自费维修？",
		Content: `自费维修说明：

**需要租户自费的情况**：
1. **人为损坏**：
   - 因操作不当造成的损坏
   - 故意破坏
   - 疏忽造成的故障

2. **私人物品维修**：
   - 租户自带家电故障
   - 私人物品损坏

3. **改造成本**：
   - 额外装修请求
   - 改造费用

**物业负责维修的情况**：
1. **正常使用磨损**：
   - 灯泡自然损坏
   - 水龙头老化
   - 门锁正常磨损

2. **房屋质量问题**：
   - 管道漏水（非人为）
   - 墙面脱落
   - 门窗故障（非人为）

3. **公共设施**：
   - 电梯故障
   - 公共区域照明
   - 消防设施

**判断标准**：
- 是否正常使用
- 是否人为造成
- 是否在质保期内

如有争议，可与物业协商或查看合同约定。`,
		Keywords:     []string{"自费维修", "付费维修", "自费", "费用承担", "谁付钱", "维修费用"},
		RelatedAPIs:  []string{"/api/tenant/maintenance"},
		Priority:     7,
		Tags:         []string{"维修", "费用"},
	},
	{
		ID:       "maint-004",
		Category: CategoryMaintenance,
		Title:    "如何查看维修工单状态？",
		Content: `查看维修工单状态：

**查询方式**：
1. 登录租户端系统
2. 进入"维修工单"页面
3. 查看所有维修记录

**工单状态说明**：
- **待处理**：物业已收到申请，等待分配维修人员
- **处理中**：维修人员正在处理中
- **已完成**：维修已完成
- **已取消**：工单已取消

**工单信息包括**：
- 工单编号
- 维修类型
- 问题描述
- 提交时间
- 处理状态
- 维修人员
- 完成时间
- 维修备注

**注意事项**：
- 保持联系方式畅通
- 维修人员会提前联系预约
- 完成后可在系统中评价`,
		Keywords:     []string{"维修进度", "工单状态", "维修状态", "查看维修", "查询维修", "维修记录"},
		RelatedAPIs:  []string{"/api/tenant/maintenance"},
		Priority:     7,
		Tags:         []string{"维修", "查询"},
	},

	// ==================== 房间管理 ====================
	{
		ID:       "room-001",
		Category: CategoryRoom,
		Title:    "房间状态有哪些？",
		Content: `房间状态说明：

| 状态 | 说明 |
|------|------|
| **空置** | 房间空闲，可入住 |
| **已入住** | 已有租户居住 |
| **维护中** | 正在装修或维修，暂不可入住 |

**空置房间**：
- 可随时安排入住
- 设施设备完好
- 已通过检查

**已入住房间**：
- 租户正在居住
- 合同期内不可更换
- 如需换房需退租后重新申请

**维护中房间**：
- 正在施工或维修
- 暂不接受入住
- 完工后会转为空置

**房间配置**：
- 基本配置：床、衣柜、桌椅、空调
- 厨房：灶台、橱柜
- 卫生间：马桶、淋浴、洗手台
- 其他：阳台、窗帘、家电`,
		Keywords:     []string{"房间状态", "空置", "入住", "维护中", "房间类型", "状态"},
		RelatedAPIs:  []string{"/api/rooms"},
		Priority:     6,
		Tags:         []string{"房间", "状态"},
	},
	{
		ID:       "room-002",
		Category: CategoryRoom,
		Title:    "入住和退房流程是什么？",
		Content: `入住流程：

1. **签署合同**：签订租房合同，缴纳首月租金和押金
2. **办理入住**：
   - 领取钥匙、门禁卡
   - 确认房间设施设备
   - 签署入住确认单
3. **熟悉环境**：
   - 了解公共区域位置
   - 学习门禁系统使用
   - 记住物业联系方式

退房流程：

1. **提前通知**：租期结束前30天通知物业
2. **缴清费用**：结清所有租金、水电费等
3. **物品清理**：
   - 清理个人物品
   - 恢复房间原状
   - 清理垃圾
4. **房间检查**：
   - 物业验收房间
   - 检查设施设备
   - 确认有无损坏
5. **押金退还**：
   - 无损坏：全额退还
   - 有损坏：扣除维修费后退还
6. **交回物品**：交还钥匙、门禁卡

**注意事项**：
- 入住时仔细检查房间，有问题及时报告
- 退房时拍照留证
- 保留缴费凭证`,
		Keywords:     []string{"入住", "退房", "入住流程", "退房流程", "搬入", "搬出"},
		RelatedAPIs:  []string{"/api/rooms"},
		Priority:     8,
		Tags:         []string{"入住", "退房", "流程"},
	},

	// ==================== 政策制度 ====================
	{
		ID:       "policy-001",
		Category: CategoryPolicy,
		Title:    "押金是如何收取和退还的？",
		Content: `押金收取和退还规则：

**押金收取**：
- 金额：通常为1-2个月租金
- 时间：签订合同时缴纳
- 用途：作为履行合同的保证

**押金退还条件**：
1. 租期正常结束
2. 缴清所有费用
3. 房间无损坏
4. 交还所有钥匙和门禁卡

**押金扣除情况**：
- 房间设施损坏
- 未结清的费用
- 违反合同约定的违约金

**退还时间**：
- 退房检查完成后7个工作日内
- 通常以原支付方式退还

**注意事项**：
- 入住时仔细检查房间并记录
- 保留押金收据
- 退房时拍照留证
- 如有损坏，提前了解维修费用`,
		Keywords:     []string{"押金", "押金退还", "保证金", "押金扣除", "押金退还时间"},
		RelatedAPIs:  []string{"/api/contracts"},
		Priority:     8,
		Tags:         []string{"押金", "政策"},
	},
	{
		ID:       "policy-002",
		Category: CategoryPolicy,
		Title:    "租户的权益有哪些？",
		Content: `租户权益说明：

**居住权**：
- 在租期内享有安静居住权
- 物业不得随意进入房间（紧急情况除外）
- 隐私权保护

**知情权**：
- 了解租金调整原因
- 了解各项费用明细
- 了解维修进度

**安全保障**：
- 小区24小时安保
- 公共区域监控
- 消防设施完备

**维修服务**：
- 正常使用下的免费维修
- 紧急维修响应
- 维修质量保障

**退租权利**：
- 提前通知后可退租
- 押金按规定退还
- 公平处理争议

**投诉权利**：
- 对物业服务不满意可投诉
- 投诉渠道：物业前台、电话、在线系统
- 合理投诉会得到处理

**维权途径**：
- 协商解决
- 消费者协会
- 房管部门
- 法律途径`,
		Keywords:     []string{"权益", "权利", "租户权益", "隐私", "安全", "保障"},
		RelatedAPIs:  []string{},
		Priority:     6,
		Tags:         []string{"权益", "政策"},
	},

	// ==================== 系统使用 ====================
	{
		ID:       "system-001",
		Category: CategorySystem,
		Title:    "如何登录系统？",
		Content: `登录系统方法：

**管理端登录**：
1. 访问系统登录页面
2. 输入用户名和密码
3. 点击登录按钮
4. 首次登录可修改密码

**租户端登录**：
- 使用物业分配的账号登录
- 用户名通常是手机号或工号
- 初始密码由物业提供

**忘记密码**：
- 点击"忘记密码"链接
- 输入注册手机号
- 通过验证码重置密码
- 或联系物业协助重置

**账号安全**：
- 定期修改密码
- 不要泄露密码
- 退出登录时点击"退出"按钮
- 不要在公共电脑上保存密码

**系统要求**：
- 推荐使用Chrome、Firefox、Safari浏览器
- 开启JavaScript
- 允许Cookies`,
		Keywords:     []string{"登录", "注册", "密码", "忘记密码", "账号", "用户名"},
		RelatedAPIs:  []string{"/api/auth/login"},
		Priority:     9,
		Tags:         []string{"登录", "系统"},
	},
	{
		ID:       "system-002",
		Category: CategorySystem,
		Title:    "系统主要功能有哪些？",
		Content: `系统主要功能：

**管理端功能**：
1. **仪表盘**：查看整体运营数据
2. **租户管理**：管理所有租户信息
3. **合同管理**：合同签订、续签、终止
4. **房间管理**：房间状态、分配管理
5. **费用管理**：费用生成、收缴管理
6. **维修管理**：工单分配、处理进度
7. **报表统计**：各类数据统计分析

**租户端功能**：
1. **我的首页**：查看个人数据概览
2. **个人信息**：查看和修改个人信息
3. **我的合同**：查看合同详情
4. **我的账单**：查看和缴纳费用
5. **维修工单**：提交和查看维修申请
6. **AI助手**：智能问答服务

**使用技巧**：
- 页面右上角有操作提示
- 数据表格支持搜索和筛选
- 重要信息可导出保存`,
		Keywords:     []string{"系统功能", "功能介绍", "菜单", "首页", "功能模块", "操作"},
		RelatedAPIs:  []string{"/api/tenant/dashboard", "/api/tenant/profile", "/api/tenant/contracts", "/api/tenant/fees", "/api/tenant/maintenance"},
		Priority:     7,
		Tags:         []string{"系统", "功能"},
	},
	{
		ID:       "system-003",
		Category: CategorySystem,
		Title:    "如何修改个人信息？",
		Content: `修改个人信息：

**可修改的信息**：
- 联系方式（手机号、邮箱）
- 紧急联系人
- 头像（如有）

**修改步骤**：
1. 登录租户端系统
2. 进入"个人信息"页面
3. 点击"编辑"按钮
4. 修改需要更新的内容
5. 点击"保存"按钮

**不可修改的信息**：
- 姓名（如需修改请联系物业）
- 身份证号
- 房间号
- 合同信息

**注意事项**：
- 手机号用于接收重要通知
- 确保联系方式准确
- 如更换手机号请及时更新

**联系物业**：
如需修改不可自行修改的信息，请联系物业协助处理。`,
		Keywords:     []string{"修改信息", "个人信息", "修改资料", "更新信息", "联系方式", "手机号"},
		RelatedAPIs:  []string{"/api/tenant/profile"},
		Priority:     6,
		Tags:         []string{"个人信息", "系统"},
	},
}

// SearchResult represents a search result
type SearchResult struct {
	Item  KnowledgeItem
	Score float64
}

// Search searches the knowledge base for relevant items
func Search(query string, topK int) []SearchResult {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil
	}

	queryWords := strings.Fields(query)
	results := make([]SearchResult, 0)

	for _, item := range KnowledgeBase {
		score := calculateScore(queryWords, item)
		if score > 0 {
			results = append(results, SearchResult{
				Item:  item,
				Score: score,
			})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}

	return results
}

// calculateScore calculates the relevance score between query and knowledge item
func calculateScore(queryWords []string, item KnowledgeItem) float64 {
	var score float64

	titleLower := strings.ToLower(item.Title)
	contentLower := strings.ToLower(item.Content)

	// Check title match (highest weight)
	titleMatchCount := 0
	for _, word := range queryWords {
		if strings.Contains(titleLower, word) {
			titleMatchCount++
			score += 3.0 // Title match is very important
		}
	}

	// Check content match
	contentMatchCount := 0
	for _, word := range queryWords {
		if strings.Contains(contentLower, word) {
			contentMatchCount++
			score += 1.0
		}
	}

	// Check keyword match (highest weight per match)
	keywordMatchCount := 0
	for _, word := range queryWords {
		for _, keyword := range item.Keywords {
			if strings.Contains(strings.ToLower(keyword), word) {
				keywordMatchCount++
				score += 2.0
			}
		}
	}

	// Check tag match
	for _, word := range queryWords {
		for _, tag := range item.Tags {
			if strings.Contains(strings.ToLower(tag), word) {
				score += 1.5
			}
		}
	}

	// Bonus for priority
	score += float64(item.Priority) * 0.1

	// Require at least some match
	if titleMatchCount == 0 && contentMatchCount == 0 && keywordMatchCount == 0 {
		return 0
	}

	return score
}

// GetItemByID retrieves a knowledge item by its ID
func GetItemByID(id string) *KnowledgeItem {
	for i := range KnowledgeBase {
		if KnowledgeBase[i].ID == id {
			return &KnowledgeBase[i]
		}
	}
	return nil
}

// GetByCategory retrieves all knowledge items in a category
func GetByCategory(category Category) []KnowledgeItem {
	var items []KnowledgeItem
	for i := range KnowledgeBase {
		if KnowledgeBase[i].Category == category {
			items = append(items, KnowledgeBase[i])
		}
	}
	return items
}

// BuildContext builds a knowledge context string from search results
func BuildContext(results []SearchResult) string {
	if len(results) == 0 {
		return "（知识库中未找到相关信息）"
	}

	var builder strings.Builder
	builder.WriteString("【相关知识】\n\n")

	for i, result := range results {
		builder.WriteString(strings.Repeat("=", 40))
		builder.WriteString("\n")
		builder.WriteString(result.Item.Title)
		builder.WriteString(" [")
		builder.WriteString(string(result.Item.Category))
		builder.WriteString("]\n\n")
		builder.WriteString(result.Item.Content)

		if len(result.Item.RelatedAPIs) > 0 {
			builder.WriteString("\n\n相关接口：")
			builder.WriteString(strings.Join(result.Item.RelatedAPIs, ", "))
		}

		if i < len(results)-1 {
			builder.WriteString("\n\n")
		}
	}

	return builder.String()
}