package com.anxcye.service.impl;

import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import com.anxcye.domain.entity.Category;
import com.anxcye.service.CategoryService;
import com.anxcye.mapper.CategoryMapper;
import org.springframework.stereotype.Service;

/**
 * @author axy
 * @description 针对表【ancy_category(分类表)】的数据库操作Service实现
 * @createDate 2024-09-05 16:27:33
 */
@Service
public class CategoryServiceImpl extends ServiceImpl<CategoryMapper, Category>
        implements CategoryService {

}




