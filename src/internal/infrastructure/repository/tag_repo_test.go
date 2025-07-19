package repository

// import (
// 	"context"
// 	"poetry/src/internal/domain/entity"
// 	"poetry/src/internal/infrastructure"
// 	"testing"

// 	"gorm.io/gorm"
// 	"trpc.group/trpc-go/trpc-go"
// )

// func TestTagRepository_CreateTag(t *testing.T) {
// 	conf, err := trpc.LoadConfig("/root/code/poetry/src/cmd/trpc_go.yaml")
// 	if err != nil {
// 		panic(err)
// 	}
// 	trpc.Setup(conf)
// 	infrastructure.InitDB()
// 	type fields struct {
// 		db *gorm.DB
// 	}
// 	type args struct {
// 		ctx     context.Context
// 		tagList []*entity.Tag
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{
// 			name: "addTag",
// 			fields: fields{
// 				db: infrastructure.DB,
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 				tagList: []*entity.Tag{
// 					{
// 						Name:     "五代诗词",
// 						Category: "诗词",
// 						Level:    3,
// 					}, {
// 						Name:     "楚辞",
// 						Category: "诗词",
// 						Level:    4,
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &TagRepository{
// 				db: tt.fields.db,
// 			}
// 			if err := r.CreateTag(tt.args.ctx, tt.args.tagList); (err != nil) != tt.wantErr {
// 				t.Errorf("TagRepository.CreateTag() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
