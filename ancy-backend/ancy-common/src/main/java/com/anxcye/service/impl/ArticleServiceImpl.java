package com.anxcye.service.impl;


import com.anxcye.constants.RedisConstant;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ArticleCardVo;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.anxcye.mapper.ArticleMapper;
import com.anxcye.service.ArticleService;
import com.anxcye.service.CategoryService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.RedisCache;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.stream.Collectors;

/**
 * @author axy
 * @description 针对表【ancy_article(文章表)】的数据库操作Service实现
 * @createDate 2024-09-05 11:29:50
 */
@Service
public class ArticleServiceImpl extends ServiceImpl<ArticleMapper, Article>
        implements ArticleService {

    @Autowired
    private CategoryService categoryService;

    @Autowired
    private RedisCache redisCache;


    private void updateViewCount(Long id) {
        redisCache.incrementCacheMapValue(RedisConstant.ARTICLE_VIEW_COUNT, id.toString(), 1);
     
    }

    @Override
    public List<HotArticleVo> hot() {
        LambdaQueryWrapper<Article> lambdaQueryWrapper = new LambdaQueryWrapper<>();
        lambdaQueryWrapper.eq(Article::getStatus, SystemConstants.ARTICLE_STATUS_NORMAL);
        lambdaQueryWrapper.orderByDesc(Article::getViewCount);

        Page<Article> page = new Page<>(1, 10);
        page(page, lambdaQueryWrapper);

        return BeanCopyUtils.copyList(page.getRecords(), HotArticleVo.class);
    }

    @Override
    public PageResult getList(Integer pageNum, Integer pageSize, Integer categoryId) {
        LambdaQueryWrapper<Article> articleLambdaQueryWrapper = new LambdaQueryWrapper<>();
        articleLambdaQueryWrapper.eq(Article::getStatus, SystemConstants.ARTICLE_STATUS_NORMAL);
        articleLambdaQueryWrapper.eq(Objects.nonNull(categoryId) && categoryId > 0, Article::getCategoryId, categoryId);
        articleLambdaQueryWrapper.orderByDesc(Article::getIsTop)
                .orderByDesc(Article::getCreateTime);
        Page<Article> page = new Page<>(pageNum, pageSize);
        page(page, articleLambdaQueryWrapper);
        List<ArticleCardVo> articleCardVos = BeanCopyUtils.copyList(page.getRecords(), ArticleCardVo.class);

        articleCardVos.forEach(articleCardVo -> {
            articleCardVo.setCategoryName(categoryService.getById(articleCardVo.getCategoryId()).getName());
        });

        return new PageResult(page.getTotal(), articleCardVos);
    }

    @Override
    public ArticleDetailVo getArticleById(Long id) {
        Article article = getById(id);

        ArticleDetailVo articleDetailVo = BeanCopyUtils.copyBean(article, ArticleDetailVo.class);

        if (Objects.isNull(articleDetailVo)) {
            return null;
        }

        updateViewCount(id);

        articleDetailVo.setCategoryName(categoryService.getById(article.getCategoryId()).getName());

        return articleDetailVo;

    }

    @Override
    public Map<String, Integer> getViewCount() {
        LambdaQueryWrapper<Article> lambdaQueryWrapper = new LambdaQueryWrapper<>();
        lambdaQueryWrapper.select(Article::getId, Article::getViewCount);
        lambdaQueryWrapper.orderByAsc(Article::getId);
        List<Article> articles = list(lambdaQueryWrapper);
        Map<String, Integer> viewCountMap = new HashMap<>();
        articles.forEach(article -> {
            viewCountMap.put(article.getId().toString(), article.getViewCount().intValue());
        });
        return viewCountMap;
    }

    @Override
    public void updateViewCount(Map<String, Integer> viewCountMap) {
        List<Article> articles = viewCountMap.entrySet().stream()
                .map(entry -> new Article(Long.valueOf(entry.getKey()), entry.getValue().longValue()))
                .collect(Collectors.toList());

        updateBatchById(articles);
    }
}
