package com.anxcye.controller;

import com.anxcye.domain.dto.TagDto;
import com.anxcye.domain.dto.TagListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.TagService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/tag")
public class TagController {

    @Autowired
    private TagService tagService;

    @GetMapping("/list")

    public ResponseResult<?> list() {
        return ResponseResult.success(tagService.list());
    }

    @GetMapping("/pageList")
    public ResponseResult<?> list(@ParameterObject TagListDto tagListDto) {
        return ResponseResult.success(tagService.pageList(tagListDto));
    }

    @PostMapping
    public ResponseResult<?> addTag(@RequestBody TagDto tagDto) {
        return ResponseResult.success(tagService.addTag(tagDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> deleteTag(@PathVariable Long id) {
        return ResponseResult.success(tagService.deleteTag(id));
    }

    @PutMapping("/{id}")
    public ResponseResult<?> updateTag(@PathVariable Long id, @RequestBody TagDto tagDto) {
        return ResponseResult.success(tagService.updateTag(id, tagDto));
    }
}