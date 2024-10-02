package com.anxcye.controller;

import com.anxcye.domain.dto.TagDto;
import com.anxcye.domain.dto.TagListDto;
import com.anxcye.domain.entity.Tag;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.TagService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/tag")
public class TagController {

    @Autowired
    private TagService tagService;

    @GetMapping("/list")
    public ResponseResult<List<Tag>> tagList() {
        return ResponseResult.success(tagService.list());
    }

    @GetMapping("/page")
    public ResponseResult<PageResult> tagPage(@ParameterObject TagListDto tagListDto) {
        return ResponseResult.success(tagService.pageList(tagListDto));
    }

    @PostMapping
    public ResponseResult<Long> tagAdd(@RequestBody TagDto tagDto) {
        return ResponseResult.success(tagService.addTag(tagDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> tagDelete(@PathVariable Long id) {
        return ResponseResult.success(tagService.deleteTag(id));
    }

    @PutMapping("/{id}")
    public ResponseResult<Boolean> tagUpdate(@PathVariable Long id, @RequestBody TagDto tagDto) {
        return ResponseResult.success(tagService.updateTag(id, tagDto));
    }
}
