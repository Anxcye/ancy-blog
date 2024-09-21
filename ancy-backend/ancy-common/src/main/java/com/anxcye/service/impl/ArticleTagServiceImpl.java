package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.ArticleTag;
import com.anxcye.service.ArticleTagService;
import com.anxcye.mapper.ArticleTagMapper;
import org.springframework.stereotype.Service;

/**
* @author axy
* @description 针对表【ancy_article_tag(文章标签关联表)】的数据库操作Service实现
* @createDate 2024-09-21 15:37:49
*/
@Service
public class ArticleTagServiceImpl extends ServiceImpl<ArticleTagMapper, ArticleTag>
    implements ArticleTagService{

}




