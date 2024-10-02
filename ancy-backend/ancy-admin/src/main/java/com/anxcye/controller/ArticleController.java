package com.anxcye.controller;

import com.anxcye.domain.dto.ArticleDto;
import com.anxcye.domain.dto.ArticleListDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.service.ArticleService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/article")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @PostMapping
    public ResponseResult<Boolean> articleAdd(@RequestBody ArticleDto articleDto) {
        return ResponseResult.success(articleService.addArticle(articleDto));
    }

    @GetMapping("/page")
    public ResponseResult<PageResult> articlePage(@ParameterObject ArticleListDto articleListDto) {
        return ResponseResult.success(articleService.pageList(articleListDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<ArticleDetailVo> articleGetById(@PathVariable Long id) {
        return ResponseResult.success(articleService.getArticleById(id));
    }


    @PutMapping("/{id}")
    public ResponseResult<Boolean> articleUpdate(@PathVariable Long id, @RequestBody ArticleDto addArticleDto) {
        return ResponseResult.success(articleService.updateArticleById(id, addArticleDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> articleDelete(@PathVariable Long id) {
        return ResponseResult.success(articleService.deleteArticleById(id));
    }


}
