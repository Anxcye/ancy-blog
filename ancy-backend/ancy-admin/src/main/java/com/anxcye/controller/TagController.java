package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.TagService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/tag")
public class TagController {

    @Autowired
    private TagService tagService;

    @GetMapping("/list")
    public ResponseResult<?> list(){
        return  ResponseResult.success(tagService.list());
    }
}
