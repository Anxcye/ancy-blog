package com.anxcye.service;

import com.anxcye.domain.entity.Article;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;
import java.util.Map;

/**
* @author axy
* @description 针对表【ancy_article(文章表)】的数据库操作Service
* @createDate 2024-09-05 11:29:50
*/
public interface ArticleService extends IService<Article> {

    List<HotArticleVo> hot();

    PageResult getList(Integer pageNum, Integer pageSize, Integer categoryId);

    ArticleDetailVo getArticleById(Long id);

    Map<String, Integer> getViewCount();

    void updateViewCount(Map<String, Integer> viewCountMap);
}
