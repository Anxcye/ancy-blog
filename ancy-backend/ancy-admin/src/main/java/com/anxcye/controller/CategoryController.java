package com.anxcye.controller;

import com.anxcye.domain.result.ResponseResult;
import com.anxcye.service.CategoryService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/categories")
public class CategoryController {
    @Autowired
    private CategoryService categoryService;

    @GetMapping("/list")
    public ResponseResult<?> list() {
        return ResponseResult.success(categoryService.getAllCategories());
    }
}
