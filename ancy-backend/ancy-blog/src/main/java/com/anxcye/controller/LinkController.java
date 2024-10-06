package com.anxcye.controller;

import com.anxcye.domain.dto.LinkDto;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.LinkVo;
import com.anxcye.service.ArticleService;
import com.anxcye.service.LinkService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/link")
public class LinkController {

    @Autowired
    private LinkService linkService;

    @Autowired
    private ArticleService articleService;

    @GetMapping("/list")
    public ResponseResult<List<LinkVo>> linkList() {
        return ResponseResult.success(linkService.getApprovedLinks());
    }

    @GetMapping("/article")
    public ResponseResult<ArticleDetailVo> linkGetArticle() {
        return ResponseResult.success(articleService.getArticleLink());
    }


    @PostMapping
    public ResponseResult<Long> linkAdd(@RequestBody LinkDto linkDto) {
        return ResponseResult.success(linkService.addLinkBlog(linkDto));
    }
}
