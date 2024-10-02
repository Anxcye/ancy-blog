package com.anxcye.controller;

import com.anxcye.domain.dto.ArticleDto;
import com.anxcye.domain.dto.ArticleListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.service.ArticleService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/article")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @PreAuthorize("@ps.hasPermission('content:article:add')")
    @PostMapping
    public ResponseResult<Boolean> articleAdd(@RequestBody ArticleDto articleDto) {
        return ResponseResult.success(articleService.addArticle(articleDto));
    }

    @PreAuthorize("@ps.hasPermission('content:article:list')")
    @GetMapping("/page")
    public ResponseResult<PageResult> articlePage(@ParameterObject ArticleListDto articleListDto) {
        return ResponseResult.success(articleService.pageList(articleListDto));
    }

    @PreAuthorize("@ps.hasPermission('content:article:query')")
    @GetMapping("/{id}")
    public ResponseResult<ArticleDetailVo> articleGetById(@PathVariable Long id) {
        return ResponseResult.success(articleService.getArticleById(id));
    }


    @PreAuthorize("@ps.hasPermission('content:article:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> articleUpdate(@PathVariable Long id, @RequestBody ArticleDto addArticleDto) {
        return ResponseResult.success(articleService.updateArticleById(id, addArticleDto));
    }

    @PreAuthorize("@ps.hasPermission('content:article:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> articleDelete(@PathVariable Long id) {
        return ResponseResult.success(articleService.deleteArticleById(id));
    }
}
