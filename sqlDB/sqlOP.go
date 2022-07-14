package main //指明文件属于哪个包，package main表示一个可独立执行的程序

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql" //即后续接口需要用到的包，或者包中的函数
	"gorm.io/gorm"
	"time"
)

type Users struct {
	UserId    uint
	Name      string
	Type      string
	UserRight uint
}

type Goods struct {
	GoodsId     uint
	GoodsName   string
	GoodsPrice  float64
	GoodsStock  uint
	InputUserID uint
}

type Promotion struct {
	PromotionId        uint
	PromotionStartDate time.Time
	PromotionEndDate   time.Time
	PromotionRule      string
	InputUserId        uint
}

type GoodsAndPromotion struct {
	GoodsId        uint
	PromotionId    uint
	PromotionPrice float64
	GoodNumber     uint
}

type TabelInfor struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

var (
	dsn   = "root:123456@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
)

func userLogin(loginName string) {
	fmt.Println("---------UserLogin Func Start---------")
	//先检测用户名是否存在，不存在则进行创建，并返回用户权限
	//若用户已经存在，则直接返回用户权限
	//其余查询结果，则返回error
	var users []Users
	var usersAgain []Users
	//查询用户名是否已存在
	err := db.Where("name = ?", loginName).Find(&users).Error
	if err != nil {
		panic(err)
	} else {
		//若数据库users表中返回数据长度为0，即无此用户信息存在
		//则进行用户新增创建操作
		if len(users) == 0 {
			u1 := Users{Name: loginName, Type: "男", UserRight: 1}
			//自动迁移
			db.AutoMigrate(&Users{})
			db.Create(&u1) //创建
			fmt.Println("登录用户为:", loginName, ",用户首次登陆,自动进行创建")
			db.Where("name = ?", loginName).Find(&usersAgain)
			//返回用户权限，用于后续操作验证
			userRights := usersAgain[0].UserRight
			fmt.Println(loginName, "的用户权限为:", userRights)
		} else {
			//若用户已存在，则直接返回用户权限
			if len(users) == 1 {
				fmt.Println("登录用户为:", loginName)
				userRights := users[0].UserRight
				fmt.Println(loginName, "的用户权限为:", userRights)
			} else {
				panic(err) //待补充
			}
		}
	}
	fmt.Println("---------UserLogin Func End---------")
}

func tableQuery(args ...interface{}) {
	fmt.Println("---------Query Func Start---------")
	//传参为表名，主键值，及其他值
	var tableName []TabelInfor
	var lineName []string
	var PriName string
	var argsMap map[string]interface{}
	//获取table name
	str, ok := args[0].(string)
	if ok != true {
		panic(errors.New("error"))
	}
	//主键值
	strPri := args[1]
	//获取表字段名
	sqlNum := "desc " + str
	db.Raw(sqlNum).Scan(&tableName)
	for _, v := range tableName {
		//数组append()
		lineName = append(lineName, v.Field)
	}
	fmt.Println(str+"表字段包含：", lineName, "，请按此顺序进行参数传递！")
	//创建map 这里创建map的原因是，gorm中查询语句。scan()需要添加对应table的字段名称，且传参个数需要与table中相符
	//update，delete中不需要，甚至定义一个空字符串都可
	argsMap = make(map[string]interface{})
	for k, v := range lineName {
		argsMap[v] = k
	}
	//查询主键
	sqlPri := "SELECT column_name FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` WHERE table_name=" + "'" + str + "'" +
		" AND CONSTRAINT_SCHEMA='mysql' AND constraint_name='PRIMARY'"
	db.Raw(sqlPri).Scan(&PriName)
	queryStr := PriName + " = ?"
	err := db.Table(str).Where(queryStr, strPri).Find(&argsMap).Error
	if err != nil {
		panic(err)
	}
	fmt.Println("查询结果为:")
	for k, v := range argsMap {
		fmt.Println(k, ":", v)
	}
	fmt.Println("---------Query Func End---------")
}

func userInsert(args ...interface{}) {
	fmt.Println("---------UserInsert Func Start---------")
	var tableName []TabelInfor
	var lineName []string
	var argsMap map[string]interface{}
	argsMap = make(map[string]interface{})
	//获取table name
	str, ok := args[0].(string)
	if ok != true {
		panic(errors.New("error"))
	}
	//获取表字段名
	sqlNum := "desc " + str
	db.Raw(sqlNum).Scan(&tableName)
	for _, v := range tableName {
		//数组append()
		lineName = append(lineName, v.Field)
	}
	fmt.Println(str+"表字段包含：", lineName, "，请按此顺序进行参数传递！")
	for k, v := range lineName {
		//调试时使用
		//fmt.Println(k, v, args[k+1])
		if args[k+1] == "" {
			continue
		}
		argsMap[v] = args[k+1]
	}
	err := db.Table(str).Create(argsMap).Error
	if err != nil {
		panic(err)
	} else {
		fmt.Println("数据新增成功！")
	}

	fmt.Println("---------UserInsert Func End---------")
}

func userUpgrade(args ...interface{}) {
	fmt.Println("---------Upgrade Func Start---------")
	//传参依次为表名，主键值，及其他参数
	var tableName []TabelInfor
	var lineName []string
	var PriName string
	var argsMap map[string]interface{}
	argsMap = make(map[string]interface{})
	//表名获取
	str, ok := args[0].(string)
	if ok != true {
		panic(errors.New("error"))
	}
	//主键值
	strPri := args[1]
	//获取数据type
	//fmt.Printf("%T", str)
	//获取表字段名
	sqlNum := "desc " + str
	db.Raw(sqlNum).Scan(&tableName)
	for _, v := range tableName {
		lineName = append(lineName, v.Field)
	}
	fmt.Println(str+"表字段包含：", lineName, "，请按此顺序进行参数传递！")
	//将入参进行key-value赋值
	for k, v := range lineName {
		//调试时使用
		//fmt.Println(k, v, args[k+1])
		if args[k+1] == "" {
			continue
		}
		argsMap[v] = args[k+1]
	}
	//查询主键
	sqlPri := "SELECT column_name FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` WHERE table_name=" + "'" + str + "'" +
		" AND CONSTRAINT_SCHEMA='mysql' AND constraint_name='PRIMARY'"
	db.Raw(sqlPri).Scan(&PriName)

	fmt.Println("PriName is :", PriName)
	//查询语句拼接
	queryStr := PriName + " = ?"
	//另一种更新语句的写法
	//db.Model(args[0]).Updates(map[string]interface{}{})
	//where语句指定需要更新的主键，argsMap为指定更新的数据内容
	err := db.Table(str).Where(queryStr, strPri).Updates(argsMap).Error
	if err != nil {
		panic(err)
	}
	errQuery := db.Table(str).Where(queryStr, strPri).Find(&argsMap).Error
	if errQuery != nil {
		panic(errQuery)
	} else {
		fmt.Println("字段更新成功！")
	}
	fmt.Println("校验更新结果:")
	for k, v := range argsMap {
		fmt.Println(k, ":", v)
	}
	fmt.Println("---------Upgrade Func End---------")
}

func userDelete(args ...interface{}) {
	fmt.Println("---------Delete Func Start---------")
	var tableName []TabelInfor
	var lineName []string
	var PriName string
	var email string
	//表名获取
	str, ok := args[0].(string)
	if ok != true {
		panic(errors.New("error"))
	}
	//获取数据type
	//fmt.Printf("%T", str)
	sqlNum := "desc " + str
	//获取表字段名
	db.Raw(sqlNum).Scan(&tableName)
	for _, v := range tableName {
		lineName = append(lineName, v.Field)
	}
	//fmt.Println(lineName)
	//条件语句拼接及参数赋值
	strPri := args[1].(string) + " = ?"
	strArg := args[2]
	//查询主键
	sqlPri := "SELECT column_name FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` WHERE table_name=" + "'" + str + "'" +
		" AND CONSTRAINT_SCHEMA='mysql' AND constraint_name='PRIMARY'"
	db.Raw(sqlPri).Scan(&PriName)
	//fmt.Println(PriName)
	//删除语句
	err := db.Table(str).Where(strPri, strArg).Delete(&email).Error
	if err != nil {
		fmt.Println(err)
		panic(err)
	} else {
		fmt.Println("对应数据已删除")
	}
	//两种写法均可使用
	//err := db.Where("name = ?", "李四四四四").Delete(&Users{}).Error
	//db.Where("name = ?", "张三").Delete(&email)
	fmt.Println("---------Delete Func End---------")
}

//return 正确写法
//func userLogin(loginName string) []Users {
//	var users []Users
//	db.Where("name = ?", loginName).Find(&users)
//	return users
//}

func main() { // func main() 是程序开始执行的函数。main 函数是每一个可执行程序所必须包含的，一般来说都是在启动后第一个执行的函数（如果有 init() 函数则会先执行该函数）
	// 连接数据库
	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	//	用户登陆
	userLogin("李四四四四五五五")

	//	用户新增
	userInsert("users", "", "李四四四五个九", "男", 2)

	//  用户查询
	tableQuery("users", 29)

	//  用户更新
	userUpgrade("users", 29, "李四四四五个ßß", "男", 6)

	//  用户删除
	userDelete("users", "name", "李四四四四五五五")
}
