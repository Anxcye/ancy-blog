package com.anxcye.service;

import com.anxcye.domain.dto.ArticleDto;
import com.anxcye.domain.dto.ArticleListDto;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【ancy_article(文章表)】的数据库操作Service
* @createDate 2024-09-05 11:29:50
*/
public interface ArticleService extends IService<Article> {

    List<HotArticleVo> hot();

    PageResult getList(Integer pageNum, Integer pageSize, Integer categoryId);

    ArticleDetailVo getArticleById(Long id);

    void syncFromRedisToDB();

    void initViewCount();

    boolean addArticle(ArticleDto articleDto);

    PageResult pageList(ArticleListDto articleListDto);

    boolean updateArticleById(Long id, ArticleDto addArticleDto);

    boolean deleteArticleById(Long articleId);
}
