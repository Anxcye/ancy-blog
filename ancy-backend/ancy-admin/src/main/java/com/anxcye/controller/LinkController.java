package com.anxcye.controller;

import com.anxcye.domain.dto.ArticleDto;
import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.dto.LinkListDto;
import com.anxcye.domain.entity.Link;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.LinkVo;
import com.anxcye.service.ArticleService;
import com.anxcye.service.LinkService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/links")
public class LinkController {
    @Autowired
    private LinkService linkService;

    @Autowired
    private ArticleService articleService;

    @PreAuthorize("@ps.hasPermission('content:link:list')")
    @GetMapping("/list")
    public ResponseResult<List<Link>> linkList() {
        return ResponseResult.success(linkService.list());
    }

    @PreAuthorize("@ps.hasPermission('content:link:list')")
    @GetMapping("/page")
    public ResponseResult<PageResult> linkPage(@ParameterObject LinkListDto linkListDto) {
        return ResponseResult.success(linkService.pageList(linkListDto));
    }

    @PreAuthorize("@ps.hasPermission('content:link:add')")
    @PostMapping
    public ResponseResult<Long> linkAdd(@RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.addLink(linkDto));
    }

    @PreAuthorize("@ps.hasPermission('content:link:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> linkDelete(@PathVariable Long id) {
        return ResponseResult.success(linkService.deleteLink(id));
    }

    @PreAuthorize("@ps.hasPermission('content:link:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> linkUpdate(@PathVariable Long id, @RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.updateLink(id, linkDto));
    }

    @PreAuthorize("@ps.hasPermission('content:link:query')")
    @GetMapping("/{id}")
    public ResponseResult<LinkVo> linkGetById(@PathVariable Long id) {
        return ResponseResult.success(linkService.getLink(id));
    }

    @PreAuthorize("@ps.hasPermission('content:link:query')")
    @PutMapping("/article")
    public ResponseResult<Boolean> linkUpdateArticle(@RequestBody ArticleDto articleDto) {
        return ResponseResult.success(articleService.updateLinkArticle(articleDto));
    }

    @PreAuthorize("@ps.hasPermission('content:link:edit')")
    @GetMapping("/article")
    public ResponseResult<ArticleDetailVo> linkGetArticle() {
        return ResponseResult.success(articleService.getArticleLink());
    }

}
