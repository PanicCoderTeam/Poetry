package handler

import (
	"context"
	"poetry/pb/tag"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/log"
	"strconv"
)

type TagHandler struct {
	tagService service.TagService
}

type TagHandlerOption struct {
	TagService service.TagService
}

func NewTagHandler(option *TagHandlerOption) *TagHandler {
	return &TagHandler{
		tagService: option.TagService,
	}
}

func (TagHandler *TagHandler) DescribeTagInfo(ctx context.Context,
	req *tag.DescribeTagRequest) (*tag.DescribeTagInfoResponse, error) {
	log.DebugContextEx(ctx, "got hello request", req)
	nameList := []string{}
	categoryList := []string{}
	if req.Limit == 0 {
		req.Limit = 20
	}
	parentTagList := []int64{}
	for _, filter := range req.Filter {
		if filter.Name == "name" {
			nameList = filter.Value
		}
		if filter.Name == "category" {
			categoryList = filter.Value
		}
		if filter.Name == "parent-tag-id" {
			for _, value := range filter.Value {
				tagID, _ := strconv.ParseInt(value, 10, 64)
				parentTagList = append(parentTagList, tagID)
			}
		}
	}
	count, tagList, err := TagHandler.tagService.DescribeTagInfo(ctx, nameList, categoryList, parentTagList, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	return &tag.DescribeTagInfoResponse{
		Total:       int32(count),
		TagInfoList: tagList,
	}, nil
}
