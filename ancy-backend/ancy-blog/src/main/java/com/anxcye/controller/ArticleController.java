package com.anxcye.controller;

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
    public ResponseResult list(Integer pageNum,
                               Integer pageSize,
                               @RequestParam(required = false) Integer categoryId) {
        return ResponseResult.okResult(articleService.getList(pageNum, pageSize, categoryId));
    }

    @GetMapping("/hot")
    public ResponseResult hot() {
        return ResponseResult.okResult(articleService.hot());
    }

    @GetMapping("/{id}")
    public ResponseResult getArticleDetail(@PathVariable Long id) {
        return ResponseResult.okResult(articleService.getArticleById(id));
    }
}
