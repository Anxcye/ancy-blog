package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CommentService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/comment")
public class CommentController {

    @Autowired
    private CommentService commentService;

    @GetMapping("/article/{articleId}")
    public ResponseResult<?> selectCommentByArticleId(@PathVariable Long articleId, Integer pageNum, Integer pageSize) {
         return ResponseResult.success(commentService.selectCommentByArticleId(articleId, pageNum, pageSize));
    }

    @GetMapping("/{parentId}/children")
    public ResponseResult<?> selectChildren(@PathVariable Long parentId, Integer pageNum, Integer pageSize) {
        return  ResponseResult.success(commentService.getChildren(parentId, pageNum, pageSize));
    }
}
