package com.anxcye.controller;

import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.exception.SystemException;
import com.anxcye.service.CategoryService;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;

@RestController
@RequestMapping("/categories")
public class CategoryController {
    @Autowired
    private CategoryService categoryService;

    @GetMapping("/list")
    public ResponseResult<?> list() {
        return ResponseResult.success(categoryService.getAllCategories());
    }

//    @GetMapping("/page")
//    public ResponseResult<?> pageList(@ParameterObject CategoryListDto categoryListDto) {
//        return ResponseResult.success(categoryService.pageList(categoryListDto));
//    }
//
//    @PostMapping
//    public ResponseResult<?> addCategory(@RequestBody CategoryDto categoryDto) {
//        return ResponseResult.success(categoryService.addCategory(categoryDto));
//    }
//
//    @PutMapping("/{id}")
//    public ResponseResult<?> updateCategory(@PathVariable Long id, @RequestBody CategoryDto categoryDto) {
//        return ResponseResult.success(categoryService.updateCategory(id, categoryDto));
//    }
//
//    @DeleteMapping("/{id}")
//    public ResponseResult<?> deleteCategory(@PathVariable Long id) {
//        return ResponseResult.success(categoryService.deleteCategory(id));
//    }

    @GetMapping("/export/xlsx")
    public void exportToXlsx(HttpServletResponse response) {
        try {
            categoryService.exportToXlsx(response);
        } catch (IOException e) {
            throw new SystemException(AppHttpCodeEnum.EXPORT_FAILED);
        }
    }

}
