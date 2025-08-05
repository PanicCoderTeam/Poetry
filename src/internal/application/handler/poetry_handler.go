package handler

import (
	"context"
	"fmt"
	"poetry/pb/poetry"
	"poetry/src/internal/domain/service"
	"poetry/src/pkg/log"
	"poetry/src/pkg/trpc/codec/capi_error"
	"strconv"
)

type PoetryHandler struct {
	peotryService service.PoetryService
}

type PoetryHandlerOption struct {
	PeotryService service.PoetryService
}

func NewPoetryHandler(option *PoetryHandlerOption) *PoetryHandler {
	return &PoetryHandler{
		peotryService: option.PeotryService,
	}
}

func (poetryHandler *PoetryHandler) DescribePoetryInfo(ctx context.Context,
	req *poetry.DescribePoetryInfoRequest) (*poetry.DescribePoetryInfoResponse, error) {
	log.DebugContextEx(ctx, "got hello request", req)
	titleList := []string{}
	authorList := []string{}
	paragraphsList := []string{}
	dynestyList := []string{}
	tagIdList := []int64{}
	poetryTypeList := []string{}
	for _, filter := range req.Filter {
		if filter.Name == "title" {
			titleList = filter.Value
		}
		if filter.Name == "author" {
			authorList = filter.Value
		}
		if filter.Name == "paragraphs" {
			paragraphsList = filter.Value
		}
		if filter.Name == "dynasty" {
			dynestyList = filter.Value
		}
		if filter.Name == "poetry-type" {
			poetryTypeList = filter.Value
		}
		if filter.Name == "tag-id" {
			for _, value := range filter.Value {
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, capi_error.NewErr(capi_error.INVAILD_PARAM_CODE, fmt.Sprintf("tag value %s is not int", value))
				}
				tagIdList = append(tagIdList, int64(v))
			}
		}
	}
	count, poetryList, err := poetryHandler.peotryService.DescribePoetryInfo(ctx, titleList, authorList, paragraphsList, dynestyList, poetryTypeList, tagIdList, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}
	return &poetry.DescribePoetryInfoResponse{
		Total:          int32(count),
		PoetryInfoList: poetryList,
	}, nil
}
