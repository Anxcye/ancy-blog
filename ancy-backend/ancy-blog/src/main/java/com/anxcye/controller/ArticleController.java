package com.anxcye.controller;

import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;


@RestController
@RequestMapping("/articles")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @GetMapping("/list")
    public ResponseResult<PageResult> list(Integer pageNum,
                               Integer pageSize,
                               @RequestParam(required = false) Integer categoryId) {
        return ResponseResult.success(articleService.getList(pageNum, pageSize, categoryId));
    }

    @GetMapping("/hot")
    public ResponseResult<?> hot() {
        return ResponseResult.success(articleService.hot());
    }

    @GetMapping("/{id}")
    public ResponseResult<?> getArticleDetail(@PathVariable Long id) {
        return ResponseResult.success(articleService.getArticleById(id));
    }
}
