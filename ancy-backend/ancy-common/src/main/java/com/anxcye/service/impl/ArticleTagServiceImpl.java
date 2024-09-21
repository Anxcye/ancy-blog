package com.anxcye.service.impl;

import com.anxcye.domain.entity.ArticleTag;
import com.anxcye.mapper.ArticleTagMapper;
import com.anxcye.service.ArticleTagService;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.stream.Collectors;

/**
* @author axy
* @description 针对表【ancy_article_tag(文章标签关联表)】的数据库操作Service实现
* @createDate 2024-09-21 15:37:49
*/
@Service
public class ArticleTagServiceImpl extends ServiceImpl<ArticleTagMapper, ArticleTag>
    implements ArticleTagService{

    @Override
    @Transactional
    public void saveArticleTag(Long id, List<Long> tagIds) {
        List<ArticleTag> articleTags = tagIds.stream()
                .map(tagId -> new ArticleTag(id, tagId))
                .collect(Collectors.toList());
        saveBatch(articleTags);
    }
}




