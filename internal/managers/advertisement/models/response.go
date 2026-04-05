package models

type Advertisement struct {
	ID        string    `json:"id"`
	SellerID  int64     `json:"sellerId"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Statistic Statistic `json:"statistics"`
	CreatedAt string    `json:"createdAt"`
}

type StatisticResponse struct {
	Likes     int64 `json:"likes"`
	ViewCount int64 `json:"viewCount"`
	Contacts  int64 `json:"contacts"`
}
