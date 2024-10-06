package com.anxcye.service.impl;

import com.alibaba.excel.util.StringUtils;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.CommentDto;
import com.anxcye.domain.entity.Comment;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.CommentVo;
import com.anxcye.exception.SystemException;
import com.anxcye.mapper.CommentMapper;
import com.anxcye.mapper.UserMapper;
import com.anxcye.service.CommentService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;

/**
 * @author axy
 * @description 针对表【ancy_comment(评论表)】的数据库操作Service实现
 * @createDate 2024-09-17 10:17:36
 */
@Service
public class CommentServiceImpl extends ServiceImpl<CommentMapper, Comment>
        implements CommentService {

    @Autowired
    private UserMapper userMapper;


    private List<CommentVo> toCommentVoList(List<Comment> commentList) {
        List<CommentVo> commentVos = BeanCopyUtils.copyList(commentList, CommentVo.class);
        commentVos.forEach(commentVo -> {
            String userName = userMapper.selectById(commentVo.getCreateBy()).getNickName();
            commentVo.setUserName(userName);
            if (commentVo.getRootId() != SystemConstants.COMMENT_IS_ROOT) {
                String toCommentUserName = userMapper.selectById(commentVo.getToCommentUserId()).getNickName();
                commentVo.setToCommentUserName(toCommentUserName);
            } else {
                commentVo.setChildren(getChildren(commentVo.getId(), 1, 3));
            }
        });

        return commentVos;
    }

    @Override
    public PageResult selectComment(String commentType, Long articleId, Integer pageNum, Integer pageSize) {
        LambdaQueryWrapper<Comment> commentLambdaQueryWrapper = new LambdaQueryWrapper<>();
        commentLambdaQueryWrapper.eq(Comment::getStatus, SystemConstants.COMMENT_STATUS_NORMAL);
        commentLambdaQueryWrapper.eq(Objects.nonNull(articleId), Comment::getArticleId, articleId);
        commentLambdaQueryWrapper.eq(Comment::getRootId, SystemConstants.COMMENT_IS_ROOT);
        commentLambdaQueryWrapper.eq(Comment::getType, commentType);
        commentLambdaQueryWrapper.orderByAsc(Comment::getCreateTime);

        Page<Comment> commentPage = new Page<>(pageNum, pageSize);
        page(commentPage, commentLambdaQueryWrapper);

        List<CommentVo> commentVos = toCommentVoList(commentPage.getRecords());
        return new PageResult(commentPage.getTotal(), commentVos);
    }

    @Override
    public PageResult getChildren(Long parentId, Integer pageNum, Integer pageSize) {
        LambdaQueryWrapper<Comment> commentLambdaQueryWrapper = new LambdaQueryWrapper<>();
        commentLambdaQueryWrapper.eq(Comment::getStatus, SystemConstants.COMMENT_STATUS_NORMAL);
        commentLambdaQueryWrapper.eq(Comment::getToCommentId, parentId);
        commentLambdaQueryWrapper.orderByAsc(Comment::getCreateTime);

        Page<Comment> commentPage = new Page<>(pageNum, pageSize);
        page(commentPage, commentLambdaQueryWrapper);
        List<CommentVo> commentVos = toCommentVoList(commentPage.getRecords());
        return new PageResult(commentPage.getTotal(), commentVos);
    }

    @Override
    public Long add(CommentDto commentDto) {
        if (StringUtils.isEmpty(commentDto.getType())) {
            throw new SystemException(AppHttpCodeEnum.CONTENT_NOT_NULL);
        }
        Comment comment = BeanCopyUtils.copyBean(commentDto, Comment.class);
        comment.setStatus(SystemConstants.COMMENT_STATUS_NORMAL);
        save(comment);
        return comment.getId();
    }
}




