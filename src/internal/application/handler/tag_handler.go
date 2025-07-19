package handler

import (
	"context"
	"fmt"
	"poetry/pb/tag"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
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
	levelList := []int64{}
	if req.Limit == 0 {
		req.Limit = 20
	}
	for _, filter := range req.Filter {
		if filter.Name == "name" {
			nameList = filter.Value
		}
		if filter.Name == "category" {
			categoryList = filter.Value
		}
		if filter.Name == "level" {
			for _, value := range filter.Value {
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, capi_error.NewErr(capi_error.INVAILD_PARAM_CODE, fmt.Sprintf("tag value %s is not int", value))
				}
				levelList = append(levelList, int64(v))
			}
		}
	}
	count, tagList, err := TagHandler.tagService.DescribeTagInfo(ctx, nameList, categoryList, levelList, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	return &tag.DescribeTagInfoResponse{
		Total:       int32(count),
		TagInfoList: tagList,
	}, nil
}
