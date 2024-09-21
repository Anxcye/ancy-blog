package com.anxcye.controller;

import com.anxcye.domain.dto.AddArticleDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/articles")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @PostMapping
    public ResponseResult<?> addArticle(@RequestBody AddArticleDto articleDto){
        return ResponseResult.success(articleService.addArticle(articleDto));
    }

}
