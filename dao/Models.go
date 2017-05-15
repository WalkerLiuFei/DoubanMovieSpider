package dao


type Movie struct {
	ID          string `gorm :not null;Auto Increment`
	Name        string `gorm :not null;primary key;` //电影名字
	Link        string `gorm :not null`             //链接
	PostImgLink string `gorm :not null`             //海报地址
	Score       string `gorm :not null`             //评分
	JudgeNumber string `gorm :not null`             // 评分人数
	Director    string `gorm :not null`             //导演
	LeadingStar string `gorm :not null`             //主演
	MovieType   string `gorm :not null`             //电影类型
	Country     string `gorm :not null`             //国家
	PostYear    string `gorm :not null`             //发布时间
}
