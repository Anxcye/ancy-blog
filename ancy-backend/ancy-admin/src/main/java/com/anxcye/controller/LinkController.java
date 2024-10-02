package com.anxcye.controller;

import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.dto.LinkListDto;
import com.anxcye.domain.entity.Link;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.LinkVo;
import com.anxcye.service.LinkService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/links")
public class LinkController {
    @Autowired
    private LinkService linkService;

    @GetMapping("/list")
    public ResponseResult<List<Link>> linkList() {
        return ResponseResult.success(linkService.list());
    }

    @GetMapping("/page")
    public ResponseResult<PageResult> linkPage(@ParameterObject LinkListDto linkListDto) {
        return ResponseResult.success(linkService.pageList(linkListDto));
    }

    @PostMapping
    public ResponseResult<Long> linkAdd(@RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.addLink(linkDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> linkDelete(@PathVariable Long id) {
        return ResponseResult.success(linkService.deleteLink(id));
    }

    @PutMapping("/{id}")
    public ResponseResult<Boolean> linkUpdate(@PathVariable Long id, @RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.updateLink(id, linkDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<LinkVo> linkGetById(@PathVariable Long id) {
        return ResponseResult.success(linkService.getLink(id));
    }
}
