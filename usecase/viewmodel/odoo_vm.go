package viewmodel

// TravelPackageVM ...
type TravelPackageVM struct {
	LastUpdate  string        `xmlrpc:"__last_update,omptempty"`
	ArrivalDate string        `xmlrpc:"arrival_date,omptempty"`
	CreateDate  string        `xmlrpc:"create_date,omptempty"`
	DisplayName string        `xmlrpc:"display_name,omptempty"`
	ID          int64         `xmlrpc:"id,omptempty"`
	IsActive    bool          `xmlrpc:"is_active,omptempty"`
	JamaahList  []interface{} `xmlrpc:"jamaah_list,omptempty"`
	GuideList   []interface{} `xmlrpc:"guide_ids,omptempty"`
	GuideListID []int64       `xmlrpc:"guide_list_id,omptempty"`
	UserList    []int64       `xmlrpc:"user_list,omptempty"`
	Name        string        `xmlrpc:"name,omptempty"`
}

// PartnerVM ...
type PartnerVM struct {
	LastUpdate            string        `xmlrpc:"__last_update,omptempty"`
	Age                   int64         `xmlrpc:"age,omptempty"`
	BirthOfDate           string        `xmlrpc:"birth_of_date,omptempty"`
	City                  string        `xmlrpc:"city,omptempty"`
	CommercialCompanyName string        `xmlrpc:"commercial_company_name,omptempty"`
	CompanyName           string        `xmlrpc:"company_name,omptempty"`
	CompanyType           string        `xmlrpc:"company_type,omptempty"`
	ContactAddress        string        `xmlrpc:"contact_address,omptempty"`
	CreateDate            string        `xmlrpc:"create_date,omptempty"`
	Credit                float64       `xmlrpc:"credit,omptempty"`
	CreditLimit           float64       `xmlrpc:"credit_limit,omptempty"`
	Date                  string        `xmlrpc:"date,omptempty"`
	Debit                 float64       `xmlrpc:"debit,omptempty"`
	DebitLimit            float64       `xmlrpc:"debit_limit,omptempty"`
	DisplayName           string        `xmlrpc:"display_name,omptempty"`
	Email                 string        `xmlrpc:"email,omptempty"`
	Employee              bool          `xmlrpc:"employee,omptempty"`
	Gender                string        `xmlrpc:"gender,omptempty"`
	Image                 string        `xmlrpc:"image,omptempty"`
	IsCompany             bool          `xmlrpc:"is_company,omptempty"`
	IsJamaah              bool          `xmlrpc:"is_jamaah,omptempty"`
	IsMultiJamaah         bool          `xmlrpc:"is_multi_jamaah,omptempty"`
	IsTourLeader          bool          `xmlrpc:"is_tour_leader,omptempty"`
	MaritalStatus         string        `xmlrpc:"martial_status,omptempty"`
	Mobile                string        `xmlrpc:"mobile,omptempty"`
	Name                  string        `xmlrpc:"name,omptempty"`
	PackageList           []interface{} `xmlrpc:"package_id,omptempty"`
	PackageListID         []int64       `xmlrpc:"package)list_id,omptempty"`
}

// GuideVM ...
type GuideVM struct {
	LastUpdate     string        `xmlrpc:"__last_update,omptempty"`
	CreateDate     string        `xmlrpc:"create_date,omptempty"`
	DisplayName    string        `xmlrpc:"display_name,omptempty"`
	GuideContactID int64         `xmlrpc:"guide_contact_id,omptempty"`
	ID             int64         `xmlrpc:"id,omptempty"`
	PartnerID      []interface{} `xmlrpc:"partner_id,omptempty"`
	PriceSubtotal  float64       `xmlrpc:"price_subtotal,omptempty"`
	PriceUnit      float64       `xmlrpc:"price_unit,omptempty"`
}
