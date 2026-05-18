package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

// Project titles
var moreTitles = []string{
	"企业官网开发", "电商小程序", "移动端APP设计", "后台管理系统",
	"数据分析平台", "在线教育系统", "CRM客户管理", "库存管理系统",
	"社交媒体APP", "在线支付系统", "物流追踪系统", "医院挂号系统",
	"餐饮点餐系统", "酒店预订平台", "旅游出行APP", "健身打卡小程序",
	"外卖配送系统", "团购平台", "二手交易市场", "人才招聘网站",
	"在线测评系统", "知识付费平台", "直播电商系统", "社区团购小程序",
	"智慧停车系统", "智能家居控制", "企业OA办公", "项目管理工具",
	"文档协作平台", "代码审查系统", "API接口开发", "微服务架构",
	"数据可视化", "大数据分析", "AI客服机器人", "图像识别系统",
	"语音识别API", "推荐算法", "搜索优化", "性能调优",
	"Web安全审计", "代码重构", "技术文档编写", "测试自动化",
	"DevOps部署", "容器化改造", "云原生应用", "区块链应用",
	"NFT交易平台", "智能合约", "DeFi应用", "游戏开发",
	"小游戏开发", "元宇宙项目", "AR/VR应用", "3D建模",
	"UI设计", "UX优化", "品牌设计", "产品设计",
}

// Project descriptions
var moreDescriptions = []string{
	"我们需要一个功能完善的系统，要求：响应式设计，支持多端适配，有良好的用户体验。",
	"项目预算充足，希望找到经验丰富的开发者。要求：按时交付，代码规范，有相关经验者优先。",
	"我们是一家初创公司，需要开发我们的第一个产品。希望找到能长期合作的开发者。",
	"这是一个企业级项目，要求：架构设计合理，性能优异，可扩展性强。有相关案例的开发者优先。",
	"项目需要快速上线，希望开发者能投入足够的时间。需求清晰，有详细的PRD文档。",
}

// Status options
var statuses = []string{"published", "bidding", "in_progress", "delivered", "completed"}

// Budget types
var budgetTypes = []string{"fixed", "hourly"}

// Tech stack options
var techStackOptions = []string{
	"Go", "Python", "JavaScript", "TypeScript", "React", "Vue.js", "Node.js",
	"Java", "Spring Boot", "C#", ".NET", "PHP", "Laravel", "Ruby", "Rails",
	"PostgreSQL", "MySQL", "MongoDB", "Redis", "Docker", "Kubernetes",
	"AWS", "GCP", "Azure", "Git", "CI/CD", "Linux", "Windows",
	"HTML", "CSS", "Tailwind", "Bootstrap", "Next.js", "Nuxt.js",
	"Angular", "Svelte", "Flutter", "React Native", "Swift", "Kotlin",
}

func main() {
	// Initialize config
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	// Database connection string
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	// Connect to database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	}
	defer pool.Close()

	fmt.Println("Connected to database successfully!")

	// Get existing client IDs and category IDs
	var clientIDs []string
	var categoryIDs []string

	// Get client IDs
	rows, err := pool.Query(ctx, "SELECT id FROM users WHERE user_type IN ('client', 'both')")
	if err != nil {
		fmt.Printf("Error fetching clients: %v\n", err)
		return
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			fmt.Printf("Error scanning client: %v\n", err)
			continue
		}
		clientIDs = append(clientIDs, id)
	}
	rows.Close()
	if len(clientIDs) == 0 {
		fmt.Println("No client users found")
		return
	}

	// Get category IDs
	rows, err = pool.Query(ctx, "SELECT id FROM project_categories")
	if err != nil {
		fmt.Printf("Error fetching categories: %v\n", err)
		return
	}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			fmt.Printf("Error scanning category: %v\n", err)
			continue
		}
		categoryIDs = append(categoryIDs, id)
	}
	rows.Close()
	if len(categoryIDs) == 0 {
		fmt.Println("No categories found")
		return
	}

	fmt.Printf("Found %d clients and %d categories\n", len(clientIDs), len(categoryIDs))

	// Get current project count
	var currentCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&currentCount)
	if err != nil {
		fmt.Printf("Error getting project count: %v\n", err)
		return
	}
	fmt.Printf("Current project count: %d\n", currentCount)

	// Calculate how many more projects we need
	need := 100 - currentCount
	if need <= 0 {
		fmt.Printf("Already have %d projects, need no more\n", currentCount)
		return
	}
	fmt.Printf("Need to add %d more projects\n", need)

	// Add new projects
	addedCount := 0
	for i := 1; i <= need; i++ {
		clientID := clientIDs[rand.Intn(len(clientIDs))]
		categoryID := categoryIDs[rand.Intn(len(categoryIDs))]
		title := moreTitles[rand.Intn(len(moreTitles))]
		description := moreDescriptions[rand.Intn(len(moreDescriptions))] + fmt.Sprintf("\n\n项目编号：P%d", currentCount+i)
		status := statuses[rand.Intn(len(statuses))]
		budgetType := budgetTypes[rand.Intn(len(budgetTypes))]

		var budgetMin, budgetMax float64
		if budgetType == "fixed" {
			budgetMin = float64(rand.Intn(50000-5000+1) + 5000)
			budgetMax = budgetMin + float64(rand.Intn(20000))
		} else {
			budgetMin = float64(rand.Intn(500-100+1) + 100)
			budgetMax = budgetMin + float64(rand.Intn(200))
		}

		// Random tech stack (1-5 items)
		numTech := rand.Intn(5) + 1
		techStack := make([]string, numTech)
		for t := 0; t < numTech; t++ {
			idx := rand.Intn(len(techStackOptions))
			techStack[t] = techStackOptions[idx]
		}

		// Random view count (10-5000)
		viewCount := rand.Intn(5000-10+1) + 10
		bidCount := rand.Intn(20)
		bookmarkCount := rand.Intn(50)

		// Random created time (last 30 days)
		createdAt := time.Now().Add(-time.Duration(rand.Intn(30)*24) * time.Hour)

		// Insert project
		_, err = pool.Exec(ctx, `
			INSERT INTO projects (
				client_id, category_id, title, description, budget_min, budget_max,
				budget_type, tech_stack, status, view_count, bid_count, bookmark_count,
				created_at, featured
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8::jsonb, $9, $10, $11, $12, $13, $14)
		`,
			clientID, categoryID, title, description, budgetMin, budgetMax,
			budgetType, toJSONArray(techStack), status, viewCount, bidCount, bookmarkCount,
			createdAt, rand.Intn(10) == 0, // 10% chance to be featured
		)
		if err != nil {
			fmt.Printf("Error inserting project %d: %v\n", i, err)
			continue
		}
		addedCount++
		if i%20 == 0 {
			fmt.Printf("  Added %d/%d projects\n", i, need)
		}
	}

	fmt.Printf("\nSuccessfully added %d projects!\n", addedCount)

	// Verify final count
	var finalCount int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM projects").Scan(&finalCount)
	if err == nil {
		fmt.Printf("Total projects now: %d\n", finalCount)
	}
}

func toJSONArray(arr []string) string {
	result := "["
	for i, s := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf(`"%s"`, s)
	}
	result += "]"
	return result
}
