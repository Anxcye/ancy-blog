package com.anxcye.controller;

import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.dto.LinkListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.LinkService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/links")
public class LinkController {
    @Autowired
    private LinkService linkService;

    @GetMapping("/list")
    public ResponseResult<?> list() {
        return ResponseResult.success(linkService.list());
    }

    @GetMapping("/page")
    public ResponseResult<?> pageList(@ParameterObject LinkListDto linkListDto) {
        return ResponseResult.success(linkService.pageList(linkListDto));
    }

    @PostMapping
    public ResponseResult<?> addLink(@RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.addLink(linkDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> deleteLink(@PathVariable Long id) {
        return ResponseResult.success(linkService.deleteLink(id));
    }

    @PutMapping("/{id}")
    public ResponseResult<?> updateLink(@PathVariable Long id, @RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.updateLink(id, linkDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<?> getLink(@PathVariable Long id) {
        return ResponseResult.success(linkService.getLink(id));
    }
}
