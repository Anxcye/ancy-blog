package com.anxcye.service.impl;

import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.vo.categoryVo;
import com.anxcye.service.ArticleService;
import com.anxcye.utils.BeanCopyUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Category;
import com.anxcye.service.CategoryService;
import com.anxcye.mapper.CategoryMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Objects;
import java.util.stream.Collectors;

/**
 * @author axy
 * @description 针对表【ancy_category(分类表)】的数据库操作Service实现
 * @createDate 2024-09-05 16:27:33
 */
@Service
public class CategoryServiceImpl extends ServiceImpl<CategoryMapper, Category> implements CategoryService {

    @Autowired
    private ArticleService articleService;

    @Override
    public List<categoryVo> getList() {
        LambdaQueryWrapper<Article> articleWrapper = new LambdaQueryWrapper<>();
        articleWrapper.select(Article::getId, Article::getCategoryId);
        articleWrapper.eq(Article::getStatus, SystemConstants.ARTICLE_STATUS_NORMAL);
        List<Article> articles = articleService.list(articleWrapper);

        List<Long> categoryIds = articles.stream().map(Article::getCategoryId).distinct().collect(Collectors.toList());

        if (categoryIds.isEmpty()) {
            return List.of();
        }
        LambdaQueryWrapper<Category> categoryWrapper = new LambdaQueryWrapper<>();
        categoryWrapper.in(Category::getId, categoryIds);

        List<Category> categories =  list(categoryWrapper).stream()
                .filter(category -> Objects.equals(category.getStatus(), SystemConstants.CATEGORY_STATUS_NORMAL))
                .collect(Collectors.toList());

        return BeanCopyUtils.copyList(categories, categoryVo.class);
    }
}




