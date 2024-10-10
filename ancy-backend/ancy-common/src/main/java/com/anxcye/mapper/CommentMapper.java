package com.anxcye.mapper;

import com.anxcye.domain.entity.Comment;
import com.baomidou.mybatisplus.core.mapper.BaseMapper;

/**
* @author axy
* @description 针对表【ancy_comment(评论表)】的数据库操作Mapper
* @createDate 2024-10-08 14:41:59
* @Entity com.anxcye.domain.entity.Comment
*/
public interface CommentMapper extends BaseMapper<Comment> {

    int updateLikeCount(Long id, Integer increment);

    Long countTotal(Long articleId, String commentType);
}




