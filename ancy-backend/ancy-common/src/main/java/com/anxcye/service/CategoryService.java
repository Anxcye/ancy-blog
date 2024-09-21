package com.anxcye.service;

import com.anxcye.domain.entity.Category;
import com.anxcye.domain.vo.categoryVo;
import com.baomidou.mybatisplus.extension.service.IService;

import java.util.List;

/**
* @author axy
* @description 针对表【ancy_category(分类表)】的数据库操作Service
* @createDate 2024-09-05 16:27:33
*/
public interface CategoryService extends IService<Category> {

    List<categoryVo> getUsingCategories();

    List<categoryVo> getAllCategories();
}
