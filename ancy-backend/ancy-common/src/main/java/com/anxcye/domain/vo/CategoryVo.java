package com.anxcye.domain.vo;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class CategoryVo {
    private Long id;
    private String name;
    private Long parentId;
    private String description;
    private Long createBy;
    private Date createTime;
    private Long updateBy;
    private Date updateTime;
}
