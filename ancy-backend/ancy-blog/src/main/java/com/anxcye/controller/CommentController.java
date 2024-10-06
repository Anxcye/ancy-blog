package com.anxcye.controller;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.CommentDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CommentService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/comment")
public class CommentController {

    @Autowired
    private CommentService commentService;

    @GetMapping("/article/{articleId}")
    public ResponseResult<PageResult> commentByArticleId(@PathVariable Long articleId, Integer pageNum, Integer pageSize) {
        return ResponseResult.success(commentService.selectComment(
                SystemConstants.COMMENT_TYPE_ARTICLE,
                articleId,
                pageNum,
                pageSize));
    }

    @GetMapping("/{parentId}/children")
    public ResponseResult<PageResult> commentChildrenByParentId(@PathVariable Long parentId, Integer pageNum, Integer pageSize) {
        return ResponseResult.success(commentService.getChildren(parentId, pageNum, pageSize));
    }

    @PostMapping
    public ResponseResult<Long> commentAdd(@RequestBody CommentDto commentDto) {
        return ResponseResult.success(commentService.add(commentDto));
    }

    @GetMapping("/link")
    public ResponseResult<PageResult> commentLink(Integer pageNum, Integer pageSize) {
        return  ResponseResult.success(commentService.selectComment(
                SystemConstants.COMMENT_TYPE_LINK,
                null,
                pageNum,
                pageSize));
    }
}
