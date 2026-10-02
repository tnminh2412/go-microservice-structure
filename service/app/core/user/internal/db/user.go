package db

import (
	"context"
	"fmt"
	"proto-repo/generated/cmm/service/cmm/core/user/core_user"
	"time"

	"github.com/tnminh2412/go-microservice-structure/service/app/core/user/internal/pkg/error_helper"

	service_error "go-pkg/helper/error_helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username              string   `bson:"_id"`
	FirstName             string   `bson:"firstName"`
	LastName              string   `bson:"lastName"`
	DisplayName           string   `bson:"displayName"`
	Roles                 []string `bson:"roles"`
	Tasks                 []string `bson:"tasks"`
	Password              string   `bson:"password"`
	Status                string   `bson:"status"`
	Level                 int32    `bson:"level"`
	ReferrerID            string   `bson:"referrerID,omitempty"`
	Contact               bson.Raw `bson:"contact,omitempty"`
	IsSystem              bool     `bson:"isSystem"`
	FailLoginAttempt      int32    `bson:"failLoginAttempt"`
	CurrencyCodes         []string `bson:"currencyCodes"`
	IsEmailVerified       bool     `bson:"isEmailVerified"`
	IsPhoneNumberVerified bool     `bson:"isPhoneNumberVerified"`

	LastLoginCountry string    `bson:"lastLoginCountry"`
	LastLoginIP      string    `bson:"lastLoginIP"`
	LastLoginTime    time.Time `bson:"lastLoginTime"`

	ParentL1           string    `bson:"parent_l1,omitempty"`
	ParentL2           string    `bson:"parent_l2,omitempty"`
	ParentL3           string    `bson:"parent_l3,omitempty"`
	ParentL4           string    `bson:"parent_l4,omitempty"`
	ParentL5           string    `bson:"parent_l5,omitempty"`
	ParentL6           string    `bson:"parent_l6,omitempty"`
	ParentL7           string    `bson:"parent_l7,omitempty"`
	ParentL8           string    `bson:"parent_l8,omitempty"`
	ParentL9           string    `bson:"parent_l9,omitempty"`
	ParentL10          string    `bson:"parent_l10,omitempty"`
	ParentL11          string    `bson:"parent_l11,omitempty"`
	ParentL12          string    `bson:"parent_l12,omitempty"`
	ParentL13          string    `bson:"parent_l13,omitempty"`
	ParentL14          string    `bson:"parent_l14,omitempty"`
	ParentL15          string    `bson:"parent_l15,omitempty"`
	DirectParent       string    `bson:"directParent"`
	Domains            []*Domain `bson:"domains"`
	LandingRootDomains []string  `bson:"landingRootDomains,omitempty"`
	Images             []*Image  `bson:"images"`
	IsExternal         bool      `bson:"isExternal,omitempty"`

	PlayerIDPrefix string `bson:"playerIDPrefix,omitempty"`

	AccessLines []string `bson:"accessLines,omitempty"`

	IsRequireChangePassword bool              `bson:"isRequireChangePassword"`
	Data                    map[string]string `bson:"data,omitempty"`
	Email                   string            `bson:"email,omitempty"`
	Phone                   string            `bson:"phone,omitempty"`
	PhoneCode               string            `bson:"phoneCode,omitempty"`

	CreatedBy   string    `bson:"createdBy"`
	CreatedTime time.Time `bson:"createdTime"`
	UpdatedBy   string    `bson:"updatedBy"`
	UpdatedTime time.Time `bson:"updatedTime"`
}

type Mask struct {
	Fields    []string
	IsExclude bool
}

type Provider struct {
	Provider string `bson:"provider"`
	ID       string `bson:"id"`
}
type Domain struct {
	Domain    string      `bson:"domain"`
	IsTest    bool        `bson:"isTest"`
	Active    bool        `bson:"active"`
	IsLobby   bool        `bson:"isLobby"`
	Countries []string    `bson:"countries"`
	Image     string      `bson:"image"`
	Name      string      `bson:"name"`
	Language  string      `bson:"language"`
	Providers []*Provider `bson:"providers"`
	CreatedBy string      `bson:"createdBy"`
}

type Image struct {
	Image    string `bson:"image"`
	Name     string `bson:"name"`
	Language string `bson:"language"`
}

const userCollectionName = "user"

func CreateUser(ctx context.Context, mongodb *mongo.Database, user *User) error {
	_, err := mongodb.Collection(userCollectionName).InsertOne(ctx, user, &options.InsertOneOptions{
		BypassDocumentValidation: nil,
		Comment:                  nil,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return transformDuplicateErr(err)
		}
		return error_helper.New(service_error.Database, "", err)
	}
	return err
}

type UpdateUserParams struct {
	*User
	ByFirstName        bool
	ByLastName         bool
	ByRoles            bool
	ByTasks            bool
	ByStatus           bool
	ByPassword         bool
	ByContact          bool
	ByFailLoginAttempt bool
	ByCurrencyCodes    bool
	ByReferrerID       bool
	ByDisplayName      bool

	ByLastLoginTime    bool
	ByLastLoginCountry bool
	ByLastLoginIP      bool

	ByLandingRootDomains bool
	ByDomains            bool
	ByImages             bool

	ByIsRequireChangePassword bool
	ByData                    bool
	ByEmail                   bool
	ByPhone                   bool
	ByPhoneCode               bool
	ByIsPhoneNumberVerified   bool
	ByIsEmailVerified         bool
	ByPlayerIDPrefix          bool
}

func UpdateUser(ctx context.Context, mongodb *mongo.Database, params *UpdateUserParams) (int64, error) {
	set := bson.M{}

	if params.ByFirstName == true {
		set["firstName"] = params.FirstName
	}
	if params.ByLastName == true {
		set["lastName"] = params.LastName
	}
	if params.ByReferrerID == true {
		set["referrerID"] = params.ReferrerID
	}

	if params.ByLastLoginTime == true {
		set["lastLoginTime"] = params.LastLoginTime
	}
	if params.ByLastLoginCountry == true {
		set["lastLoginCountry"] = params.LastLoginCountry
	}
	if params.ByLastLoginIP == true {
		set["lastLoginIP"] = params.LastLoginIP
	}

	if params.ByRoles == true {
		set["roles"] = params.Roles
	}
	if params.ByTasks == true {
		set["tasks"] = params.Tasks
	}
	if params.ByStatus == true {
		set["status"] = params.Status
	}
	if params.ByPassword == true {
		set["password"] = params.Password
	}
	if params.ByContact == true {
		set["contact"] = params.Contact
	}
	if params.ByFailLoginAttempt == true {
		set["failLoginAttempt"] = params.FailLoginAttempt
	}
	if params.ByCurrencyCodes == true {
		set["currencyCodes"] = params.CurrencyCodes
	}
	if params.ByDisplayName == true {
		set["displayName"] = params.DisplayName
	}
	if params.ByIsRequireChangePassword == true {
		set["isRequireChangePassword"] = params.IsRequireChangePassword
	}
	if params.ByData {
		for k, v := range params.Data {
			set[fmt.Sprintf("data.%s", k)] = v
		}
	}
	if params.ByIsPhoneNumberVerified == true {
		set["isPhoneNumberVerified"] = params.IsPhoneNumberVerified
	}
	if params.ByIsEmailVerified == true {
		set["isEmailVerified"] = params.IsEmailVerified
	}
	if params.ByPlayerIDPrefix == true {
		set["playerIDPrefix"] = params.PlayerIDPrefix
	}

	if params.ByEmail == true {
		set["email"] = params.Email
	}

	if params.ByPhoneCode == true {
		set["phoneCode"] = params.PhoneCode
	}
	if params.ByPhone == true {
		set["phone"] = params.Phone
	}

	if params.ByLandingRootDomains == true {
		set["landingRootDomains"] = params.LandingRootDomains
	}

	if params.ByDomains == true {
		set["domains"] = params.Domains
	}

	if params.ByImages == true {
		set["images"] = params.Images
	}

	filter := bson.M{"_id": params.Username}

	set["updatedBy"] = params.UpdatedBy
	set["updatedTime"] = params.UpdatedTime

	update, err := mongodb.Collection(userCollectionName).UpdateOne(ctx, filter, bson.M{"$set": set})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return 0, transformDuplicateErr(err)
		}
		return 0, error_helper.New(service_error.Database, "", err)
	}

	return update.MatchedCount, err
}

func transformDuplicateErr(err error) error {
	mgErr, ok := isMongoServerErr(err)
	if !ok {
		return error_helper.New(service_error.Duplicate, "", err)
	}
	if isDuplicateUsername(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidUsernameDuplicated.String(), "", nil)
	}
	if isDuplicateEmail(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidEmailDuplicated.String(), "", nil)
	}
	if isDuplicatePhone(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidPhoneNumberDuplicated.String(), "", nil)
	}
	if isDuplicateReferrerID(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidReferrerIDDuplicated.String(), "", nil)
	}
	if isDuplicatePlayerIDPrefix(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidPlayerIDPrefixDuplicated.String(), "", nil)
	}
	if isDuplicateDomain(mgErr) {
		return error_helper.New(core_user.ErrorMessage_InvalidDomainDuplicated.String(), "", nil)
	}
	return error_helper.New(service_error.Duplicate, "", err)
}

func isMongoServerErr(err error) (mongo.ServerError, bool) {
	for ; err != nil; err = unwrap(err) {
		if e, ok := err.(mongo.ServerError); ok {
			return e, true
		}
	}
	return nil, false
}

func unwrap(err error) error {
	u, ok := err.(interface {
		Unwrap() error
	})
	if !ok {
		return nil
	}
	return u.Unwrap()
}

func isDuplicateUsername(err mongo.ServerError) bool {
	return err.HasErrorMessage("index: _id_ dup key")
}
func isDuplicateEmail(err mongo.ServerError) bool {
	return err.HasErrorMessage("index: user_email_unique dup key")
}
func isDuplicatePhone(err mongo.ServerError) bool {
	return err.HasErrorMessage("index: user_phone_unique dup key")
}
func isDuplicateReferrerID(err mongo.ServerError) bool {
	return err.HasErrorMessage("index: user_referrer_unique dup key")
}
func isDuplicatePlayerIDPrefix(err mongo.ServerError) bool {
	return err.HasErrorMessage("playerIDPrefix")
}
func isDuplicateDomain(err mongo.ServerError) bool {
	return err.HasErrorMessage("domains.domain")
}

type Parent struct {
	Username string
	Level    int32
}

type LevelParams struct {
	ByEq bool
	ByGt bool
	Eq   int32
	Gt   int32
}
type ListUserParams struct {
	ByLevel         bool
	BySkip          bool
	ByLimit         bool
	ByIDs           bool
	ByStatuses      bool
	ByDirectParent  bool
	ByExcludeSystem bool
	ByRoles         bool
	ByExcludeRoles  bool
	ByIsSystem      bool
	ByExcludeHidden bool

	Skip          int64
	Limit         int64
	IDs           []string
	Statuses      []string
	DirectParent  string
	ExcludeSystem bool
	Mask          *Mask
	Level         *LevelParams
	Roles         []string
	ExcludeRoles  []string
	IsSystem      bool
	OrderBy       string
	IsAsc         bool
	ExcludeHidden bool
}

func ListUser(ctx context.Context, mongodb *mongo.Database, params ListUserParams) ([]*User, error) {

	filter := bson.M{}
	if params.ByIDs {
		filter["_id"] = bson.M{"$in": params.IDs}
	}
	if params.ByStatuses {
		filter["status"] = bson.M{"$in": params.Statuses}
	}

	if params.ByDirectParent {
		filter["directParent"] = params.DirectParent
	}
	if params.ByIsSystem {
		filter["isSystem"] = params.IsSystem
	}
	if params.ByExcludeSystem && params.ExcludeSystem {
		filter["isSystem"] = bson.M{
			"$ne": true,
		}
	}
	if params.ByExcludeHidden && params.ExcludeHidden {
		filter["isHidden"] = bson.M{
			"$ne": true,
		}
	}
	if params.ByLevel {
		level := bson.M{}
		if params.Level.ByEq {
			level["$eq"] = params.Level.Eq
		}
		if params.Level.ByGt {
			level["$gt"] = params.Level.Gt
		}
		filter["level"] = level
	}
	if params.ByRoles || params.ByExcludeRoles {
		filterRoles := bson.M{}
		if params.ByRoles {
			filterRoles["$in"] = params.Roles
		}
		if params.ByExcludeRoles {
			filterRoles["$nin"] = params.ExcludeRoles
		}
		filter["roles"] = filterRoles
	}

	opts := &options.FindOptions{}
	if params.ByLimit == true && params.BySkip {
		opts.SetSkip(params.Skip)
		opts.SetLimit(params.Limit)
	}

	sorting := bson.M{}
	sortValue := 1
	if params.IsAsc == false {
		sortValue = -1
	}
	if len(params.OrderBy) > 0 {
		sorting[params.OrderBy] = sortValue
		opts.SetSort(sorting)
	}

	if params.Mask != nil && len(params.Mask.Fields) != 0 {
		mask := params.Mask
		project := bson.M{}
		maskVal := 1
		if mask.IsExclude {
			maskVal = 0
		}
		for _, m := range mask.Fields {
			project[m] = maskVal
		}
		opts.SetProjection(project)
	}

	cursor, err := mongodb.Collection(userCollectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var reply []*User

	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		} else {
			reply = append(reply, &user)
		}
	}

	return reply, nil
}

type GetUserParams struct {
	ByID         bool
	ByReferrerID bool
	ByEmail      bool
	ByDomain     bool
	ByPhone      bool
	ByPhoneCode  bool

	ID         string
	ReferrerID string
	Email      string
	Domain     string
	Phone      string
	PhoneCode  string
	Mask       *Mask
}

func GetUser(ctx context.Context, mongodb *mongo.Database, params GetUserParams) (*User, error) {
	var user User

	filter := bson.M{}
	if params.ByID == true {
		filter = bson.M{"_id": params.ID}
	}
	if params.ByReferrerID == true {
		filter = bson.M{"referrerID": params.ReferrerID}
	}
	if params.ByEmail == true {
		filter = bson.M{"email": params.Email}
	}
	if params.ByPhone == true {
		filter["phone"] = params.Phone

		if params.ByPhoneCode == true {
			filter["phoneCode"] = params.PhoneCode
		}
	}

	if params.ByDomain == true {
		filter = bson.M{"domains.domain": params.Domain}
	}

	opts := options.FindOne()
	if params.Mask != nil {
		mask := params.Mask
		maskVal := 1
		project := bson.M{}

		if mask.IsExclude {
			maskVal = 0
		}
		for _, m := range mask.Fields {
			project[m] = maskVal
		}
		opts.SetProjection(project)
	}

	err := mongodb.Collection(userCollectionName).FindOne(ctx, filter, opts).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
