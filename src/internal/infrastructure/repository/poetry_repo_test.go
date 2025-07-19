package repository

import (
	"context"
	"fmt"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/infrastructure"
	"testing"

	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-go"
)

// func TestPoetryRepository_DescribePeotryInfo(t *testing.T) {
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
// 		ctx        context.Context
// 		title      string
// 		author     string
// 		paragraphs string
// 		dynasty    string
// 		limit      int
// 		offset     int
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    int64
// 		want1   []*entity.Poetry
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 		{name: "abs",
// 			fields: fields{
// 				db: infrastructure.DB,
// 			},
// 			args: args{
// 				ctx: context.Background(),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &PoetryRepository{
// 				db: tt.fields.db,
// 			}
// 			offset := 0
// 			limit := 1000

//				for {
//					got, got1, err := r.DescribePeotryInfo(tt.args.ctx, "", "", "", "", limit, offset)
//					offset += limit
//					if offset >= int(got) {
//						break
//					}
//					// fmt.Printf("%+v %+v %+v\n", got, got1, err)
//					if err != nil {
//						t.Errorf("DescribePeotryInfo() error = %v, wantErr %v", err, tt.wantErr)
//						return
//					}
//					tryList := []*entity.Poetry{}
//					for _, poetry := range got1 {
//						if len(poetry.TitleTradition) > 0 {
//							continue
//						}
//						poetry.AuthorTradition = utils.ConvertChinsesSimplified2T(poetry.Author)
//						poetry.Author = utils.ConvertChinsesTraditional2S(poetry.Author)
//						poetry.TitleTradition = utils.ConvertChinsesSimplified2T(poetry.Title)
//						poetry.Title = utils.ConvertChinsesTraditional2S(poetry.Title)
//						tryList = append(tryList, poetry)
//						fmt.Printf("title:%+v\n", poetry.Title)
//						fmt.Printf("title:%+v\n", poetry.TitleTradition)
//						// fmt.Printf("%+v\n", poetry.ParagraphsTradition)
//						// fmt.Printf("%+v\n", poetry.Paragraphs)
//					}
//					if len(tryList) > 0 {
//						// r.UpdatePoetryInfo(tt.args.ctx, tryList)
//					}
//				}
//			})
//		}
//	}
func TestPoetryRepository_DescribePeotryInfo1(t *testing.T) {
	conf, err := trpc.LoadConfig("/root/code/poetry/src/cmd/trpc_go.yaml")
	if err != nil {
		panic(err)
	}
	trpc.Setup(conf)
	infrastructure.InitDB()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx        context.Context
		title      string
		author     string
		paragraphs string
		dynasty    string
		limit      int
		offset     int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		want1   []*entity.Poetry
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "abs",
			fields: fields{
				db: infrastructure.DB,
			},
			args: args{
				ctx: context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &PoetryRepository{
				db: tt.fields.db,
			}
			rt := &PoetryTagRepository{
				db: tt.fields.db,
			}
			got, got1, err := r.DescribePeotryInfo(tt.args.ctx, []string{}, []string{}, []string{}, []string{}, []string{"楚辞"}, []int64{}, 10000, 0)
			if got == 0 {
				fmt.Printf("%+v %+v %+v\n", got, got1, err)
			}
			poetryTagInfoList := []*entity.PoetryTag{}
			for _, poetry := range got1 {
				fmt.Printf("%v\n", poetry)
				poetryTagInfo := &entity.PoetryTag{
					PoetryID: poetry.ID,
					Tag:      "楚辞",
					Category: "诗词",
					TagID:    7,
				}
				poetryTagInfoList = append(poetryTagInfoList, poetryTagInfo)
				fmt.Printf("%+v\n", poetryTagInfo)

			}
			if len(poetryTagInfoList) > 0 {
				rt.CreatePoetryTagInfo(tt.args.ctx, poetryTagInfoList)
				if err != nil {
					t.Errorf("DescribePeotryInfo() error = %v, wantErr %v", err, tt.wantErr)
					// return
				}
			}
		})
	}
}
