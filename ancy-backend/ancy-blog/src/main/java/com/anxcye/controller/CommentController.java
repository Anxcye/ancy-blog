package com.anxcye.controller;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.CommentDto;
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
    public ResponseResult<?> selectCommentByArticleId(@PathVariable Long articleId, Integer pageNum, Integer pageSize) {
        return ResponseResult.success(commentService.selectComment(
                SystemConstants.COMMENT_TYPE_ARTICLE,
                articleId,
                pageNum,
                pageSize));
    }

    @GetMapping("/{parentId}/children")
    public ResponseResult<?> selectChildren(@PathVariable Long parentId, Integer pageNum, Integer pageSize) {
        return ResponseResult.success(commentService.getChildren(parentId, pageNum, pageSize));
    }

    @PostMapping
    public ResponseResult<?> add(@RequestBody CommentDto commentDto) {
        commentService.add(commentDto);
        return ResponseResult.success();
    }

    @GetMapping("/link")
    public ResponseResult<?> selectLinkComment( Integer pageNum, Integer pageSize) {
        return  ResponseResult.success(commentService.selectComment(
                SystemConstants.COMMENT_TYPE_LINK,
                null,
                pageNum,
                pageSize));
    }
}
