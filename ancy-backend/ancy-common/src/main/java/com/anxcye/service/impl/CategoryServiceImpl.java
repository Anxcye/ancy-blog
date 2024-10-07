package com.anxcye.service.impl;

import com.alibaba.excel.EasyExcel;
import com.anxcye.annotation.Log;
import com.anxcye.constants.SystemConstants;
import com.anxcye.domain.dto.CategoryDto;
import com.anxcye.domain.dto.CategoryListDto;
import com.anxcye.domain.entity.Article;
import com.anxcye.domain.entity.Category;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.vo.CategoryVo;
import com.anxcye.domain.vo.ExcelCategoryVo;
import com.anxcye.exception.SystemException;
import com.anxcye.mapper.CategoryMapper;
import com.anxcye.service.ArticleService;
import com.anxcye.service.CategoryService;
import com.anxcye.utils.BeanCopyUtils;
import com.anxcye.utils.WebUtils;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import java.io.IOException;
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

    @Lazy
    @Autowired
    private ArticleService articleService;

    @Override
    public List<CategoryVo> getUsingCategories() {
        LambdaQueryWrapper<Article> articleWrapper = new LambdaQueryWrapper<>();
        articleWrapper.select(Article::getId, Article::getCategoryId);
        articleWrapper.eq(Article::getStatus, SystemConstants.ARTICLE_STATUS_NORMAL);
        List<Article> articles = articleService.list(articleWrapper);

        List<Long> categoryIds = articles.stream()
                .map(Article::getCategoryId)
                .distinct()
                .collect(Collectors.toList());

        if (categoryIds.isEmpty()) return List.of();

        LambdaQueryWrapper<Category> categoryWrapper = new LambdaQueryWrapper<>();
        categoryWrapper.in(Category::getId, categoryIds);

        List<Category> categories = this.list(categoryWrapper).stream()
                .filter(
                        category ->
                                Objects.equals(category.getStatus(), SystemConstants.CATEGORY_STATUS_NORMAL))
                .collect(Collectors.toList());


        return BeanCopyUtils.copyList(categories, CategoryVo.class);
    }

    @Override
    public List<CategoryVo> getAllCategories() {
        LambdaQueryWrapper<Category> categoryWrapper = new LambdaQueryWrapper<>();
        List<Category> categories = list(categoryWrapper);
        return BeanCopyUtils.copyList(categories, CategoryVo.class);
    }

    @Log
    @Override
    public void exportToXlsx(HttpServletResponse response) throws IOException {
        WebUtils.setDownLoadHeader(SystemConstants.EXPORT_CATEGORY_FILE_NAME, response);

        List<Category> categories = list();
        List<ExcelCategoryVo> excelCategoryVos = BeanCopyUtils.copyList(categories, ExcelCategoryVo.class);
        EasyExcel.write(response.getOutputStream(), ExcelCategoryVo.class)
                .autoCloseStream(Boolean.FALSE)
                .sheet()
                .doWrite(excelCategoryVos);
    }

    @Override
    public PageResult pageList(CategoryListDto categoryListDto) {
        LambdaQueryWrapper<Category> queryWrapper = new LambdaQueryWrapper<>();

        queryWrapper.like(StringUtils.hasText(categoryListDto.getName()), Category::getName, categoryListDto.getName());
        queryWrapper.eq(Objects.nonNull(categoryListDto.getStatus()), Category::getStatus, categoryListDto.getStatus());

        Page<Category> page = new Page<>(categoryListDto.getPageNum(), categoryListDto.getPageSize());
        page(page, queryWrapper);

        List<Category> categories = page.getRecords();
        List<CategoryVo> categoryVos = BeanCopyUtils.copyList(categories, CategoryVo.class);
        return new PageResult(page.getTotal(), categoryVos);
    }

    @Log
    @Override
    public Long addCategory(CategoryDto categoryDto) {
        Category category = BeanCopyUtils.copyBean(categoryDto, Category.class);
        save(category);
        return category.getId();
    }

    @Log
    @Override
    public boolean updateCategory(Long id, CategoryDto categoryDto) {
        Category category = BeanCopyUtils.copyBean(categoryDto, Category.class);
        category.setId(id);
        updateById(category);
        return true;
    }

    @Log
    @Override
    public boolean deleteCategory(Long id) {
        if (!articleService.getByCategoryId(id).isEmpty()){
            throw new SystemException(AppHttpCodeEnum.CATEGORY_EXIST_ARTICLE);
        }
        removeById(id);
        return true;
    }
}




