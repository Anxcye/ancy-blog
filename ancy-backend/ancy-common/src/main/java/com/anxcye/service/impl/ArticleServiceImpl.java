package com.anxcye.service.impl;

import com.anxcye.constants.RedisConstant;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.ArticleDto;
import com.anxcye.domain.dto.ArticleListDto;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ArticleCardVo;
import com.anxcye.domain.vo.ArticleDetailVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.anxcye.exception.SystemException;
import com.anxcye.mapper.ArticleMapper;
import com.anxcye.service.ArticleService;
import com.anxcye.service.ArticleTagService;
import com.anxcye.service.CategoryService;
import com.anxcye.service.TagService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.RedisCache;
import com.anxcye.utils.SecurityUtil;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

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

    @Autowired
    private ArticleTagService articleTagService;

    @Autowired
    private TagService tagService;

    private void updateViewCount(Long id) {
        redisCache.incrementCacheMapValue(RedisConstant.ARTICLE_VIEW_COUNT, id.toString(), 1);

    }

    private Long getViewCount(Long id) {
        Integer viewCount = redisCache.getCacheMapValue(RedisConstant.ARTICLE_VIEW_COUNT, id.toString());
        return viewCount.longValue();
    }

    private Map<String, Integer> getViewCount() {
        LambdaQueryWrapper<Article> wrapper = new LambdaQueryWrapper<>();
        wrapper.select(Article::getId, Article::getViewCount);
        wrapper.orderByAsc(Article::getId);
        List<Article> articles = list(wrapper);
        Map<String, Integer> viewCountMap = new HashMap<>();
        articles.forEach(article -> {
            viewCountMap.put(article.getId().toString(), article.getViewCount().intValue());
        });
        return viewCountMap;
    }

    private LambdaQueryWrapper<Article> getArticleWrapper() {
        LambdaQueryWrapper<Article> wrapper = new LambdaQueryWrapper<>();
        if (!SecurityUtil.isAdmin()) {
            wrapper.eq(Article::getStatus, SystemConstants.ARTICLE_STATUS_NORMAL);
        }
        return wrapper;
    }

    @Override
    @Transactional
    public void syncFromRedisToDB() {
        Map<String, Integer> viewCountMap = redisCache.getCacheMap(RedisConstant.ARTICLE_VIEW_COUNT);
        List<Article> articles = viewCountMap.entrySet().stream()
                .map(entry -> new Article(Long.valueOf(entry.getKey()), entry.getValue().longValue()))
                .collect(Collectors.toList());

        // updateBatchById(articles);
        // avoid auto fill updatetime by mybatis plus
        getBaseMapper().updateViewCountById(articles);
    }

    @Override
    @Transactional
    public void initViewCount() {
        Map<String, Integer> viewCountMap = redisCache.getCacheMap(RedisConstant.ARTICLE_VIEW_COUNT);
        if (!viewCountMap.isEmpty()) {
            syncFromRedisToDB();
        }
        viewCountMap = getViewCount();
        redisCache.setCacheMap(RedisConstant.ARTICLE_VIEW_COUNT, viewCountMap);
    }

    @Override
    @Transactional
    public boolean addArticle(ArticleDto articleDto) {
        Article article = BeanCopyUtils.copyBean(articleDto, Article.class);
        save(article);
        initViewCount();

        List<Long> tagIds = articleDto.getTags();
        if (tagIds != null) {
            articleTagService.saveArticleTag(article.getId(), tagIds);
        }

        return true;

    }

    // admin
    @Override
    public PageResult pageList(ArticleListDto articleListDto) {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.like(StringUtils.hasText(articleListDto.getTitle()), Article::getTitle, articleListDto.getTitle());
        wrapper.like(StringUtils.hasText(articleListDto.getSummary()), Article::getSummary,
                articleListDto.getSummary());

        Page<Article> page = new Page<>(articleListDto.getPageNum(), articleListDto.getPageSize());
        page(page, wrapper);

        List<ArticleCardVo> articleCardVos = BeanCopyUtils.copyList(page.getRecords(), ArticleCardVo.class);

        articleCardVos.forEach(articleCardVo -> {
            articleCardVo.setViewCount(getViewCount(articleCardVo.getId()));
            articleCardVo.setCategoryName(categoryService.getById(articleCardVo.getCategoryId()).getName());
            articleCardVo.setTags(tagService.selectTagsByArticleId(articleCardVo.getId()));
        });

        return new PageResult(page.getTotal(), articleCardVos);
    }

    @Override
    @Transactional
    public boolean updateArticleById(Long id, ArticleDto addArticleDto) {
        Article article = BeanCopyUtils.copyBean(addArticleDto, Article.class);
        article.setId(id);
        updateById(article);
        if (addArticleDto.getTags() != null) {
            articleTagService.updateByArticleId(id, addArticleDto.getTags());
        }
        return true;
    }

    @Override
    public boolean deleteArticleById(Long articleId) {
        removeById(articleId);
        return true;
    }

    @Override
    public List<ArticleCardVo> getByCategoryId(Long categoryId) {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getCategoryId, categoryId);
        List<Article> articles = list(wrapper);
        return BeanCopyUtils.copyList(articles, ArticleCardVo.class);
    }

    @Override
    public List<ArticleCardVo> getArticleFront() {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getType, SystemConstants.ARTICLE_TYPE_FRONT);
        wrapper.orderByAsc(Article::getOrderNum);

        List<Article> articles = list(wrapper);
        return BeanCopyUtils.copyList(articles, ArticleCardVo.class);
    }

    @Override
    public List<ArticleCardVo> recent() {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getType, SystemConstants.ARTICLE_TYPE_NORMAL);
        wrapper.orderByDesc(Article::getCreateTime);
        wrapper.last("limit 5");
        return BeanCopyUtils.copyList(list(wrapper), ArticleCardVo.class);
    }



    @Override
    public List<HotArticleVo> hot() {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();

        wrapper.eq(Article::getType, SystemConstants.ARTICLE_TYPE_NORMAL);
        wrapper.orderByDesc(Article::getViewCount);

        Page<Article> page = new Page<>(1, 10);
        page(page, wrapper);

        List<HotArticleVo> hotArticleVos = BeanCopyUtils.copyList(page.getRecords(), HotArticleVo.class);

        hotArticleVos.forEach(hotArticleVo -> {
            hotArticleVo.setViewCount(getViewCount(hotArticleVo.getId()));
        });

        return hotArticleVos;
    }

    // blog
    @Override
    public PageResult getList(Integer pageNum, Integer pageSize, Integer categoryId) {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getType, SystemConstants.ARTICLE_TYPE_NORMAL);
        wrapper.eq(Objects.nonNull(categoryId) && categoryId > 0,
                Article::getCategoryId, categoryId);
        wrapper.orderByDesc(Article::getIsTop)
                .orderByAsc(Article::getOrderNum)
                .orderByDesc(Article::getCreateTime);
        Page<Article> page = new Page<>(pageNum, pageSize);
        page(page, wrapper);
        List<ArticleCardVo> articleCardVos = BeanCopyUtils.copyList(page.getRecords(), ArticleCardVo.class);

        articleCardVos.forEach(articleCardVo -> {
            articleCardVo.setCategoryName(categoryService.getById(articleCardVo.getCategoryId()).getName());
            articleCardVo.setViewCount(getViewCount(articleCardVo.getId()));
            articleCardVo.setTags(tagService.selectTagsByArticleId(articleCardVo.getId()));
        });

        return new PageResult(page.getTotal(), articleCardVos);
    }

    private ArticleDetailVo getDetailById(Long id, Integer type) {
           LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getId, id);

        wrapper.eq(Article::getType, type);

        Article article = getOne(wrapper);

        ArticleDetailVo articleDetailVo = BeanCopyUtils.copyBean(article, ArticleDetailVo.class);

        if (Objects.isNull(articleDetailVo)) {
            throw new SystemException(AppHttpCodeEnum.NOT_FOUND);
        }

        updateViewCount(id);

        articleDetailVo.setViewCount(getViewCount(id));

        articleDetailVo.setCategoryName(categoryService.getById(article.getCategoryId()).getName());

        articleDetailVo.setTags(tagService.selectTagsByArticleId(id));

        return articleDetailVo;
    }

    @Override
    public ArticleDetailVo getArticleById(Long id) {
      return getDetailById(id, SystemConstants.ARTICLE_TYPE_NORMAL);
    }

    @Override
    public ArticleDetailVo getHomeById(Long id) {
        return getDetailById(id, SystemConstants.ARTICLE_TYPE_FRONT);
    }

    @Override
    public ArticleDetailVo getArticleLink() {
        LambdaQueryWrapper<Article> wrapper = getArticleWrapper();
        wrapper.eq(Article::getType, SystemConstants.ARTICLE_TYPE_LINK);
        Article article = getOne(wrapper);
        return BeanCopyUtils.copyBean(article, ArticleDetailVo.class);
    }
}
