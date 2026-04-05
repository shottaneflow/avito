package models

type Statistic struct {
	Likes     int64 `json:"likes"`
	ViewCount int64 `json:"viewCount"`
	Contacts  int64 `json:"contacts"`
}

type CreateAdvertisementRequest struct {
	SellerID  int64     `json:"sellerID"`
	Name      string    `json:"name"`
	Price     int64     `json:"price"`
	Statistic Statistic `json:"statistics"`
}
