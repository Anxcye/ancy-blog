package com.anxcye.controller;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.CommentDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CommentService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
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
        return ResponseResult.success(commentService.addComment(commentDto));
    }

    @GetMapping("/note/{id}")
    public ResponseResult<PageResult> commentNote(@PathVariable Long id, Integer pageNum, Integer pageSize) {
        return  ResponseResult.success(commentService.selectComment(
                SystemConstants.COMMENT_TYPE_NOTE,
                id,
                pageNum,
                pageSize));
    }

    @PreAuthorize("@ps.hasPermission('content:comment:edit')")
    @PutMapping("/admin/{id}")
    public ResponseResult<Boolean> commentUpdate(@PathVariable Long id, @RequestBody CommentDto commentDto) {
        return ResponseResult.success(commentService.updateComment(id, commentDto));
    }

    @PutMapping("/{id}/like")
    public ResponseResult<Boolean> commentLike(@PathVariable Long id, @RequestParam Boolean increase) {
        return ResponseResult.success(commentService.updateCommentLike(id, increase));
    }
}
