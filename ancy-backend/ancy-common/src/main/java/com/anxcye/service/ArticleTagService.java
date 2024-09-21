package com.anxcye.service;

import com.anxcye.domain.entity.ArticleTag;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【ancy_article_tag(文章标签关联表)】的数据库操作Service
* @createDate 2024-09-21 15:37:49
*/
public interface ArticleTagService extends IService<ArticleTag> {

    void saveArticleTag(Long id, List<Long> tagIds);
}
