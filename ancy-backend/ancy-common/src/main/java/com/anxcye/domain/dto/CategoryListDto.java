package com.anxcye.domain.dto;

import lombok.Data;

import java.io.Serializable;

@Data
public class CategoryListDto extends PageListDto implements Serializable {
    private String name;
    private String status;
}