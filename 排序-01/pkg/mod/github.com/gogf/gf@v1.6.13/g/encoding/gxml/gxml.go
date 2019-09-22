// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gxml provides accessing and converting for XML content.
//
// XML数据格式解析。
package gxml

import (
	"fmt"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/third/github.com/axgle/mahonia"
	"github.com/gogf/gf/third/github.com/clbanning/mxj"
	"strings"
)

// 将XML内容解析为map变量
func Decode(content []byte) (map[string]interface{}, error) {
	res, err := convert(content)
	if err != nil {
		return nil, err
	}
	return mxj.NewMapXml(res)
}

// 将map变量解析为XML格式内容
func Encode(v map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(v).Xml(rootTag...)
}

func EncodeWithIndent(v map[string]interface{}, rootTag ...string) ([]byte, error) {
	return mxj.Map(v).XmlIndent("", "\t", rootTag...)
}

// XML格式内容直接转换为JSON格式内容
func ToJson(content []byte) ([]byte, error) {
	res, err := convert(content)
	if err != nil {
		fmt.Println("convert error. ", err)
		return nil, err
	}

	mv, err := mxj.NewMapXml(res)
	if err == nil {
		return mv.Json()
	} else {
		return nil, err
	}
}

// XML字符集预处理
// @author wenzi1
// @date 20180604  修复并发安全问题,改为如果非UTF8字符集则先做字符集转换
func convert(xmlbyte []byte) (res []byte, err error) {
	patten := `<\?xml.*encoding\s*=\s*['|"](.*?)['|"].*\?>`
	matchStr, err := gregex.MatchString(patten, string(xmlbyte))
	if err != nil {
		return nil, err
	}

	xmlEncode := "UTF-8"
	if len(matchStr) == 2 {
		xmlEncode = matchStr[1]
	}

	s := mahonia.GetCharset(xmlEncode)
	if s == nil {
		return nil, fmt.Errorf("not support charset:%s\n", xmlEncode)
	}

	res, err = gregex.Replace(patten, []byte(""), []byte(xmlbyte))
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(s.Name, "UTF-8") {
		res = []byte(s.NewDecoder().ConvertString(string(res)))
	}

	return res, nil
}
