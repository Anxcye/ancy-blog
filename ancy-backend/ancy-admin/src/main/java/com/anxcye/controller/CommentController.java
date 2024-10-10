package com.anxcye.controller;

import com.anxcye.domain.dto.CommentDto;
import com.anxcye.domain.dto.CommentPageDto;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CommentService;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/comment")
public class CommentController {

    @Autowired
    private CommentService commentService;

    @PreAuthorize("@ps.hasPermission('content:comment:list')")
    @GetMapping
    public ResponseResult<PageResult> commentPage(@ParameterObject CommentPageDto commentPageDto){
        return ResponseResult.success(commentService.getCommentPage(commentPageDto));
    }

    @PreAuthorize("@ps.hasPermission('content:comment:add')")
    @PostMapping
    public ResponseResult<Long> commentAdd(@RequestBody CommentDto commentDto) {
        return ResponseResult.success(commentService.addComment(commentDto));
    }

    @PreAuthorize("@ps.hasPermission('content:comment:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> commentUpdate(@PathVariable Long id, @RequestBody CommentDto commentDto) {
        return ResponseResult.success(commentService.updateComment(id, commentDto));
    }

    @PreAuthorize("@ps.hasPermission('content:comment:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> commentDelete(@PathVariable Long id) {
        return ResponseResult.success(commentService.deleteComment(id));
    }

}
