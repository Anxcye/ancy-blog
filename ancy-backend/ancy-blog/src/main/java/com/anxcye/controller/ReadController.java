package com.anxcye.controller;

import com.anxcye.domain.dto.ReadPageDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ReadService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/read")
public class ReadController {

    @Autowired
    private ReadService readService;

    @GetMapping("/page")
    public ResponseResult<PageResult> readPage(@ParameterObject ReadPageDto readPageDto) {
        return ResponseResult.success(readService.getReadPage(readPageDto));
    }
}
