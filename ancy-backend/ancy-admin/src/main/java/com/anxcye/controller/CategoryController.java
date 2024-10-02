package com.anxcye.controller;

import com.anxcye.domain.dto.CategoryDto;
import com.anxcye.domain.dto.CategoryListDto;
import com.anxcye.domain.enums.AppHttpCodeEnum;
import com.anxcye.domain.result.PageResult;
import com.anxcye.domain.result.ResponseResult;
import com.anxcye.domain.vo.CategoryVo;
import com.anxcye.exception.SystemException;
import com.anxcye.service.CategoryService;
import jakarta.servlet.http.HttpServletResponse;
import org.springdoc.core.annotations.ParameterObject;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.access.prepost.PreAuthorize;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.List;

@RestController
@RequestMapping("/category")
public class CategoryController {
    @Autowired
    private CategoryService categoryService;

    @PreAuthorize("@ps.hasPermission('content:category:list')")
    @GetMapping("/list")
    public ResponseResult<List<CategoryVo>> categoryList() {
        return ResponseResult.success(categoryService.getAllCategories());
    }

    @PreAuthorize("@ps.hasPermission('content:category:list')")
    @GetMapping("/page")
    public ResponseResult<PageResult> categoryPage(@ParameterObject CategoryListDto categoryListDto) {
        return ResponseResult.success(categoryService.pageList(categoryListDto));
    }

    @PreAuthorize("@ps.hasPermission('content:category:add')")
    @PostMapping
    public ResponseResult<Long> CategoryAdd(@RequestBody CategoryDto categoryDto) {
        return ResponseResult.success(categoryService.addCategory(categoryDto));
    }

    @PreAuthorize("@ps.hasPermission('content:category:edit')")
    @PutMapping("/{id}")
    public ResponseResult<Boolean> categoryUpdate(@PathVariable Long id, @RequestBody CategoryDto categoryDto) {
        return ResponseResult.success(categoryService.updateCategory(id, categoryDto));
    }

    @PreAuthorize("@ps.hasPermission('content:category:remove')")
    @DeleteMapping("/{id}")
    public ResponseResult<Boolean> CategoryDelete(@PathVariable Long id) {
        return ResponseResult.success(categoryService.deleteCategory(id));
    }

    @PreAuthorize("@ps.hasPermission('content:category:export')")
    @GetMapping("/export/xlsx")
    public void exportToXlsx(HttpServletResponse response) {
        try {
            categoryService.exportToXlsx(response);
        } catch (IOException e) {
            throw new SystemException(AppHttpCodeEnum.EXPORT_FAILED);
        }
    }

}
