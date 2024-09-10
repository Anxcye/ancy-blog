package com.anxcye.service.impl;


import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.ArticleCardVo;
import com.anxcye.domain.vo.HotArticleVo;
import com.anxcye.mapper.ArticleMapper;
import com.anxcye.service.ArticleService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;

/**
 * @author axy
 * @description 针对表【ancy_article(文章表)】的数据库操作Service实现
 * @createDate 2024-09-05 11:29:50
 */
@Service
public class ArticleServiceImpl extends ServiceImpl<ArticleMapper, Article>
        implements ArticleService {

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

        PageResult pageResult = new PageResult(page.getTotal(), articleCardVos);


        return pageResult;
    }
}
