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

    @GetMapping("/export/xlsx")
    public void exportToXlsx(HttpServletResponse response) {
        try {
            categoryService.exportToXlsx(response);
        } catch (IOException e) {
            throw new SystemException(AppHttpCodeEnum.EXPORT_FAILED);
        }
    }

}
