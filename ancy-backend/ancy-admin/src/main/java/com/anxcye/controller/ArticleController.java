package com.anxcye.controller;

import com.anxcye.domain.dto.AddArticleDto;
import com.anxcye.domain.dto.ArticleListDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ArticleService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/articles")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @PostMapping
    public ResponseResult<?> addArticle(@RequestBody AddArticleDto articleDto) {
        return ResponseResult.success(articleService.addArticle(articleDto));
    }

    @GetMapping("/page")
    public ResponseResult<?> pageList(@ParameterObject ArticleListDto articleListDto) {
        return ResponseResult.success(articleService.pageList(articleListDto));
    }

    @GetMapping("/{id}")
    public ResponseResult<?> getArticleById(@PathVariable Long id) {
        return ResponseResult.success(articleService.getArticleById(id));
    }


    @PutMapping("/{id}")
    public ResponseResult<?> updateArticle(@PathVariable Long id, @RequestBody AddArticleDto addArticleDto) {
        return ResponseResult.success(articleService.updateArticleById(id, addArticleDto));
    }

    @DeleteMapping("/{id}")
    public ResponseResult<?> deleteArticle(@PathVariable Long id) {
        return ResponseResult.success(articleService.deleteArticleById(id));
    }


}
