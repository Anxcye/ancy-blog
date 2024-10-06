package com.anxcye.controller;

import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ArticleCardVo;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.anxcye.service.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;


@RestController
@RequestMapping("/article")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @GetMapping("/page")
    public ResponseResult<PageResult> articlePage(Integer pageNum,
                               Integer pageSize,
            @RequestParam(required = false) Integer categoryId) {
        return ResponseResult.success(articleService.getList(pageNum, pageSize, categoryId));
    }

    @GetMapping("/front")
    public ResponseResult<List<ArticleCardVo>> articleFront() {
        return ResponseResult.success(articleService.getArticleFront());
    }
    
    @GetMapping("/recent")
    public ResponseResult<List<ArticleCardVo>> articleRecent() {
        return ResponseResult.success(articleService.recent());
    }

    @GetMapping("/hot")
    public ResponseResult<List<HotArticleVo>> articleHot() {
        return ResponseResult.success(articleService.hot());
    }

    @GetMapping("/{id}")
    public ResponseResult<ArticleDetailVo> articleGetById(@PathVariable Long id) {
        return ResponseResult.success(articleService.getArticleById(id));
    }

    @GetMapping("/home/{id}")
    public ResponseResult<ArticleDetailVo> articleGetFrontList(@PathVariable Long id) {
        return ResponseResult.success(articleService.getHomeById(id));
    }




}
