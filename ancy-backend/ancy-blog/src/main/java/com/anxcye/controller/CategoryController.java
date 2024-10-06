package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.CategoryVo;
import com.anxcye.service.CategoryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/category")
public class CategoryController {
    @Autowired
    private CategoryService categoryService;

    @GetMapping("/list")
    public ResponseResult<List<CategoryVo>> categoryList() {
        return ResponseResult.success(categoryService.getUsingCategories());
    }
}
