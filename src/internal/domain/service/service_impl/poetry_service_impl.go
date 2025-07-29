package serviceimpl

import (
	"context"
	"poetry/pb/poetry"
	"poetry/src/internal/domain/entity"
	"poetry/src/internal/domain/repository"
	"poetry/src/internal/domain/service"
	"sort"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type PoetryServiceImpl struct {
	poetryRepo repository.PeotryRepo
}

type PoetryServiceImplOption struct {
	PoetryRepo repository.PeotryRepo
}

var _ service.PoetryService = &PoetryServiceImpl{}

var chineseNumMap = map[int]string{
	1:  "一",
	2:  "二",
	3:  "三",
	4:  "四",
	5:  "五",
	6:  "六",
	7:  "七",
	8:  "八",
	9:  "九",
	10: "十",
	11: "十一",
	12: "十二",
	13: "十三",
	14: "十四",
	15: "十五",
	16: "十六",
	17: "十七",
	18: "十八",
	19: "十九",
	20: "二十",
	21: "二十一",
	22: "二十二",
	23: "二十三",
	24: "二十四",
	25: "二十五",
	26: "二十六",
	27: "二十七",
	28: "二十八",
	29: "二十九",
	30: "三十",
	31: "三十一",
	32: "三十二",
	33: "三十三",
	34: "三十四",
}

func NewPoetryServiceImpl(option PoetryServiceImplOption) *PoetryServiceImpl {
	return &PoetryServiceImpl{
		poetryRepo: option.PoetryRepo,
	}
}

func (psi *PoetryServiceImpl) DescribePoetryInfo(ctx context.Context, title, author, paragraphs, dynasty, poetryType []string, tagId []int64, limit, offset int) (int64, []*poetry.PoetryInfo, error) {

	count, poetryDaoList, err := psi.poetryRepo.DescribePeotryInfo(ctx, title, author, paragraphs, dynasty, poetryType, tagId, limit, offset)
	poetryList := []*poetry.PoetryInfo{}
	if err != nil {
		return count, poetryList, err
	}
	c := collate.New(language.Chinese, collate.Numeric)
	for _, poetryInfo := range poetryDaoList {
		poetryList = append(poetryList, &poetry.PoetryInfo{
			Id:         poetryInfo.ID,
			Author:     poetryInfo.Author,
			Title:      poetryInfo.Title,
			Paragraphs: poetryInfo.Paragraphs,
			Rhythmic:   poetryInfo.Notes,
			Notes:      "",
			PoetryType: poetryInfo.Dynasty,
			Dynasty:    poetryInfo.Dynasty,
		})
	}
	sort.SliceStable(poetryList, func(i, j int) bool {
		if poetryList[i].Author != poetryList[j].Author {
			return poetryList[i].Author < poetryList[j].Author
		}
		return c.CompareString(poetryList[i].Title, poetryList[j].Title) < 0
	})

	begin := 1
	for i := 0; i < len(poetryList); i++ {
		if i+1 < len(poetryList) && poetryList[i].Author == poetryList[i+1].Author && poetryList[i].Title == poetryList[i+1].Title {
			poetryList[i].Title += "·" + chineseNumMap[begin]
			begin++
			continue
		}
		if begin > 1 {
			poetryList[i].Title += chineseNumMap[begin]
			begin = 1
		}
	}
	return count, poetryList, nil
}

func (psi *PoetryServiceImpl) CreatePoetryInfo(ctx context.Context, poetryList []*entity.Poetry) error {
	return psi.poetryRepo.CreatePoetryInfo(ctx, poetryList)
}
