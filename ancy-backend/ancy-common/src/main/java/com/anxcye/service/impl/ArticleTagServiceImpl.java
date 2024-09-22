package com.anxcye.service.impl;

import com.anxcye.domain.entity.ArticleTag;
import com.anxcye.mapper.ArticleTagMapper;
import com.anxcye.service.ArticleTagService;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
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
    public boolean saveArticleTag(Long id, List<Long> tagIds) {
        if (tagIds == null || tagIds.isEmpty()) {
            return true;
        }
        List<ArticleTag> articleTags = tagIds.stream()
                .map(tagId -> new ArticleTag(id, tagId))
                .collect(Collectors.toList());
        saveBatch(articleTags);
        return true;
    }

    @Transactional
    @Override
    public boolean deleteByArticleId(Long articleId) {
        LambdaQueryWrapper<ArticleTag> wrapper = new LambdaQueryWrapper<>();
        wrapper.eq(ArticleTag::getArticleId, articleId);
        return remove(wrapper);
    }

    @Override
    @Transactional
    public boolean updateByArticleId(Long articleId, List<Long> tags) {
        deleteByArticleId(articleId);
        saveArticleTag(articleId, tags);
        return true;
    }
}




