package com.anxcye.domain.vo;

import com.alibaba.excel.annotation.ExcelProperty;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class ExcelCategoryVo {
    @ExcelProperty("分类ID")
    private Long id;
    @ExcelProperty("分类名称")
    private String name;
    @ExcelProperty("父分类ID")
    private Long parentId;
    @ExcelProperty("分类描述")
    private String description;
}
