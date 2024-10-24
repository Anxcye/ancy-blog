package com.anxcye.service;

import com.anxcye.domain.dto.CommentDto;
import com.anxcye.domain.dto.CommentPageDto;
import com.anxcye.domain.entity.Comment;
import com.anxcye.domain.result.PageResult;
import com.baomidou.mybatisplus.extension.service.IService;

/**
* @author axy
* @description 针对表【ancy_comment(评论表)】的数据库操作Service
* @createDate 2024-09-17 10:17:36
*/
public interface CommentService extends IService<Comment> {

    PageResult selectComment(String commentType, Long articleId, Integer pageNum, Integer pageSize);

    PageResult getChildren(Long parentId, Integer pageNum, Integer pageSize);

    Long addComment(CommentDto commentDto);

    Boolean updateComment(Long id, CommentDto commentDto);

    Boolean updateCommentLike(Long id, Boolean increase);

    Long countTotal(String commentTypeNote, Long id);

    PageResult getCommentPage(CommentPageDto commentPageDto);

    boolean deleteComment(Long id);
}
