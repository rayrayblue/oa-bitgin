package usecase

import (
	"github.com/stretchr/testify/require"
	"oa-bitgin/pkg/domain"
	repo "oa-bitgin/pkg/repository"
	"testing"
	"time"
)

func Test_cashierUsecase_NewUser(t *testing.T) {
	type fields struct {
		userRepo domain.UserRepository
	}
	type args struct {
		name        string
		memberLevel int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "OK",
			fields: fields{
				userRepo: repo.NewUserRepository(),
			},
			args: args{
				name:        "testName",
				memberLevel: 0,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cashierUsecase{
				userRepo: tt.fields.userRepo,
			}
			got, err := c.NewUser(tt.args.name, tt.args.memberLevel)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cashierUsecase_BuyToken(t *testing.T) {
	type fields struct {
		userRepo     domain.UserRepository
		activityRepo domain.ActivityRepository
	}
	type args struct {
		userID int
		token  int64
	}
	tests := []struct {
		name       string
		buildStubs func(usecase domain.CashierUsecase)
		fields     fields
		args       args
		want       int
		wantErr    bool
	}{
		{
			name: "OKNormalMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 0)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    100,
			wantErr: false,
		},
		{
			name: "OKLevel1Member",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 1)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    95,
			wantErr: false,
		},
		{
			name: "OKLevel2Member",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 2)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    90,
			wantErr: false,
		},
		{
			name: "OKLevel3Member",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 3)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    85,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cashierUsecase{
				userRepo:     tt.fields.userRepo,
				activityRepo: tt.fields.activityRepo,
			}
			tt.buildStubs(c)
			got, err := c.BuyToken(tt.args.userID, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuyToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cashierUsecase_BuyTokenWithActivity(t *testing.T) {
	type fields struct {
		userRepo     domain.UserRepository
		activityRepo domain.ActivityRepository
		productRepo  domain.ProductRepository
	}
	type args struct {
		userID int
		token  int64
	}
	tests := []struct {
		name       string
		buildStubs func(usecase domain.CashierUsecase)
		fields     fields
		args       args
		want       int
		wantErr    bool
	}{
		{
			name: "OKNormalMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 0)
				_, _ = usecase.NewBuyTokenActivity(0, time.Now(), time.Now().Add(time.Hour*24*30), 50)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    50,
			wantErr: false,
		},
		{
			name: "OKLevel3Member",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 3)
				_, _ = usecase.NewBuyTokenActivity(3, time.Now(), time.Now().Add(time.Hour*24*30), 50)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    50,
			wantErr: false,
		},
		{
			name: "OKLevel3MemberMultipleActivity",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser", 3)
				_, _ = usecase.NewBuyTokenActivity(3, time.Now(), time.Now().Add(time.Hour*24*30), 90)
				_, _ = usecase.NewBuyTokenActivity(3, time.Now(), time.Now().Add(time.Hour*24*30), 85)
				_, _ = usecase.NewBuyTokenActivity(3, time.Now(), time.Now().Add(time.Hour*24*30), 80)
				_, _ = usecase.NewBuyTokenActivity(3, time.Now(), time.Now().Add(time.Hour*24*30), 75)
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID: 1,
				token:  100,
			},
			want:    75,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cashierUsecase{
				userRepo:     tt.fields.userRepo,
				activityRepo: tt.fields.activityRepo,
				productRepo:  tt.fields.productRepo,
			}
			tt.buildStubs(c)
			got, err := c.BuyTokenWithActivity(tt.args.userID, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuyTokenWithActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuyTokenWithActivity() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cashierUsecase_BuyProduct(t *testing.T) {
	type fields struct {
		userRepo     domain.UserRepository
		activityRepo domain.ActivityRepository
		productRepo  domain.ProductRepository
	}
	type args struct {
		userID    int
		productID int
	}
	tests := []struct {
		name       string
		buildStubs func(usecase domain.CashierUsecase)
		fields     fields
		args       args
		want       int
		wantErr    bool
		check      func(t *testing.T, usecase domain.CashierUsecase)
	}{
		{
			name: "OKNormalMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 0) // id = 1
				_, _ = usecase.NewUser("testUser2", 0) // id = 2
				_, _ = usecase.BuyToken(1, 1000)
				_, _ = usecase.BuyToken(2, 1000)
				_, _ = usecase.NewProduct("testProduct1", 100) // id = 1
				_, _ = usecase.NewProduct("testProduct2", 150) // id = 2
				_, _ = usecase.NewProduct("testProduct3", 200) // id = 3
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:    1,
				productID: 1,
			},
			want:    100,
			wantErr: false,
			check: func(t *testing.T, usecase domain.CashierUsecase) {
				remain, err := usecase.GetUserToken(1)
				require.NoError(t, err)
				require.Equal(t, 900, remain)
				require.Equal(t, int64(2000), usecase.GetTotalAmount())
			},
		},
		{
			name: "OKNormalMemberNotEnoughToken",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 0) // id = 1
				_, _ = usecase.NewUser("testUser2", 0) // id = 2
				_, _ = usecase.BuyToken(1, 10)
				_, _ = usecase.NewProduct("testProduct1", 100) // id = 1
				_, _ = usecase.NewProduct("testProduct2", 150) // id = 2
				_, _ = usecase.NewProduct("testProduct3", 200) // id = 3
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:    1,
				productID: 1,
			},
			want:    -1,
			wantErr: true,
			check: func(t *testing.T, usecase domain.CashierUsecase) {

			},
		},
		{
			name: "OKLevel1Member",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 1) // id = 1
				_, _ = usecase.NewUser("testUser2", 1) // id = 2
				_, _ = usecase.BuyToken(1, 1000)
				_, _ = usecase.NewProduct("testProduct1", 100) // id = 1
				_, _ = usecase.NewProduct("testProduct2", 150) // id = 2
				_, _ = usecase.NewProduct("testProduct3", 200) // id = 3
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:    1,
				productID: 1,
			},
			want:    100,
			wantErr: false,
			check: func(t *testing.T, usecase domain.CashierUsecase) {
				remain, err := usecase.GetUserToken(1)
				require.NoError(t, err)
				require.Equal(t, 900, remain)
				require.Equal(t, int64(950), usecase.GetTotalAmount())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cashierUsecase{
				userRepo:     tt.fields.userRepo,
				activityRepo: tt.fields.activityRepo,
				productRepo:  tt.fields.productRepo,
			}

			tt.buildStubs(c)
			got, err := c.BuyProduct(tt.args.userID, tt.args.productID)
			tt.check(t, c)
			if (err != nil) != tt.wantErr {
				t.Errorf("BuyProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuyProduct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cashierUsecase_BuyProductWithActivity(t *testing.T) {
	type fields struct {
		TotalAmount  int64
		userRepo     domain.UserRepository
		activityRepo domain.ActivityRepository
		productRepo  domain.ProductRepository
	}
	type args struct {
		userID     int
		productID  int
		activityID int
	}
	tests := []struct {
		name       string
		buildStubs func(usecase domain.CashierUsecase)
		fields     fields
		args       args
		want       int
		wantErr    bool
		check      func(t *testing.T, usecase domain.CashierUsecase)
	}{
		{
			name: "OKNormalMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 0) // id = 1
				_, _ = usecase.BuyToken(1, 1000)
				_ = usecase.AddPoint(1, 100)
				_, _ = usecase.NewProduct("testProduct1", 100)                                        // id = 1
				_, _ = usecase.NewProduct("testProduct2", 150)                                        // id = 2
				_, _ = usecase.NewProduct("testProduct3", 200)                                        // id = 3
				_, _ = usecase.NewBuyProductActivity(time.Now(), time.Now().Add(time.Hour*24*30), 90) // id = 1
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:     1,
				productID:  1,
				activityID: 1,
			},
			want:    90,
			wantErr: false,
			check: func(t *testing.T, usecase domain.CashierUsecase) {
				remainToken, err := usecase.GetUserToken(1)
				require.NoError(t, err)
				require.Equal(t, 910, remainToken)
				remainPoint, err := usecase.GetUserPoint(1)
				require.NoError(t, err)
				require.Equal(t, 90, remainPoint)
			},
		},
		{
			name: "OKLevelMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 1) // id = 1
				_, _ = usecase.BuyToken(1, 10000)
				_ = usecase.AddPoint(1, 1000)
				_, _ = usecase.NewProduct("testProduct1", 100)                                        // id = 1
				_, _ = usecase.NewBuyProductActivity(time.Now(), time.Now().Add(time.Hour*24*30), 90) // id = 1
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:     1,
				productID:  1,
				activityID: 1,
			},
			want:    90,
			wantErr: false,
			check: func(t *testing.T, usecase domain.CashierUsecase) {
				remainToken, err := usecase.GetUserToken(1)
				require.NoError(t, err)
				require.Equal(t, 9910, remainToken)
				remainPoint, err := usecase.GetUserPoint(1)
				require.NoError(t, err)
				require.Equal(t, 990, remainPoint)
				require.Equal(t, int64(9500), usecase.GetTotalAmount())
			},
		},
		{
			name: "OKLevelMember",
			buildStubs: func(usecase domain.CashierUsecase) {
				_, _ = usecase.NewUser("testUser1", 1) // id = 1
				_, _ = usecase.BuyToken(1, 10000)
				_ = usecase.AddPoint(1, 1000)
				_, _ = usecase.NewProduct("testProduct1", 1000)                                       // id = 1
				_, _ = usecase.NewBuyProductActivity(time.Now(), time.Now().Add(time.Hour*24*30), 80) // id = 1
			},
			fields: fields{
				userRepo:     repo.NewUserRepository(),
				activityRepo: repo.NewActivityRepository(),
				productRepo:  repo.NewProductRepository(),
			},
			args: args{
				userID:     1,
				productID:  1,
				activityID: 1,
			},
			want:    720,
			wantErr: false,
			check: func(t *testing.T, usecase domain.CashierUsecase) {
				require.Equal(t, int64(9500), usecase.GetTotalAmount())
				remainToken, err := usecase.GetUserToken(1)
				require.NoError(t, err)
				require.Equal(t, 9280, remainToken)
				remainPoint, err := usecase.GetUserPoint(1)
				require.NoError(t, err)
				require.Equal(t, 800, remainPoint)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cashierUsecase{
				TotalAmount:  tt.fields.TotalAmount,
				userRepo:     tt.fields.userRepo,
				activityRepo: tt.fields.activityRepo,
				productRepo:  tt.fields.productRepo,
			}

			tt.buildStubs(c)
			got, err := c.BuyProductWithActivity(tt.args.userID, tt.args.productID, tt.args.activityID)
			tt.check(t, c)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuyProductWithActivity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuyProductWithActivity() got = %v, want %v", got, tt.want)
			}
		})
	}
}
