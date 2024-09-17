package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Comment;
import com.anxcye.service.CommentService;
import com.anxcye.mapper.CommentMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【ancy_comment(评论表)】的数据库操作Service实现
* @createDate 2024-09-17 10:17:36
*/
@Service
public class CommentServiceImpl extends ServiceImpl<CommentMapper, Comment>
    implements CommentService{

}




