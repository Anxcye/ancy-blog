package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;


@Data
public class CategoryDto implements Serializable {
    /**
     * 分类名
     */
    private String name;

    /**
     * 父分类id，如果没有父分类为-1
     */
    private Long parentId;

    /**
     * 描述
     */
    private String description;

    /**
     * 状态0:正常,1禁用
     */
    private String status;
}