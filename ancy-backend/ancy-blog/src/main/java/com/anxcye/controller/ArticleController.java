package com.anxcye.controller;

import com.anxcye.domin.ResponseResult;
import com.anxcye.domin.entity.Article;
import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/articles")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @GetMapping("/list")
    public List<Article> list() {
        return articleService.list();
    }

    @GetMapping("/hot")
    public ResponseResult hot() {
        return ResponseResult.okResult(articleService.hot());
    }
}
