package com.anxcye.domain.vo;

import lombok.Data;

@Data
public class categoryVo {
    private Long id;
    private String name;
    private Long parentId;
    private String description;
}
